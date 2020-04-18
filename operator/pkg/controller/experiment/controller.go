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

	requeue, err := r.reconcileFederation(instance)
	if requeue {
		return reconcile.Result{Requeue: true}, err
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Everything is healthy and ready for the experiment.
	if instance.Spec.DryRun != true {
		foundFederation := &api.Federation{}
		err = r.client.Get(context.TODO(), types.NamespacedName{
			Name:      instance.Name,
			Namespace: instance.Namespace}, foundFederation)
		if err != nil {
			return reconcile.Result{}, err
		}

		if err := r.reconcileEvaluatorJob(instance, foundFederation); err != nil {
			return reconcile.Result{}, err
		}
	}
	reqLogger.Info("Reached the end of the reconciling loop for the kobe Experiment %s/%s\n", instance.Name, instance.Namespace)
	return reconcile.Result{}, nil
}

func (r *ReconcileExperiment) reconcileEvaluatorJob(instance *api.Experiment, fed *api.Federation) error {
	reqLogger := log
	fedEndpoint :=
		util.EndpointURL(fed.Name, fed.Namespace, int(fed.Spec.Template.Port), fed.Spec.Template.Path)
	fedName := fed.Name

	// Create the new job that will run the EVAL client for this experiment
	if instance.Status.CurrentRun <= instance.Spec.TimesToRun {
		foundJob := &batchv1.Job{}
		err := r.client.Get(context.TODO(), types.NamespacedName{
			Namespace: instance.Namespace,
			Name:      instance.Name + "-" + strconv.Itoa(identifier)},
			foundJob)

		if err == nil {
			if &foundJob.Status.Succeeded == nil || foundJob.Status.Succeeded == 0 {
				return nil
			}
			reqLogger.Info("All past jobs are done\n")
			identifier++
		}
		experimentJob := r.createEvaluatorJob(instance, identifier, fedEndpoint, fedName)
		reqLogger.Info("Creating a new job to run the experiment for this setup")
		if err := r.client.Create(context.TODO(), experimentJob); err != nil {
			reqLogger.Info("FAILED to create the job to run this experiment  %s/%s\n", experimentJob.Name, experimentJob.Namespace)
			return err
		}
		instance.Status.CurrentRun = instance.Status.CurrentRun + 1
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Info("Failed to update the times to run of the experiment")
			return err
		}
	}
	return nil
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

func (r *ReconcileExperiment) reconcileFederation(instance *api.Experiment) (bool, error) {
	reqLogger := log

	foundFederation := &api.Federation{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      instance.Name,
		Namespace: instance.Namespace}, foundFederation)

	if err != nil && errors.IsNotFound(err) {
		foundBenchmark := &api.Benchmark{}
		err = r.client.Get(context.TODO(), types.NamespacedName{
			Name:      instance.Spec.Benchmark,
			Namespace: instance.Namespace}, foundBenchmark)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Did not found a kobebenchmark resource with this name please define that first")
			return true, err
		}

		foundFederator := &api.Federator{}
		err = r.client.Get(context.TODO(), types.NamespacedName{
			Name:      instance.Spec.Federator,
			Namespace: instance.Namespace}, foundFederator)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("No federator with this name is defined in the cluster")
			return false, err
		}
		if err != nil {
			reqLogger.Info("Error at getting this federator resource from the cluster.")
			return false, err
		}

		newFederation := r.newFederation(instance, foundFederator, foundBenchmark.Spec.Datasets)
		reqLogger.Info("Creating a new federation based on this experiments datasets and federator")
		err = r.client.Create(context.TODO(), newFederation)
		if err != nil {
			reqLogger.Info("Failed to create the federation")
			return false, err
		}
	}
	if err != nil {
		return false, err
	}

	podList := &corev1.PodList{}
	listOps := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels{"kobeoperator_cr": instance.Name},
	}
	err = r.client.List(context.TODO(), podList, listOps...)
	if err != nil {
		reqLogger.Info("Failed to list pods: %v", err)
		return false, err
	}
	podNames := getPodNames(podList.Items)
	//for _, podname := range foundFederation.Status.PodNames {
	for _, podname := range podNames {
		foundPod := &corev1.Pod{}
		err := r.client.Get(context.TODO(), types.NamespacedName{
			Namespace: instance.Namespace,
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
	return false, nil
}

//function that creates a new kobefederation custom resource from the federator and benchmark  in experiment.
//The native objects that kobefederation needs are created by kobefederation controller .
func (r *ReconcileExperiment) newFederation(m *api.Experiment,
	fed *api.Federator, datasets []string) *api.Federation {

	federation := &api.Federation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: fed.Namespace,
		},
		Spec: api.FederationSpec{
			Template:      fed.Spec.FederatorTemplate,
			InitPolicy:    api.ForceInit,
			FederatorName: fed.Name,
			Datasets:      datasets,
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
