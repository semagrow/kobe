package federation

import (
	"context"
	"reflect"
	"strconv"

	api "github.com/semagrow/kobe/operator/pkg/apis/kobe/v1alpha1"
	"github.com/semagrow/kobe/operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_federation")

// Add creates a new Federation Controller and adds it to the Manager. The Manager will set fields on the Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileFederation{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("federation-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource KobeFederator
	err = c.Watch(&source.Kind{Type: &api.Federation{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner KobeFederator
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &api.Federation{},
	})

	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &api.Federation{},
	})

	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner KobeDataset
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &api.Federation{},
	})

	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &batchv1.Job{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &api.Federation{},
	})
	if err != nil {
		return err
	}
	return nil
}

// blank assignment to verify that ReconcileFederation implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileFederation{}

// ReconcileFederation reconciles a Federation object
type ReconcileFederation struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Federation object and
// makes changes based on the state read and what is in the Federation.Spec
//
// Note:
// The Controller will requeue the Request to be processed again if the returned
// error is non-nil or Result.Requeue is true, otherwise upon completion it will
// remove the work from the queue.
func (r *ReconcileFederation) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Federation")

	// fetch the Federation instance
	instance := &api.Federation{}
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
	//setting fields internally to default if not present
	instance.SetDefaults()

	endpoints := []string{}
	datasets := []string{}

	// Check if every kobedataset of the benchmark is healthy.
	// Create a list of the endpoints and of the names of the datasets
	for _, datasetInfo := range instance.Spec.Datasets {
		foundDataset := &api.Dataset{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: instance.Namespace, Name: datasetInfo}, foundDataset)
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
			if foundPod.Status.Phase != corev1.PodRunning {
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
			util.EndpointURL(foundDataset.Name, foundDataset.Namespace, int(foundDataset.Spec.Port), foundDataset.Spec.Path))
		datasets = append(datasets, foundDataset.Name)
	}

	datasetsForInit := []string{}  //here we will collect only datasets that get init containers for metadata creation
	endpointsForInit := []string{} //here we will collect the endpoints that correspond to the selected datasets in the above slice

	// getting plan for metadata creation
	if instance.Status.Phase == api.FederationInitializing {
		// the federation controller still runs the init loop as long as this
		// flag is true

		// create a job that will make the necessary directories to save the
		// config files for future caching ( in dataset-name/federator/ for all
		// datasets)
		foundJob := &batchv1.Job{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundJob)
		if err != nil && errors.IsNotFound(err) {
			job := r.newJobForFederation(instance)
			err = r.client.Create(context.TODO(), job)
			if err != nil {
				reqLogger.Info("Failed to create the init job that will make the directories in the server for caching")
				return reconcile.Result{}, err
			}
		} else if err != nil {
			reqLogger.Info("Failed to retrieve the job that makes the directories")
			return reconcile.Result{}, err
		}

		//hang till it finishes successfully (this controller listens to job
		//changes so he will awake if the job status changes /no need to
		//requeue)
		if &foundJob.Status.Succeeded == nil || foundJob.Status.Succeeded == 0 {
			return reconcile.Result{}, nil
		}

		//----------------------experimental jobs-------------------------------------
		//create jobs for the federation datasets that will check if those
		//datasets have init files for this federator already by either failing
		//or succeeding
		for _, dataset := range datasets {
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: dataset, Namespace: instance.Namespace}, foundJob)
			if err != nil && errors.IsNotFound(err) {
				job := r.newJobForDataset(instance, dataset)
				err = r.client.Create(context.TODO(), job)
				if err != nil {
					reqLogger.Info("Failed to create the dataset job that checks if config file already exists")
					return reconcile.Result{}, err
				}
			} else if err != nil {
				reqLogger.Info("Failed to retrieve the job ")
				return reconcile.Result{}, err
			}
		}

		// wait till they all finish either with error or successfully and
		// collect a list with those that errored ->which means they didn't find
		// init files if forcenewinit is true then the list contains all the
		// datasets since we will initialize for all of them again if
		// forcenewinit is false only those that errored will get passed to the
		// list to make init containers
		for i, dataset := range datasets { //loop through all datasets of this federation
			foundJob := &batchv1.Job{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: dataset, Namespace: instance.Namespace}, foundJob)
			if err != nil && errors.IsNotFound(err) {
				return reconcile.Result{Requeue: true}, nil
			} else if err != nil { //some other error
				return reconcile.Result{RequeueAfter: 10}, err
			}
			//fetch the pod of the init - job for this dataset to check its status
			podList := &corev1.PodList{}
			listOps := []client.ListOption{
				client.InNamespace(instance.Namespace),
				client.MatchingLabels{"job-name": dataset},
			}
			err = r.client.List(context.TODO(), podList, listOps...)
			if err != nil {
				reqLogger.Info("Failed to list pods: %v", err)
				return reconcile.Result{}, err
			}
			//if the job-pod doesn't exist yet then requeue (we got here faster than we should and must wait)
			podNames := getPodNames(podList.Items)
			if podNames == nil || len(podNames) == 0 {
				return reconcile.Result{RequeueAfter: 15}, nil

			}
			pod := &corev1.Pod{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: podNames[0], Namespace: instance.Namespace}, pod) //fetch the pod
			if err != nil {
				reqLogger.Info("Failed to get the pod that checks if config file for dataset exist")
				return reconcile.Result{}, err
			}
			//decide whether to include this dataset in the initialization based on the status of the job pod  and the forceNewInit flag
			if instance.Spec.InitPolicy != api.ForceInit { //we make a choice
				if pod.Status.Phase == corev1.PodRunning {

				} else if pod.Status.Phase == corev1.PodFailed {
					datasetsForInit = append(datasetsForInit, dataset)
					endpointsForInit = append(endpointsForInit, endpoints[i])
				} else { //pod is still running so we again need to wait for it before seeing if it failed or succeededs
					return reconcile.Result{RequeueAfter: 5}, nil
				}
			} else if instance.Spec.InitPolicy == api.ForceInit { //we dont make a choice we gather all of them
				datasetsForInit = append(datasetsForInit, dataset)
				endpointsForInit = append(endpointsForInit, endpoints[i])
			}
		}
		//clean up the jobs that checked for the files
		for _, dataset := range instance.Spec.Datasets {
			foundJob := &batchv1.Job{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: dataset, Namespace: instance.Namespace}, foundJob)
			err = r.client.Delete(context.TODO(), foundJob, client.PropagationPolicy(metav1.DeletionPropagation("Background")))
		}

		//------------------------------/experimental jobs------------------------------------------------------------

		//clean up the job that made the necessary directories to safe keep the init files
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundJob)
		err = r.client.Delete(context.TODO(), foundJob, client.PropagationPolicy(metav1.DeletionPropagation("Background")))
		if err != nil {
			reqLogger.Info("Failed to delete the federation job from the cluster")
			return reconcile.Result{}, err
		}

		// Never rerun the init jobs (this whole part of the loop) even if the
		// user changes an attribute of the federation object unless he redefines
		// the experiment if this flag change doesn't happen,then every time this
		// controller reruns to reconcile our federation we will get a repeat of
		// all the init process of the federation jobs again and again. Also if
		// federation pod drops and this controller relaunches it ,it will not
		// recreate the init files per dataset since datasetsToInit will be empty
		// which means we save time.
		instance.Status.Phase = api.FederationRunning
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Info("Failed to update the init flag")
			return reconcile.Result{}, err
		}
	}

	// NOTE: We currently use a Pod instead of a Deployment to avoid the respawning of
	// the Pod (and therefore re-execute the initContainers)
	// check for the healthiness of the federation pod and create it if it
	// doesn't exist
	foundPod := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundPod)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Making a new pod for federation", "instance.Namespace", instance.Namespace, "instance.Name", instance.Name)
		pod := r.newPodForFederation(instance, datasetsForInit, endpointsForInit)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			reqLogger.Info("Failed to create new Pod: %v\n", err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// check for status changes
	podList := &corev1.PodList{}
	listOps := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels(labelsForFederation(instance.Name)),
	}
	err = r.client.List(context.TODO(), podList, listOps...)
	if err != nil {
		reqLogger.Info("Failed to list pods: %v", err)
		return reconcile.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.PodNames if needed
	if !reflect.DeepEqual(podNames, instance.Status.PodNames) {
		instance.Status.PodNames = podNames
		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Info("failed to update node status: %v", err)
			return reconcile.Result{}, err
		}
	}

	//check the healthiness of the federation service that is used for name resolving
	foundService := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Making a new service for the federation", "instance.Namespace", instance.Namespace, "instance.Name", instance.Name)
		service := r.newServiceForFederation(instance)
		reqLogger.Info("Creating a new Service %s/%s\n", service.Namespace, service.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			reqLogger.Info("Failed to create new Service: %v\n", err)
			return reconcile.Result{}, err
		}

		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	//all checks are completed successfully
	reqLogger.Info("Loop went through the end for reconciling kobedataset")

	return reconcile.Result{}, nil
}

//---------------------------------functions that create native kubernetes objects that are owned by a federation -----------------------
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

func labelsForFederation(name string) map[string]string {
	return map[string]string{"app": "Kobe-Operator", "kobeoperator_cr": name}
}

// a service to find the federation by name internally in the cluster.
func (r *ReconcileFederation) newServiceForFederation(m *api.Federation) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},

		Spec: corev1.ServiceSpec{
			Selector: labelsForFederation(m.Name),
			Ports: []corev1.ServicePort{
				{
					Port: m.Spec.Template.Port,
					TargetPort: intstr.IntOrString{
						IntVal: m.Spec.Template.Port,
					},
				},
			},
		},
	}
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// A job that remove the (temporary) files and create some dirs...
func (r *ReconcileFederation) newJobForFederation(m *api.Federation) *batchv1.Job {
	times := int32(1)
	parallelism := int32(1)
	volumes := []corev1.Volume{}
	vmounts := []corev1.VolumeMount{}

	nfsPodFound := &corev1.Pod{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: "kobenfs", Namespace: m.Namespace}, nfsPodFound)
	if err != nil && errors.IsNotFound(err) {
		return nil
	}
	nfsip := nfsPodFound.Status.PodIP

	volume := corev1.Volume{
		Name: "nfs-job",
		VolumeSource: corev1.VolumeSource{
			NFS: &corev1.NFSVolumeSource{
				Server: nfsip,
				Path:   "/exports/"}}}

	vmountFinal := corev1.VolumeMount{
		Name:      "nfs-job",
		MountPath: "/kobe/"}

	volumes = append(volumes, volume)
	vmounts = append(vmounts, vmountFinal)

	job := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: batchv1.JobSpec{
			Parallelism: &parallelism,
			Completions: &times,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           "busybox",
						Name:            m.Name,
						ImagePullPolicy: corev1.PullIfNotPresent,
						VolumeMounts:    vmounts,
						Command:         []string{"sh", "-c"},
						Args: []string{
							"cd /kobe ; rm -r " + m.Name + " ; rm -r" + " temp-" + m.Name + " ; for d in */; do   cd $d;  mkdir " + m.Spec.FederatorName +
								"; cd /kobe ; done ;" + " mkdir " + m.Name + " ; mkdir " + "temp-" + m.Name},
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

//------------------------ job that checks if init file exists for this dataset/federator by failing or succeeding
func (r *ReconcileFederation) newJobForDataset(m *api.Federation, dataset string) *batchv1.Job {
	times := int32(1)
	parallelism := int32(1)
	volumes := []corev1.Volume{}
	vmounts := []corev1.VolumeMount{}

	nfsPodFound := &corev1.Pod{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: "kobenfs", Namespace: m.Namespace}, nfsPodFound)
	if err != nil && errors.IsNotFound(err) {
		return nil
	}
	nfsip := nfsPodFound.Status.PodIP

	volume := corev1.Volume{Name: "nfs-job", VolumeSource: corev1.VolumeSource{NFS: &corev1.NFSVolumeSource{Server: nfsip, Path: "/exports/"}}}
	volumes = append(volumes, volume)

	vmountFinal := corev1.VolumeMount{Name: "nfs-job", MountPath: "/kobe/"}
	vmounts = append(vmounts, vmountFinal)

	job := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      dataset,
			Namespace: m.Namespace,
		},
		Spec: batchv1.JobSpec{
			Parallelism: &parallelism,
			Completions: &times,
			Template: corev1.PodTemplateSpec{
				metav1.ObjectMeta{},
				corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           "busybox",
						Name:            m.Name,
						ImagePullPolicy: corev1.PullIfNotPresent,
						VolumeMounts:    vmounts,
						Command:         []string{"sh", "-c"},
						Args:            []string{"cat /kobe/" + dataset + "/" + m.Spec.FederatorName + "/*"},
					}},
					RestartPolicy: corev1.RestartPolicyNever,
					Volumes:       volumes,
				},
			},
		},
	}
	controllerutil.SetControllerReference(m, job, r.scheme)
	return job

}

// creates a new federation deployment
// This is the deployment that runs the federator image
func (r *ReconcileFederation) newPodForFederation(m *api.Federation, datasets []string, endpoints []string) *corev1.Pod {
	labels := labelsForFederation(m.Name)

	nfsPodFound := &corev1.Pod{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: "kobenfs", Namespace: m.Namespace}, nfsPodFound)
	if err != nil && errors.IsNotFound(err) {
		return nil
	}
	nfsip := nfsPodFound.Status.PodIP //it seems we need this cause dns for service of the nfs doesnt work in kubernetes

	//create init containers  that make one config file for federation per dataset dump /if not needed then these can be set to do nothing
	initContainers := []corev1.Container{}
	volumes := []corev1.Volume{}
	for i, datasetname := range datasets {
		//each init container is given DATASET_NAME and DATASET_ENDPOINT environment variables to work with(needed if they create the files from quering the database directly)
		//also inputfiledir and outputfiledir both point to exports/<datasetname>/dumps/ and exports/dataset/<datasetname>/<federation>/ respectively to nfs server(needed if they make the config files from the dumps)
		vmounts := []corev1.VolumeMount{}
		envs := []corev1.EnvVar{
			{Name: "DATASET_NAME", Value: datasetname},
			{Name: "DATASET_ENDPOINT", Value: endpoints[i]}}

		if m.Spec.InitPolicy == api.ForceInit { //optional variable to skip creating the files if they already exist in /exports/<dataset-name>/<federator-name>.Is passed by the experiment yaml
			env := corev1.EnvVar{Name: "INITIALIZE", Value: "yes"}
			envs = append(envs, env)
		}
		volumeIn := corev1.Volume{Name: "nfs-in-" + datasetname, VolumeSource: corev1.VolumeSource{NFS: &corev1.NFSVolumeSource{Server: nfsip, Path: "/exports/" + datasetname + "/dump"}}}
		volumeOut := corev1.Volume{Name: "nfs-out-" + datasetname, VolumeSource: corev1.VolumeSource{NFS: &corev1.NFSVolumeSource{Server: nfsip, Path: "/exports/" + datasetname + "/" + m.Spec.FederatorName}}}
		volumes = append(volumes, volumeIn, volumeOut)

		vmountIn := corev1.VolumeMount{Name: "nfs-in-" + datasetname, MountPath: m.Spec.Template.InputDumpDir}
		vmountOut := corev1.VolumeMount{Name: "nfs-out-" + datasetname, MountPath: m.Spec.Template.OutputDumpDir}
		vmounts = append(vmounts, vmountIn, vmountOut)

		container := corev1.Container{
			Image:        m.Spec.Template.ConfFromFileImage,
			Name:         "initcontainer" + strconv.Itoa(i),
			Env:          envs,
			VolumeMounts: vmounts,
		}
		initContainers = append(initContainers, container)
	}
	//create a helper init container that will choose the config files for this set of datasets only and move them in a temps directory
	envs := []corev1.EnvVar{}
	vmounts := []corev1.VolumeMount{}
	count := 0
	for i, datasetname := range m.Spec.Datasets {
		env := corev1.EnvVar{Name: "DATASET_NAME_" + strconv.Itoa(i), Value: datasetname}
		envs = append(envs, env)
		env = corev1.EnvVar{Name: "DATASET_ENDPOINT_" + strconv.Itoa(i), Value: endpoints[i]}
		envs = append(envs, env)
		count++
	}
	env := corev1.EnvVar{Name: "N", Value: strconv.Itoa(count - 1)}
	envs = append(envs, env)

	env = corev1.EnvVar{Name: "FEDERATION_NAME", Value: m.Name}
	envs = append(envs, env)

	env = corev1.EnvVar{Name: "FEDERATOR_NAME", Value: m.Spec.FederatorName}
	envs = append(envs, env)

	volumeHouse := corev1.Volume{Name: "nfs-housekeep", VolumeSource: corev1.VolumeSource{NFS: &corev1.NFSVolumeSource{Server: nfsip, Path: "/exports"}}}
	volumes = append(volumes, volumeHouse)

	vmountHouse := corev1.VolumeMount{Name: "nfs-housekeep", MountPath: "/kobe"}
	vmounts = append(vmounts, vmountHouse)

	containerHouse := corev1.Container{
		Image:        "kostbabis/housekeeping",
		Name:         "inithouse",
		Env:          envs,
		VolumeMounts: vmounts,
	}
	initContainers = append(initContainers, containerHouse)

	//create the initcontainer that will run the image that combines many configs from the above temp directory and make appropriate config for the whole experiment/federation
	vmounts = []corev1.VolumeMount{}
	path := "/exports/temp-" + m.Name

	volumeInFinal := corev1.Volume{Name: "nfs-final-in", VolumeSource: corev1.VolumeSource{NFS: &corev1.NFSVolumeSource{Server: nfsip, Path: path}}}
	volumes = append(volumes, volumeInFinal)

	vmountInFinal := corev1.VolumeMount{Name: "nfs-final-in", MountPath: m.Spec.Template.InputDir}
	vmounts = append(vmounts, vmountInFinal)

	path = "/exports/" + m.Name

	volumeOutFinal := corev1.Volume{Name: "nfs-final-out", VolumeSource: corev1.VolumeSource{NFS: &corev1.NFSVolumeSource{Server: nfsip, Path: path}}}
	volumes = append(volumes, volumeOutFinal)

	vmountOutFinal := corev1.VolumeMount{Name: "nfs-final-out", MountPath: m.Spec.Template.OutputDir}
	vmounts = append(vmounts, vmountOutFinal)

	container := corev1.Container{
		Image:        m.Spec.Template.ConfImage,
		Name:         "init" + "final",
		Env:          envs,
		VolumeMounts: vmounts,
	}
	initContainers = append(initContainers, container)

	//create the deployment of the federation .
	//mount the config files to where the federator needs (for example etc/default/semagrow) -->passed by the yaml of federator
	volumeConf := corev1.Volume{
		Name: "volumeconf",
		VolumeSource: corev1.VolumeSource{
			NFS: &corev1.NFSVolumeSource{
				Server: nfsip,
				Path:   "/exports/" + m.Name + "/"},
		}}
	volumes = append(volumes, volumeConf)

	mountConf := corev1.VolumeMount{
		Name:      "volumeconf",
		MountPath: m.Spec.Template.FedConfDir}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{

			InitContainers: initContainers,
			Containers: []corev1.Container{{
				Image:           m.Spec.Template.Image,
				Name:            m.Name,
				ImagePullPolicy: m.Spec.Template.ImagePullPolicy,
				Ports: []corev1.ContainerPort{{
					ContainerPort: m.Spec.Template.Port,
					Name:          m.Name,
				}},
				VolumeMounts: []corev1.VolumeMount{mountConf},
			}},
			Volumes: volumes,
		},
	}

	controllerutil.SetControllerReference(m, pod, r.scheme)
	return pod
}