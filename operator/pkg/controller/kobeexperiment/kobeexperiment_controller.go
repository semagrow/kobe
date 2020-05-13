package kobeexperiment

import (
	"context"
	"strconv"

	kobebenchmarkv1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobebenchmark/v1alpha1"
	kobedatasetv1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobedataset/v1alpha1"
	kobeexperimentv1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobeexperiment/v1alpha1"
	kobefederationv1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobefederation/v1alpha1"
	kobefederatorv1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobefederator/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
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
var identifier = 0

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

	err = c.Watch(&source.Kind{Type: &batchv1.Job{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobeexperimentv1alpha1.KobeExperiment{},
	})
	if err != nil {
		return err
	}

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
	//normally i have to check for finishing initialization not just if they exist.Federator for example could be initiliazing with its init container .very important

	//check if there exist a kobe benchmark with this name in kubernetes.If not its an error .
	foundBenchmark := &kobebenchmarkv1alpha1.KobeBenchmark{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.Benchmark, Namespace: instance.Namespace}, foundBenchmark)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Did not found a kobebenchmark resource with this name please define that first")
		return reconcile.Result{RequeueAfter: 5}, err
	}
	endpoints := []string{}
	datasets := []string{}

	//check if every kobedataset of the benchmark is healthy.Create a list of the endpoints and of the names of the datasets
	for _, datasetInfo := range foundBenchmark.Spec.Datasets {
		foundDataset := &kobedatasetv1alpha1.KobeDataset{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: foundBenchmark.Namespace, Name: datasetInfo.Name}, foundDataset)
		if err != nil {
			reqLogger.Info("Failed to find a specific dataset from the list of datasets of this benchmark")
			return reconcile.Result{RequeueAfter: 5}, err
		}
		//check for the healthiness of the individual pods of the kobe dataset

		podList := &corev1.PodList{}
		labelSelector := labels.SelectorFromSet(map[string]string{"kobeoperator_cr": foundDataset.Name})
		listOps := &client.ListOptions{Namespace: instance.Namespace, LabelSelector: labelSelector}
		err = r.client.List(context.TODO(), listOps, podList)
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

		//create a list of the sparql endpoints
		endpoints = append(endpoints, "http://"+foundDataset.Name+"."+foundDataset.Namespace+".svc.cluster.local"+":"+strconv.Itoa(int(foundDataset.Spec.Port))+foundDataset.Spec.SparqlEnding)
		datasets = append(datasets, foundDataset.Name)
	}

	foundFederator := &kobefederatorv1alpha1.KobeFederator{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.Federator, Namespace: instance.Namespace}, foundFederator)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("No federator with this name is defined in the cluster")
		return reconcile.Result{}, err
	}
	if err != nil {
		reqLogger.Info("Error at getting this federator resource from the cluster.")
		return reconcile.Result{}, err
	}

	foundFederation := &kobefederationv1alpha1.KobeFederation{}
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

	//check if the pods of the federators exist and have a status of running before proceeding and get fed name and endpoint for the eval job
	fedEndpoint := "http://" + foundFederation.Name + "." + foundFederation.Namespace + ".svc.cluster.local" + ":" + strconv.Itoa(int(foundFederation.Spec.Port)) + foundFederation.Spec.SparqlEnding
	fedName := foundFederation.Name

	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(map[string]string{"kobeoperator_cr": instance.Name})
	listOps := &client.ListOptions{Namespace: instance.Namespace, LabelSelector: labelSelector}
	err = r.client.List(context.TODO(), listOps, podList)
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

	//Everything is healthy and ready for the experiment.
	if instance.Spec.DryRun == true { //dont run just yet just have it defined
		return reconcile.Result{}, nil
	}

	//Create the new job that will run the EVAL client for this experiment
	if instance.Spec.TimesToRun > 0 {
		foundJob := &batchv1.Job{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name + "-" + strconv.Itoa(identifier)}, foundJob)
		if err == nil {
			if &foundJob.Status.Succeeded == nil || foundJob.Status.Succeeded == 0 {
				return reconcile.Result{}, nil
			}
			reqLogger.Info("All past jobs are done /n")
			identifier++
		}
		experimentJob := r.newJobForExperiment(instance, identifier, fedEndpoint, fedName)
		reqLogger.Info("Creating a new job to run the experiment for this setup")
		err = r.client.Create(context.TODO(), experimentJob)
		if err != nil {
			reqLogger.Info("FAILED to create the job to run this expriment  %s/%s\n", experimentJob.Name, experimentJob.Namespace)
			return reconcile.Result{}, err
		}
		instance.Spec.TimesToRun = 0 //instance.Spec.TimesToRun - 1
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
func (r *ReconcileKobeExperiment) newJobForExperiment(m *kobeexperimentv1alpha1.KobeExperiment, i int, fedendpoint string, fedname string) *batchv1.Job {
	times := int32(1)
	parallelism := int32(1)
	//labels := map[string]string{"name": m.Name}
	envs := []corev1.EnvVar{}
	env := corev1.EnvVar{Name: "FEDERATION_NAME", Value: fedname}
	envs = append(envs, env)
	env = corev1.EnvVar{Name: "FEDERATION_ENDPOINT", Value: fedendpoint}
	envs = append(envs, env)
	env = corev1.EnvVar{Name: "ENDPOINT", Value: fedendpoint}
	envs = append(envs, env)
  env = corev1.EnvVar{Name: "EXPERIMENT", Value: m.Name}
	envs = append(envs, env)
	env = corev1.EnvVar{Name: "EVAL_RUNS", Value: strconv.Itoa(m.Spec.TimesToRun)}
	envs = append(envs, env)

	volumes := []corev1.Volume{corev1.Volume{Name: "queries",
		VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: m.Spec.Benchmark}}}}}
	vmounts := []corev1.VolumeMount{corev1.VolumeMount{Name: "queries", MountPath: "/queries"}}

	job := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-" + strconv.Itoa(i),
			Namespace: m.Namespace,
		},
		Spec: batchv1.JobSpec{
			Parallelism: &parallelism,
			Completions: &times,
			Template: corev1.PodTemplateSpec{
				metav1.ObjectMeta{},
				corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           m.Spec.EvalImage, //this is the image of the eval
						Name:            "job" + "-" + strconv.Itoa(i),
						ImagePullPolicy: corev1.PullAlways,
						Ports: []corev1.ContainerPort{{
							ContainerPort: int32(8890), //eval endpoint
							Name:          "client",
						}},
						//Command:      m.Spec.EvalCommands,
						VolumeMounts: vmounts,
						Env:          envs,
					}},
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Volumes:       volumes,
				},
			},
		},
	}
	controllerutil.SetControllerReference(m, job, r.scheme)
	return job

}

//function that creates a new kobefederation custom resource from the federator and benchmark  in kobeexperiment.
//The native objects that kobefederation needs are created by kobefederation controller .
func (r *ReconcileKobeExperiment) newFederationForExperiment(m *kobeexperimentv1alpha1.KobeExperiment,
	fed *kobefederatorv1alpha1.KobeFederator, endpoints []string, datasetnames []string) *kobefederationv1alpha1.KobeFederation {

	federation := &kobefederationv1alpha1.KobeFederation{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kobefederator.kobe.com/v1alpha1",
			Kind:       "KobeFederator",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: fed.Namespace,
		},
		Spec: kobefederationv1alpha1.KobeFederationSpec{
			Image:             fed.Spec.Image,
			ImagePullPolicy:   fed.Spec.ImagePullPolicy,
			Affinity:          fed.Spec.Affinity,
			Port:              fed.Spec.Port,
			ConfFromFileImage: fed.Spec.ConfFromFileImage,
			InputDumpDir:      fed.Spec.InputDumpDir,
			OutputDumpDir:     fed.Spec.OutputDumpDir,
			ConfImage:         fed.Spec.ConfImage,
			InputDir:          fed.Spec.InputDir,
			OutputDir:         fed.Spec.OutputDir,
			FedConfDir:        fed.Spec.FedConfDir,
			SparqlEnding:      fed.Spec.SparqlEnding,

			ForceNewInit:  m.Spec.ForceNewInit,
			Init:          true,
			FederatorName: fed.Name,
			Endpoints:     endpoints,
			DatasetNames:  datasetnames,
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
