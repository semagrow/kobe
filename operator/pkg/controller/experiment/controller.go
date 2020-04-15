package experiment

import (
	"context"
	"fmt"
	"strconv"

	api "github.com/semagrow/kobe/operator/pkg/apis/kobe/v1alpha1"
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
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_experiment")
var identifier = 0

// Add creates a new Experiment Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileExperiment{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("experiment-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Experiment
	err = c.Watch(&source.Kind{Type: &api.Experiment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &batchv1.Job{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &api.Experiment{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &api.Experiment{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileExperiment implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileExperiment{}

// ReconcileExperiment reconciles a Experiment object
type ReconcileExperiment struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Experiment object and makes changes based on the state read
// and what is in the Experiment.Spec
func (r *ReconcileExperiment) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Experiment")

	// Fetch the Experiment instance
	instance := &api.Experiment{}
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

	// Normally I have to check for finishing initialization not just if they
	// exist. Federator for example could be initiliazing with its init container.
	// very important

	// Check if there exist a kobe benchmark with this name in Kubernetes.
	// If not its an error.
	foundBenchmark := &api.Benchmark{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.Benchmark, Namespace: instance.Namespace}, foundBenchmark)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Did not found a kobebenchmark resource with this name please define that first")
		return reconcile.Result{RequeueAfter: 5}, err
	}
	endpoints := []string{}
	datasets := []string{}

	// Check if every kobedataset of the benchmark is healthy.
	// Create a list of the endpoints and of the names of the datasets
	for _, datasetInfo := range foundBenchmark.Spec.Datasets {
		foundDataset := &api.Dataset{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: foundBenchmark.Namespace, Name: datasetInfo}, foundDataset)
		if err != nil {
			reqLogger.Info("Failed to find a specific dataset from the list of datasets of this benchmark")
			return reconcile.Result{RequeueAfter: 5}, err
		}

		// Check for the healthiness of the individual pods of the kobe dataset
		podList := &corev1.PodList{}
		listOps := []client.ListOption{
			client.InNamespace(instance.Namespace),
			client.MatchingLabels{"kobeoperator_cr": foundDataset.Name},
		}
		err = r.client.List(context.TODO(), podList, listOps...)
		if err != nil {
			reqLogger.Info("Failed to list pods: %v", err)
			return reconcile.Result{}, err
		}

		podNames := getPodNames(podList.Items)
		for _, podname := range podNames {
			foundPod := &corev1.Pod{}
			err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: instance.Namespace, Name: podname}, foundPod)
			if err != nil && errors.IsNotFound(err) {
				reqLogger.Info("Failed to get the pod of the kobe dataset that experiment will use")
				return reconcile.Result{RequeueAfter: 5}, nil
			}
			var test string
			test = string(foundPod.Status.Phase)
			if test != "Running" {
				reqLogger.Info("Kobe dataset pod is not ready so experiment needs to wait")
				return reconcile.Result{RequeueAfter: 5}, nil
			}
		}

		if podNames == nil || len(podNames) == 0 {
			reqLogger.Info("Experiment waits for components initialization")
			return reconcile.Result{RequeueAfter: 5}, nil
		}

		// Create a list of the SPARQL endpoints
		endpoints = append(endpoints,
			EndpointURL(foundDataset.Name, foundDataset.Namespace, int(foundDataset.Spec.Port), foundDataset.Spec.Path))
		datasets = append(datasets, foundDataset.Name)
	}

	foundFederator := &api.Federator{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.Federator, Namespace: instance.Namespace}, foundFederator)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("No federator with this name is defined in the cluster")
		return reconcile.Result{}, err
	}
	if err != nil {
		reqLogger.Info("Error at getting this federator resource from the cluster.")
		return reconcile.Result{}, err
	}

	foundFederation := &api.Federation{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundFederation)
	if err != nil && errors.IsNotFound(err) {
		newFederation := r.newFederationForExperiment(instance, foundFederator, endpoints, datasets)
		reqLogger.Info("Creating a new federation based on this experiments datasets and federator")
		err = r.client.Create(context.TODO(), newFederation)
		if err != nil {
			reqLogger.Info("Failed to create the federation")
			return reconcile.Result{}, err
		}
	}
	if err != nil {
		return reconcile.Result{}, err
	}

	// Check if the pods of the federators exist and have a status of running
	// before proceeding and get fed name and endpoint for the eval job
	fedEndpoint :=
		EndpointURL(foundFederation.Name, foundFederation.Namespace, int(foundFederation.Spec.Template.Port), foundFederation.Spec.Template.Path)
	fedName := foundFederation.Name

	podList := &corev1.PodList{}
	listOps := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels{"kobeoperator_cr": instance.Name},
	}
	err = r.client.List(context.TODO(), podList, listOps...)
	if err != nil {
		reqLogger.Info("Failed to list pods: %v", err)
		return reconcile.Result{}, err
	}
	podNames := getPodNames(podList.Items)
	//for _, podname := range foundFederation.Status.PodNames {
	for _, podname := range podNames {
		foundPod := &corev1.Pod{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: instance.Namespace, Name: podname}, foundPod)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Failed to get the pod of the kobe federation that experiment will use")
			return reconcile.Result{RequeueAfter: 25}, nil
		}
		var test string
		test = string(foundPod.Status.Phase)
		if test != "Running" {
			reqLogger.Info("Kobe federation pod is not ready so experiment needs to wait")
			return reconcile.Result{RequeueAfter: 25}, nil
		}
	}
	if podNames == nil || len(podNames) == 0 {
		reqLogger.Info("Experiment waits for components initialization")
		return reconcile.Result{RequeueAfter: 25}, nil
	}

	// Everything is healthy and ready for the experiment.
	if instance.Spec.DryRun == true {
		// dont run just yet just have it defined
		return reconcile.Result{}, nil
	}

	// Create the new job that will run the EVAL client for this experiment
	if instance.Spec.TimesToRun > 0 {
		foundJob := &batchv1.Job{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name + "-" + strconv.Itoa(identifier)}, foundJob)
		if err == nil {
			if &foundJob.Status.Succeeded == nil || foundJob.Status.Succeeded == 0 {
				return reconcile.Result{}, nil
			}
			reqLogger.Info("All past jobs are done\n")
			identifier++
		}
		experimentJob := r.createEvaluatorJob(instance, identifier, fedEndpoint, fedName)
		reqLogger.Info("Creating a new job to run the experiment for this setup")
		err = r.client.Create(context.TODO(), experimentJob)
		if err != nil {
			reqLogger.Info("FAILED to create the job to run this experiment  %s/%s\n", experimentJob.Name, experimentJob.Namespace)
			return reconcile.Result{}, err
		}
		instance.Spec.TimesToRun = instance.Spec.TimesToRun - 1
		err = r.client.Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Info("Failed to update the times to run of the experiment")
			return reconcile.Result{}, err
		}
		reqLogger.Info("Reached the end of the reconciling loop for the kobe Experiment %s/%s\n", instance.Name, instance.Namespace)
	}
	return reconcile.Result{}, nil
}

//----------------------functions that create native kubernetes objects--------------------------------------
//create the job that will run the evaluation program
func (r *ReconcileExperiment) createEvaluatorJob(m *api.Experiment, i int, fedendpoint string, fedname string) *batchv1.Job {
	times := int32(1)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-" + strconv.Itoa(i),
			Namespace: m.Namespace,
		},
		Spec: batchv1.JobSpec{
			Parallelism: &m.Spec.Evaluator.Parallelism,
			Completions: &times,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           m.Spec.Evaluator.Image, //this is the image of the eval
						ImagePullPolicy: m.Spec.Evaluator.ImagePullPolicy,
						Name:            "job" + "-" + strconv.Itoa(i),
						Command:         m.Spec.Evaluator.Command,
						Ports: []corev1.ContainerPort{{
							ContainerPort: int32(8890), //eval endpoint
							Name:          "client",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "queries",
							MountPath: "/queries",
						}},
						Env: []corev1.EnvVar{
							{Name: "FEDERATION_NAME", Value: fedname},
							{Name: "FEDERATION_ENDPOINT", Value: fedendpoint},
							{Name: "ENDPOINT", Value: fedendpoint},
							{Name: "EVAL_RUN", Value: "1"},
						},
					}},
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Volumes: []corev1.Volume{{
						Name: "queries",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: m.Spec.Benchmark,
								}}}}},
				},
			},
		},
	}
	controllerutil.SetControllerReference(m, job, r.scheme)
	return job
}

//function that creates a new kobefederation custom resource from the federator and benchmark  in experiment.
//The native objects that kobefederation needs are created by kobefederation controller .
func (r *ReconcileExperiment) newFederationForExperiment(m *api.Experiment,
	fed *api.Federator, endpoints []string, datasetnames []string) *api.Federation {

	federation := &api.Federation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: fed.Namespace,
		},
		Spec: api.FederationSpec{
			Template:      fed.Spec.FederatorTemplate,
			ForceNewInit:  m.Spec.ForceNewInit,
			Init:          true,
			FederatorName: fed.Name,
			Endpoints:     endpoints,
			Datasets:      datasetnames,
		},
	}
	controllerutil.SetControllerReference(m, federation, r.scheme)
	return federation
}

func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

func EndpointURL(name, namespace string, port int, path string) string {
	// "http://"+foundDataset.Name+"."+foundDataset.Namespace+".svc.cluster.local"+":"+strconv.Itoa(int(foundDataset.Spec.Port))+foundDataset.Spec.Path)
	return fmt.Sprintf("http://%s.%s.svc:%d/%s", name, namespace, port, path)
}
