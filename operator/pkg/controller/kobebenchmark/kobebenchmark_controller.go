package kobebenchmark

import (
	"context"

	kobebenchmarkv1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobebenchmark/v1alpha1"
	kobedatasetv1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobedataset/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_kobebenchmark")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new KobeBenchmark Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileKobeBenchmark{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("kobebenchmark-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource KobeBenchmark
	err = c.Watch(&source.Kind{Type: &kobebenchmarkv1alpha1.KobeBenchmark{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner KobeBenchmark
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobebenchmarkv1alpha1.KobeBenchmark{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &kobedatasetv1alpha1.KobeDataset{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobebenchmarkv1alpha1.KobeBenchmark{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileKobeBenchmark implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileKobeBenchmark{}

// ReconcileKobeBenchmark reconciles a KobeBenchmark object
type ReconcileKobeBenchmark struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a KobeBenchmark object and makes changes based on the state read
// and what is in the KobeBenchmark.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileKobeBenchmark) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling KobeBenchmark")

	// Fetch the KobeBenchmark instance
	instance := &kobebenchmarkv1alpha1.KobeBenchmark{}
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
	foundDataset := &kobedatasetv1alpha1.KobeDataset{}
	for _, dataset := range instance.Spec.Datasets {

		err = r.client.Get(context.TODO(), types.NamespacedName{Name: dataset.Name, Namespace: instance.Namespace}, foundDataset)
		if err != nil && errors.IsNotFound(err) {
			// Define a new deployment
			kobedataset := r.newKobeDataset(&dataset, instance)
			reqLogger.Info("Creating a new basic KobeDataset %s/%s\n", kobedataset.Namespace, kobedataset.Name)
			err = r.client.Create(context.TODO(), kobedataset)
			if err != nil {
				reqLogger.Info("Failed to create new Kobedataset: %v\n", err)
				return reconcile.Result{}, err
			}
			// Kobedataset created successfully - return and requeue
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Info("Failed to get KobeDataset with the same name in same namespace: %v\n", err)
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

	return reconcile.Result{}, err
}
func labelsForKobeBenchmark(name string) map[string]string {
	return map[string]string{"app": "Kobe-Operator", "kobeoperator_cr": name}
}
func (r *ReconcileKobeBenchmark) newConfigMapForQueries(m *kobebenchmarkv1alpha1.KobeBenchmark, querymap map[string]string) *corev1.ConfigMap {
	configmap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},

		Data: querymap,
	}
	controllerutil.SetControllerReference(m, configmap, r.scheme)
	return configmap
}
func (r *ReconcileKobeBenchmark) newKobeDataset(dataset *kobebenchmarkv1alpha1.Dataset, m *kobebenchmarkv1alpha1.KobeBenchmark) *kobedatasetv1alpha1.KobeDataset {

	data := &kobedatasetv1alpha1.KobeDataset{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kobedataset.kobe.com/v1alpha1",
			Kind:       "KobeDataset",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      dataset.Name,
			Namespace: m.Namespace,
		},
		Spec: kobedatasetv1alpha1.KobeDatasetSpec{
			Image:           dataset.Image,
			DownloadFrom:    dataset.DownloadFrom,
			ImagePullPolicy: "Always",
			Count:           1,
			Group:           "kobedataset.kobe.com",
			Port:            80,
		},
	}
	// Set Examplekind instance as the owner and controller
	controllerutil.SetControllerReference(m, data, r.scheme)
	return data

}
