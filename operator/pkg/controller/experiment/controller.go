package experiment

import (
	"context"
	"strconv"

	api "github.com/semagrow/kobe/operator/pkg/apis/kobe/v1alpha1"
	"github.com/semagrow/kobe/operator/pkg/util"
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

const ControllerName = "kobe-experiment-controller"

var log = logf.Log.WithName(ControllerName)
var identifier = 0

// Add creates a new Experiment Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileExperiment{
		client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Experiment
	err = c.Watch(&source.Kind{Type: &api.Experiment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &api.Federation{}}, &handler.EnqueueRequestForOwner{
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

	err = c.Watch(&source.Kind{Type: &batchv1.Job{}}, &handler.EnqueueRequestForOwner{
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

	//add finalizer to the resource . If the experiments gets deleted the finalizer logic deletes the federation
	//need this cause it belongs to different namespace and doesnt seem to care that experiment is the father of the federation..

	fedFinalizer := "delete.the.fking.fed.kobe"

	// examine DeletionTimestamp to determine if object is under deletion
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containsString(instance.ObjectMeta.Finalizers, fedFinalizer) {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, fedFinalizer)
			if err := r.client.Update(context.Background(), instance); err != nil {
				return reconcile.Result{}, err
			}
		}
	} else {
		if containsString(instance.ObjectMeta.Finalizers, fedFinalizer) {
			foundFederation := &api.Federation{}
			err := r.client.Get(context.TODO(), types.NamespacedName{
				Name:      instance.Spec.FederatorName,
				Namespace: instance.Spec.Benchmark}, foundFederation)

			if err == nil {
				err = r.client.Delete(context.TODO(), foundFederation, client.PropagationPolicy(metav1.DeletionPropagation("Background")))
			}
			instance.ObjectMeta.Finalizers = removeString(instance.ObjectMeta.Finalizers, fedFinalizer)
			if err := r.client.Update(context.Background(), instance); err != nil {
				return reconcile.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return reconcile.Result{}, nil
	}
	//check if the assosciated benchmark component exists in the experiment namespace
	foundBenchmark := &api.Benchmark{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.Benchmark, Namespace: instance.Namespace}, foundBenchmark)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("The benchmark of this experiment does not exist" + instance.Spec.Benchmark)
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	//check if all datasets for the experiment are up and running
	//testsetestset
	requeue, err := r.reconcileDatasets(instance, *foundBenchmark)
	if requeue {
		reqLogger.Info("Datasets are not up yet. Wait 10 sec and  and then check again!\n")
		return reconcile.Result{RequeueAfter: 10000000000}, err
	} else if err != nil {
		return reconcile.Result{}, err
	}
	//reconcile federation also ensures federator is up and running
	requeue, err = r.reconcileFederation(instance)
	if requeue {
		return reconcile.Result{RequeueAfter: 2000000000}, err
	} else if err != nil {
		return reconcile.Result{RequeueAfter: 2000000000}, err
	}

	// Everything is healthy and ready for the experiment.
	if instance.Spec.DryRun != true {
		foundFederation := &api.Federation{}
		err = r.client.Get(context.TODO(), types.NamespacedName{
			Name:      instance.Spec.FederatorName,
			Namespace: instance.Spec.Benchmark}, foundFederation)
		if err != nil {
			return reconcile.Result{}, err
		}

		if ok, err := r.reconcileEvaluatorJob(instance, foundFederation); ok != false {
			return reconcile.Result{RequeueAfter: 5000000000}, err
		}
	}
	reqLogger.Info("Reached the end of the reconciling loop for the kobe Experiment %s/%s\n", instance.Name, instance.Namespace)
	return reconcile.Result{}, nil
}

func (r *ReconcileExperiment) reconcileEvaluatorJob(instance *api.Experiment, fed *api.Federation) (bool, error) {
	reqLogger := log
	fedEndpoint :=
		util.EndpointURL(fed.Name, fed.Namespace, int(fed.Spec.Template.Port), fed.Spec.Template.Path)
	fedName := fed.Name

	// Create the new job that will run the EVAL client for this experiment

	foundJob := &batchv1.Job{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: instance.Namespace,
		Name:      instance.Name + "-evaluationjob"},
		foundJob)

	if err == nil {
		if &foundJob.Status.Succeeded == nil || foundJob.Status.Succeeded == 0 {
			return true, nil
		}
		reqLogger.Info("The job is done\n")
		return false, nil
	}
	experimentJob := r.createEvaluatorJob(instance, fedEndpoint, fedName)
	reqLogger.Info("Creating a new job to run the experiment for this setup")
	if err := r.client.Create(context.TODO(), experimentJob); err != nil {
		reqLogger.Info("FAILED to create the job to run this experiment  %s/%s\n", experimentJob.Name, experimentJob.Namespace)
		return true, err
	}
	//instance.Status.CurrentRun = instance.Status.CurrentRun + 1
	// err = r.client.Status().Update(context.TODO(), instance)
	// if err != nil {
	// 	reqLogger.Info("Failed to update the times to run of the experiment")
	// 	return err
	// }

	return true, nil
}

//----------------------functions that create native kubernetes objects--------------------------------------
//create the job that will run the evaluation program
func (r *ReconcileExperiment) createEvaluatorJob(m *api.Experiment, fedendpoint string, fedname string) *batchv1.Job {
	times := int32(1)
	parallelism := int32(1) //hardcoded cause if not set it defaults to 0 and no pod is ever created/need defaulting asap
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-evaluationjob",
			Namespace: m.Namespace,
		},
		Spec: batchv1.JobSpec{
			Parallelism: &parallelism,
			Completions: &times,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           m.Spec.Evaluator.Image, //this is the image of the eval
						ImagePullPolicy: corev1.PullAlways,
						Name:            "job",
						//Command:         m.Spec.Evaluator.Command,
						Ports: []corev1.ContainerPort{{
							ContainerPort: int32(8890), //eval endpoint
							Name:          "client",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "queries",
							MountPath: "/queries",
						}},
						Env: append([]corev1.EnvVar{
							{Name: "FEDERATION_NAME", Value: fedname},
							{Name: "FEDERATION_ENDPOINT", Value: fedendpoint},
							{Name: "ENDPOINT", Value: fedendpoint},
							{Name: "EXPERIMENT", Value: m.Name},
							{Name: "EVAL_RUNS", Value: strconv.Itoa(m.Spec.TimesToRun)},
						}, m.Spec.Evaluator.Env...),
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

func (r *ReconcileExperiment) reconcileFederation(instance *api.Experiment) (bool, error) {
	reqLogger := log

	foundFederation := &api.Federation{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      instance.Spec.FederatorName,
		Namespace: instance.Spec.Benchmark}, foundFederation)

	if err != nil && errors.IsNotFound(err) {
		foundBenchmark := &api.Benchmark{}
		err = r.client.Get(context.TODO(), types.NamespacedName{
			Name:      instance.Spec.Benchmark,
			Namespace: instance.Namespace}, foundBenchmark)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Did not found a kobebenchmark resource with this name please define that first")
			return true, err
		}
		if instance.Spec.FederatorSpec == nil {
			foundTemplate := &api.FederatorTemplate{}
			reqLogger.Info("Finding the federator template reference specified for the experiment " + instance.Name + " %v\n")
			err := r.client.Get(context.TODO(), types.NamespacedName{
				Name:      instance.Spec.FederatorTemplateRef,
				Namespace: corev1.NamespaceDefault},
				foundTemplate)
			if err != nil && errors.IsNotFound(err) {
				reqLogger.Info("Failed to find the requested dataset template: ", err)
				return true, err
			}
			if err != nil {
				reqLogger.Info("this is the true error ", err)
				return true, err
			}
			instance.Spec.FederatorSpec = &foundTemplate.Spec
			err = r.client.Update(context.TODO(), instance)
			if err != nil {
				reqLogger.Info("failed to update the template spec of the federation: %v", instance.Spec.FederatorName)
				return true, err
			}
			return true, nil
		}

		newFederation := r.newFederation(instance, foundBenchmark)
		reqLogger.Info("Creating a new federation based on this experiments datasets and federator")
		err = r.client.Create(context.TODO(), newFederation)
		if err != nil {
			reqLogger.Info("Failed to create the federation")
			return false, err
		}
		newFederation.Status.Phase = api.FederationInitializing
		newFederation.Status.PodNames = []string{}
		err = r.client.Status().Update(context.TODO(), newFederation)
		if err != nil {
			reqLogger.Info("Failed to update the federation")
			return false, err
		}
	}
	if err != nil {
		return true, err
	}

	podList := &corev1.PodList{}
	listOps := []client.ListOption{
		client.InNamespace(instance.Spec.Benchmark),
		client.MatchingLabels{"kobeoperator_cr": instance.Spec.FederatorName},
	}
	err = r.client.List(context.TODO(), podList, listOps...)
	if err != nil {
		reqLogger.Info("Failed to list pods: %v", err)
		return true, err
	}

	podNames := getPodNames(podList.Items)
	//for _, podname := range foundFederation.Status.PodNames {
	for _, podname := range podNames {
		foundPod := &corev1.Pod{}
		err := r.client.Get(context.TODO(), types.NamespacedName{
			Namespace: instance.Spec.Benchmark,
			Name:      podname}, foundPod)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Failed to get the pod of the kobe federation that experiment will use")
			return true, nil
		}
		if foundPod.Status.Phase != corev1.PodRunning {
			reqLogger.Info("Kobe federation pod is not ready so experiment needs to wait")
			return true, nil
		}
	}
	if podNames == nil || len(podNames) == 0 {
		reqLogger.Info("Experiment waits for FEDERATOR initialization")
		return true, nil
	}
	return false, nil
}

//function that creates a new kobefederation custom resource from the federator and benchmark  in experiment.
//The native objects that kobefederation needs are created by kobefederation controller .
func (r *ReconcileExperiment) newFederation(m *api.Experiment, benchmark *api.Benchmark) *api.Federation {
	networktopology := []api.NetworkConnection{}
	datasetendpoints := []api.DatasetEndpoint{}

	for _, d := range benchmark.Spec.Datasets {
		//we need to find each dataset cause in benchmark it is possible that they dont carry the definition for their spec
		//but just the reference to a template.
		foundDataset := &api.EphemeralDataset{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: d.Name, Namespace: m.Spec.Benchmark}, foundDataset)
		if err != nil && errors.IsNotFound(err) {
			return nil
		}
		datasetendpoints = append(datasetendpoints, api.DatasetEndpoint{
			Host:      foundDataset.Name,
			Namespace: benchmark.Name,
			Port:      foundDataset.Spec.SystemSpec.Port,
			Path:      foundDataset.Spec.SystemSpec.Path})

		if d.FederatorConnection != nil {
			networktopology = append(networktopology, api.NetworkConnection{Source: &d.Name, DelayInjection: d.FederatorConnection.DelayInjection})
		}
	}

	federation := &api.Federation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Spec.FederatorName,
			Namespace: benchmark.Name,
		},
		Spec: api.FederationSpec{
			Template:        *m.Spec.FederatorSpec,
			InitPolicy:      api.ForceInit, //???????????????
			FederatorName:   m.Spec.FederatorName,
			Datasets:        datasetendpoints,
			NetworkTopology: networktopology,
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

func (r *ReconcileExperiment) reconcileDatasets(instance *api.Experiment, benchmark api.Benchmark) (bool, error) {
	reqLogger := log

	// Check if every kobedataset of the benchmark is healthy.
	// Create a list of the endpoints and of the names of the datasets
	for _, datasetInfo := range benchmark.Spec.Datasets {
		foundDataset := &api.EphemeralDataset{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: benchmark.Name, Name: datasetInfo.Name}, foundDataset)
		if err != nil {
			reqLogger.Info("Failed to find a specific dataset from the list of datasets of this benchmark")
			return true, err
		}

		// Check for the healthiness of the individual pods of the kobe dataset
		podList := &corev1.PodList{}
		listOps := []client.ListOption{
			client.InNamespace(benchmark.Name),
			client.MatchingLabels{"datasetName": foundDataset.Name},
		}
		err = r.client.List(context.TODO(), podList, listOps...)
		if err != nil {
			reqLogger.Info("Failed to list pods: %v", err)
			return true, err
		}

		podNames := getPodNames(podList.Items)
		for _, podname := range podNames {
			foundPod := &corev1.Pod{}
			err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: instance.Spec.Benchmark, Name: podname}, foundPod)
			if err != nil && errors.IsNotFound(err) {
				reqLogger.Info("Failed to get the pod of the kobe dataset that experiment will use")
				return true, nil
			}
			if foundPod.Status.Phase != corev1.PodRunning {
				reqLogger.Info("Kobe dataset pod is not ready so experiment needs to wait")
				return true, nil
			}
		}
		if podNames == nil || len(podNames) == 0 {
			reqLogger.Info("Experiment waits for components initialization")
			return true, nil
		}

	}
	return false, nil

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
