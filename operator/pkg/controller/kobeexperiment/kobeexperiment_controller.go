package kobeexperiment

import (
	"context"

	kobebenchmarkv1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobebenchmark/v1alpha1"
	kobedatasetv1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobedataset/v1alpha1"
	kobeexperimentv1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobeexperiment/v1alpha1"
	kobefederatorv1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobefederator/v1alpha1"

	batchv1 "k8s.io/api/batch/v1"
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

var log = logf.Log.WithName("controller_kobeexperiment")

// Add creates a new KobeExperiment Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileKobeExperiment{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("kobeexperiment-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource KobeExperiment
	err = c.Watch(&source.Kind{Type: &kobeexperimentv1alpha1.KobeExperiment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner KobeExperiment
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobeexperimentv1alpha1.KobeExperiment{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileKobeExperiment implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileKobeExperiment{}

// ReconcileKobeExperiment reconciles a KobeExperiment object
type ReconcileKobeExperiment struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a KobeExperiment object and makes changes based on the state read
// and what is in the KobeExperiment.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileKobeExperiment) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling KobeExperiment")

	// Fetch the KobeExperiment instance
	instance := &kobeexperimentv1alpha1.KobeExperiment{}
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
	//normally i have to check for finishing initialization not just if they exist.Federator for example could be initiliazing with its init container
	//check if there exist a benchmark with this name
	foundBenchmark := &kobebenchmarkv1alpha1.KobeBenchmark{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.Benchmark, Namespace: instance.Namespace}, foundBenchmark)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Did not found a kobebenchmark resource with this name please define that first")
		return reconcile.Result{}, err
	}
	//check if every kobedataset of the benchmark possible obsolete of this experiment is up and running
	for _, datasetInfo := range foundBenchmark.Spec.Datasets {
		dataset := &kobedatasetv1alpha1.KobeDataset{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: foundBenchmark.Namespace, Name: datasetInfo.Name}, dataset)
		if err != nil {
			reqLogger.Info("Failed to find a specific dataset from the list of datasets of this benchmark")
			return reconcile.Result{}, err
		}
	}

	//check if every federator is up and running
	for _, federatorInfo := range instance.Spec.Federators {
		foundFederator := &kobefederatorv1alpha1.KobeFederator{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: instance.Namespace, Name: federatorInfo}, foundFederator)
		if err != nil {
			reqLogger.Info("Failed to find a specific federator from the list of federators of this experiment")
			return reconcile.Result{}, err
		}
	}
	//Create more stuff here that are necessary to init the federators and the queries
	//get the queries and mount them so the job can find them
	//Everything is running and ready for the experiment
	if instance.Spec.RunFlag == false { //dont run just yet just have it defined
		return reconcile.Result{}, nil
	}
	//Create the new job that will run for this experiment
	foundJob := &batchv1.Job{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Namespace: instance.Namespace, Name: "kobejob"}, foundJob)
	if err != nil && errors.IsNotFound(err) {
		//to fix this .. do something else here cause there should be multiple jobs and there should be no problem
		reqLogger.Info("WHY AM I EVEN HERE !!!!?????? /n")
		reqLogger.Info("There is still a job\n ")
		return reconcile.Result{}, err
	}
	experimentJob := r.newJobForExperiment(instance)
	reqLogger.Info("Creating a new job to run the experiment for this setup")
	err = r.client.Create(context.TODO(), experimentJob)
	if err != nil {
		reqLogger.Info("FAILED to create the job to run this expriment  %s/%s\n", experimentJob.Name, experimentJob.Namespace)
		return reconcile.Result{}, err
	}
	reqLogger.Info("Reached the end of the reconciling loop for the kobe Experiment %s/%s\n", instance.Name, instance.Namespace)

	return reconcile.Result{}, err
}

func (r *ReconcileKobeExperiment) newJobForExperiment(m *kobeexperimentv1alpha1.KobeExperiment) *batchv1.Job {
	times := int32(1)
	parallelism := int32(1)
	job := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kobejob",
			Namespace: m.Namespace,
		},
		Spec: batchv1.JobSpec{
			Parallelism: &parallelism,
			Completions: &times,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"name": "kobeJobPod"},
			},
			Template: corev1.PodTemplateSpec{
				metav1.ObjectMeta{
					Name:      "kobeJobPod",
					Namespace: m.Namespace,
					Labels:    map[string]string{"name": "kobeJobPod"},
				},
				corev1.PodSpec{ //should these be hardfixed??
					Containers: []corev1.Container{{
						Image:           "busybox", //this is gonna be the image of client program
						Name:            "client name not important",
						ImagePullPolicy: "Never",
						Ports: []corev1.ContainerPort{{
							ContainerPort: int32(8990),
							Name:          "port of the client ",
						}},
					}},
				},
			},
		},
	}
	controllerutil.SetControllerReference(m, job, r.scheme)
	return job

}
