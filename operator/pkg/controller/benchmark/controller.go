package benchmark

import (
	"context"

	api "github.com/semagrow/kobe/operator/pkg/apis/kobe/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_benchmark")

// Add creates a new Benchmark Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileBenchmark{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("benchmark-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Benchmark
	err = c.Watch(&source.Kind{Type: &api.Benchmark{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &api.EphemeralDataset{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &api.Benchmark{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileBenchmark implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileBenchmark{}

// ReconcileBenchmark reconciles a Benchmark object
type ReconcileBenchmark struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Benchmark object and makes changes based on the state read
// and what is in the Benchmark.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileBenchmark) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Benchmark")

	// Fetch the KobeBenchmark instance
	instance := &api.Benchmark{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      instance.Name,
		Namespace: ""},
		instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// check  if a KobeUtil instance exists for this namespace and if not create it
	requeue, err := ensureNFS(r.client, corev1.NamespaceDefault)
	if requeue {
		return reconcile.Result{Requeue: requeue}, err
	} else if err != nil {
		return reconcile.Result{}, err
	}

	//create the new namespace
	//istio label for the namespace
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: instance.Name, Labels: map[string]string{"istio-injection": "enabled"}}}
	config, err := clientcmd.BuildConfigFromFlags("", "")
	// if err != nil {
	// 	reqLogger.Info("Failed client connection: %v\n", err)
	// 	//return reconcile.Result{Requeue: requeue}, err
	// }

	clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	reqLogger.Info("Failed client connection: %v\n", err)
	// 	//return reconcile.Result{Requeue: requeue}, err
	// }
	_, err = clientset.CoreV1().Namespaces().Create(ns)
	// if err != nil {
	// 	reqLogger.Info("Failed client connection: %v\n", err)
	// 	//return reconcile.Result{Requeue: requeue}, err
	// }

	//add finalizer to the resource . If the benchmark gets deleted the finalizer logic deletes the entire benchmark
	nsFinalizer := "delete.the.fking.ns.kobe"

	// examine DeletionTimestamp to determine if object is under deletion
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !containsString(instance.ObjectMeta.Finalizers, nsFinalizer) {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, nsFinalizer)
			if err := r.client.Update(context.Background(), instance); err != nil {
				return reconcile.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if containsString(instance.ObjectMeta.Finalizers, nsFinalizer) {
			// our finalizer is present, so lets handle any external dependency
			propagation := metav1.DeletePropagationBackground
			err = clientset.CoreV1().Namespaces().Delete(ns.Name, &metav1.DeleteOptions{PropagationPolicy: &propagation})

			// remove our finalizer from the list and update it.
			instance.ObjectMeta.Finalizers = removeString(instance.ObjectMeta.Finalizers, nsFinalizer)
			if err := r.client.Update(context.Background(), instance); err != nil {
				return reconcile.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return reconcile.Result{}, nil
	}

	//check if the datasets exist else create them and let dataset controller build the resources
	foundDataset := &api.EphemeralDataset{}
	for _, dataset := range instance.Spec.Datasets {

		err = r.client.Get(context.TODO(), types.NamespacedName{Name: dataset.Name, Namespace: instance.Name}, foundDataset)
		if err != nil && errors.IsNotFound(err) {
			ed := r.newEphemeralDataset(instance, dataset)
			err := r.client.Create(context.TODO(), ed)
			return reconcile.Result{}, err
		} else if err != nil {
			reqLogger.Info("Failed to get the ephemeral dataset in the namespace of the benchmark: %v\n", err)
			return reconcile.Result{}, err
		}
	}

	//check if config map exists else create it
	//config map contains the queries assosciated with this benchmark setup in seperate files .
	foundConfig := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundConfig)
	if err != nil && errors.IsNotFound(err) {
		if instance.Spec.Queries == nil {
			return reconcile.Result{}, err
		}
		//create a new config map from the queries that are defined in the yaml of this benchmark
		querymap := map[string]string{}
		for _, query := range instance.Spec.Queries {
			querymap[query.Name] = query.QueryString
		}
		configMap := r.newConfigMapForQueries(instance, querymap)
		err := r.client.Create(context.TODO(), configMap)
		if err != nil {
			reqLogger.Info("FAILED to create the configmap for this set of queries for the benchmark")
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}
	reqLogger.Info("FINISHED RECONCILING LOOP FOR BENCHMARK SUCCESSFULLY")
	return reconcile.Result{}, nil
}

func labelsForKobeBenchmark(name string) map[string]string {
	return map[string]string{"app": "Kobe-Operator", "kobeoperator_cr": name}
}

func (r *ReconcileBenchmark) newConfigMapForQueries(m *api.Benchmark, querymap map[string]string) *corev1.ConfigMap {
	configmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Data: querymap,
	}
	controllerutil.SetControllerReference(m, configmap, r.scheme)
	return configmap
}

func ensureNFS(client client.Client, ns string) (bool, error) {
	reqLogger := log

	// check  if a KobeUtil instance exists for this namespace and if not create it
	kobeUtil := &api.KobeUtil{}

	err := client.Get(context.TODO(), types.NamespacedName{Name: "kobeutil", Namespace: ns}, kobeUtil)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating the kobe utility custom resource")
		kobeutil := newKobeUtility(ns)
		err = client.Create(context.TODO(), kobeutil)
		if err != nil {
			reqLogger.Info("Failed to create the kobe utility instance: %v\n", err)
			return false, err
		}
		return true, nil
	}

	//check for  the nfs pod if it exist and wait if not
	nfsPodFound := &corev1.Pod{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: "kobenfs", Namespace: ns}, nfsPodFound)
	if err != nil && errors.IsNotFound(err) {
		return true, nil
	}
	//check if the persistent volume claim exist and wait if not
	pvcFound := &corev1.PersistentVolumeClaim{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: "kobepvc", Namespace: ns}, pvcFound)
	if err != nil && errors.IsNotFound(err) {
		return true, nil
	}
	//make sure nfs pod status is running, to take care of racing condition
	if nfsPodFound.Status.Phase != corev1.PodRunning {
		return true, nil
	}

	return false, nil
}

func newKobeUtility(ns string) *api.KobeUtil {
	kutil := &api.KobeUtil{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kobeutil",
			Namespace: ns,
		},
	}
	return kutil
}

func (r *ReconcileBenchmark) newEphemeralDataset(benchmark *api.Benchmark, dataset api.Dataset) *api.EphemeralDataset {
	ephemeralDataset := &api.EphemeralDataset{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dataset.Name,
			Namespace: benchmark.Name,
		},

		Spec: dataset,
	}
	controllerutil.SetControllerReference(benchmark, ephemeralDataset, r.scheme)
	return ephemeralDataset
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
