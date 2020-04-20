package benchmark

import (
	"context"

	api "github.com/semagrow/kobe/operator/pkg/apis/kobe/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Benchmark
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &api.Benchmark{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &api.Dataset{}}, &handler.EnqueueRequestForOwner{
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
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
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

	//check if the datasets exist else create a very basic version of those that dont
	foundDataset := &api.Dataset{}
	for _, dataset := range instance.Spec.Datasets {

		err = r.client.Get(context.TODO(), types.NamespacedName{Name: dataset.Name, Namespace: instance.Namespace}, foundDataset)
		if err != nil && errors.IsNotFound(err) {
			// Define a new deployment
			/*
				kobedataset := r.newKobeDataset(&dataset, instance)
				reqLogger.Info("Creating a new basic Dataset %s/%s\n", kobedataset.Namespace, kobedataset.Name)
				err = r.client.Create(context.TODO(), kobedataset)
				if err != nil {
					reqLogger.Info("Failed to create new Dataset: %v\n", err)
					return reconcile.Result{}, err
				}
				// Kobedataset created successfully - return and requeue
				return reconcile.Result{Requeue: true}, nil
			*/
			return reconcile.Result{}, err
		} else if err != nil {
			reqLogger.Info("Failed to get Dataset with the same name in same namespace: %v\n", err)
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

/*
func (r *ReconcileBenchmark) newKobeDataset(dataset *api.Dataset, m *api.Benchmark) *api.Dataset {

	data := &api.Dataset{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dataset.Name,
			Namespace: m.Namespace,
		},
		Spec: kobev1alpha1.KobeDatasetSpec{
			Image:           dataset.Image,
			DownloadFrom:    dataset.DownloadFrom,
			ImagePullPolicy: "Always",
			Port:            80,
		},
	}
	// Set kobe benchmark instance as the owner and controller
	controllerutil.SetControllerReference(m, data, r.scheme)
	return data
}
*/
