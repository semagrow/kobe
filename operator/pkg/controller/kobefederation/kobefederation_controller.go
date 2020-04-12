package kobefederation

import (
	"context"
	"reflect"
	"strconv"

	kobev1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobe/v1alpha1"
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

var log = logf.Log.WithName("controller_kobefederation")

// Add creates a new KobeFederation Controller and adds it to the Manager. The Manager will set fields on the Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileKobeFederation{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("kobefederation-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource KobeFederator
	err = c.Watch(&source.Kind{Type: &kobev1alpha1.KobeFederation{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner KobeFederator
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobev1alpha1.KobeFederation{},
	})

	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobev1alpha1.KobeFederation{},
	})

	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner KobeDataset
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobev1alpha1.KobeFederation{},
	})

	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &batchv1.Job{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobev1alpha1.KobeFederation{},
	})
	if err != nil {
		return err
	}
	return nil
}

// blank assignment to verify that ReconcileKobeFederation implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileKobeFederation{}

// ReconcileKobeFederation reconciles a KobeFederation object
type ReconcileKobeFederation struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a KobeFederation object and
// makes changes based on the state read and what is in the KobeFederation.Spec
//
// Note:
// The Controller will requeue the Request to be processed again if the returned
// error is non-nil or Result.Requeue is true, otherwise upon completion it will
// remove the work from the queue.
func (r *ReconcileKobeFederation) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling KobeFederation")

	// fetch the KobeFederation instance
	instance := &kobev1alpha1.KobeFederation{}
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
	if instance.Spec.InputDumpDir == "" {
		instance.Spec.InputDumpDir = "/kobe/input"
	}
	if instance.Spec.OutputDumpDir == "" {
		instance.Spec.OutputDumpDir = "/kobe/output"
	}
	if instance.Spec.InputDir == "" {
		instance.Spec.InputDir = "/kobe/input"
	}
	if instance.Spec.OutputDir == "" {
		instance.Spec.OutputDir = "/kobe/output"
	}
	datasetsForInit := []string{}  //here we will collect only datasets that get init containers for metadata creation
	endpointsForInit := []string{} //here we will collect the endpoints that correspond to the selected datasets in the above slice

	// getting plan for metadata creation
	if instance.Spec.Init == true {
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
		for _, dataset := range instance.Spec.DatasetNames {
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
		// collect a list with those that errored ->which means they didnt find
		// init files if forcenewinit is true then the list contains all the
		// datasets since we will initialize for all of them again if
		// forcenewinit is false only those that errored will get passed to the
		// list to make init containers
		for i, dataset := range instance.Spec.DatasetNames { //loop through all datasets of this federation
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
			//if the job-pod doesnt exist yet then requeue (we got here faster than we should and must wait)
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
			//decide wether to include this dataset in the initialization based on the status of the job pod  and the forceNewInit flag
			if instance.Spec.ForceNewInit == false { //we make a choice
				if pod.Status.Phase == corev1.PodPhase("Succeeded") {

				} else if pod.Status.Phase == corev1.PodPhase("Failed") {
					datasetsForInit = append(datasetsForInit, dataset)
					endpointsForInit = append(endpointsForInit, instance.Spec.Endpoints[i])
				} else { //pod is still running so we again need to wait for it before seeing if it failed or succeededs
					return reconcile.Result{RequeueAfter: 5}, nil
				}
			} else if instance.Spec.ForceNewInit == true { //we dont make a choice we gather all of them
				datasetsForInit = append(datasetsForInit, dataset)
				endpointsForInit = append(endpointsForInit, instance.Spec.Endpoints[i])
			}
		}
		//clean up the jobs that checked for the files
		for _, dataset := range instance.Spec.DatasetNames {
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
		// the experiment if this flag change doesnt happen,then every time this
		// controller reruns to reconcile our federation we will get a repeat of
		// all the init process of the federation jobs again and again. Also if
		// federation pod drops and this controller relaunches it ,it will not
		// recreate the init files per dataset since datasetsToInit will be empty
		// which means we save time.
		instance.Spec.Init = false
		err = r.client.Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Info("Failed to update the init flag")
			return reconcile.Result{}, err
		}

	}

	// NOTE: We currently use a Pod instead of a Deployment to avoid the respawning of
	// the Pod (and therefore reexecute the initContainers)
	// check for the healthiness of the federation pod and create it if it
	// doesnt exist
	foundPod := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundPod)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Making a new pod for kobefederation", "instance.Namespace", instance.Namespace, "instance.Name", instance.Name)
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
		client.MatchingLabels(labelsForKobeFederation(instance.Name)),
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

func labelsForKobeFederation(name string) map[string]string {
	return map[string]string{"app": "Kobe-Operator", "kobeoperator_cr": name}
}

/*
func (r *ReconcileKobeFederation) newDeploymentForFederation(m *kobefederationv1alpha1.KobeFederation, datasets []string) *appsv1.Deployment {
	labels := labelsForKobeFederation(m.Name)

	// First, find the pod that NFS server is running
	nfsPodFound := &corev1.Pod{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: "kobenfs", Namespace: m.Namespace}, nfsPodFound)
	if err != nil && errors.IsNotFound(err) {
		return nil
	}

	// it seems we need this cause dns for service of the nfs doesn't work in Kubernetes
	nfsip := nfsPodFound.Status.PodIP

	// create init containers  that make one config file for federation per
	// dataset dump. if not needed then these can be set to do nothing
	initContainers := []corev1.Container{}
	volumes := []corev1.Volume{}
	for i, datasetname := range datasets {
		// each init container is given DATASET_NAME and DATASET_ENDPOINT
		// environment variables to work with (needed if they create the files
		// from quering the database directly) also `inputfiledir` and
		// `outputfiledir` both point to `exports/<datasetname>/dumps/` and
		// `exports/dataset/<datasetname>/<federation>/` respectively to nfs
		// server (needed if they make the config files from the dumps)
		vmounts := []corev1.VolumeMount{}
		envs := []corev1.EnvVar{}
		env := corev1.EnvVar{Name: "DATASET_NAME", Value: datasetname}
		envs = append(envs, env)

		env = corev1.EnvVar{Name: "DATASET_ENDPOINT", Value: m.Spec.Endpoints[i]}
		envs = append(envs, env)

		// optional variable to skip creating the files if they already exist in
		// `/exports/<dataset-name>/<federator-name>`. Is passed by the experiment
		// yaml
		if m.Spec.ForceNewInit == true {
			env = corev1.EnvVar{Name: "INITIALIZE", Value: "yes"}
			envs = append(envs, env)
		}

		volumeIn := corev1.Volume{
			Name: "nfs-in-" + datasetname,
			VolumeSource: corev1.VolumeSource{
				NFS: &corev1.NFSVolumeSource{
					Server: nfsip,
					Path:   "/exports/" + datasetname + "/dump"}}}

		vmountIn := corev1.VolumeMount{
			Name:      "nfs-in-" + datasetname,
			MountPath: m.Spec.InputDumpDir}

		volumeOut := corev1.Volume{
			Name: "nfs-out-" + datasetname,
			VolumeSource: corev1.VolumeSource{
				NFS: &corev1.NFSVolumeSource{
					Server: nfsip,
					Path:   "/exports/" + datasetname + "/" + m.Spec.FederatorName}}}

		vmountOut := corev1.VolumeMount{
			Name:      "nfs-out-" + datasetname,
			MountPath: m.Spec.OutputDumpDir}

		volumes = append(volumes, volumeIn, volumeOut)
		vmounts = append(vmounts, vmountIn, vmountOut)

		container := corev1.Container{
			Image:        m.Spec.ConfFromFileImage,
			Name:         "initcontainer" + strconv.Itoa(i),
			Env:          envs,
			VolumeMounts: vmounts,
		}
		initContainers = append(initContainers, container)
	}

	// create a helper init container that will choose the config files for this
	// set of datasets only and move them in a temps directory
	envs := []corev1.EnvVar{}
	vmounts := []corev1.VolumeMount{}
	count := 0
	for i, datasetname := range m.Spec.DatasetNames {
		env := corev1.EnvVar{Name: "DATASET_NAME_" + strconv.Itoa(i), Value: datasetname}
		envs = append(envs, env)
		env = corev1.EnvVar{Name: "DATASET_ENDPOINT_" + strconv.Itoa(i), Value: m.Spec.Endpoints[i]}
		envs = append(envs, env)
		count++
	}
	env := corev1.EnvVar{Name: "N", Value: strconv.Itoa(count - 1)}
	envs = append(envs, env)

	env = corev1.EnvVar{Name: "FEDERATION_NAME", Value: m.Name}
	envs = append(envs, env)

	env = corev1.EnvVar{Name: "FEDERATOR_NAME", Value: m.Spec.FederatorName}
	envs = append(envs, env)

	volumeHouse := corev1.Volume{
		Name: "nfs-housekeep",
		VolumeSource: corev1.VolumeSource{
			NFS: &corev1.NFSVolumeSource{
				Server: nfsip,
				Path:   "/exports"}}}

	vmountHouse := corev1.VolumeMount{
		Name:      "nfs-housekeep",
		MountPath: "/kobe"}

	volumes = append(volumes, volumeHouse)
	vmounts = append(vmounts, vmountHouse)

	containerHouse := corev1.Container{
		Image:        "kostbabis/housekeeping",
		Name:         "inithouse",
		Env:          envs,
		VolumeMounts: vmounts,
	}
	initContainers = append(initContainers, containerHouse)

	// create the initcontainer that will run the image that combines many
	// configs from the above temp directory and make appropriate config for the
	// whole experiment/federation
	vmounts = []corev1.VolumeMount{}

	path := "/exports/temp-" + m.Name

	volumeInFinal := corev1.Volume{
		Name: "nfs-final-in",
		VolumeSource: corev1.VolumeSource{
			NFS: &corev1.NFSVolumeSource{
				Server: nfsip,
				Path:   path}}}

	vmountInFinal := corev1.VolumeMount{
		Name:      "nfs-final-in",
		MountPath: m.Spec.InputDir}

	volumes = append(volumes, volumeInFinal)
	vmounts = append(vmounts, vmountInFinal)

	path = "/exports/" + m.Name

	volumeOutFinal := corev1.Volume{
		Name: "nfs-final-out",
		VolumeSource: corev1.VolumeSource{
			NFS: &corev1.NFSVolumeSource{
				Server: nfsip,
				Path:   path}}}

	vmountOutFinal := corev1.VolumeMount{
		Name:      "nfs-final-out",
		MountPath: m.Spec.OutputDir}

	volumes = append(volumes, volumeOutFinal)
	vmounts = append(vmounts, vmountOutFinal)

	container := corev1.Container{
		Image:        m.Spec.ConfImage,
		Name:         "init" + "final",
		Env:          envs,
		VolumeMounts: vmounts,
	}

	initContainers = append(initContainers, container)
	// initContainers should be:
	// m.Spec.ConfFromFileImage for each dataset
	// kostbabis/housekeeping
	// m.Spec.ConfImage - that initializes the federator engine

	// create the deployment of the federation. mount the config files to where
	// the federator needs (for example, `etc/default/semagrow`) --> passed by
	// the yaml of federator
	volumeConf := corev1.Volume{
		Name: "volumeconf",
		VolumeSource: corev1.VolumeSource{
			NFS: &corev1.NFSVolumeSource{
				Server: nfsip,
				Path:   "/exports/" + m.Name + "/",
			},
		},
	}

	mountConf := corev1.VolumeMount{
		Name:      "volumeconf",
		MountPath: m.Spec.FedConfDir,
	}

	volumes = append(volumes, volumeConf)

	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					InitContainers: initContainers,
					Containers: []corev1.Container{{
						Image:           m.Spec.Image,
						Name:            m.Name,
						ImagePullPolicy: m.Spec.ImagePullPolicy,
						Ports: []corev1.ContainerPort{{
							ContainerPort: m.Spec.Port,
							Name:          m.Name,
						}},
						VolumeMounts: []corev1.VolumeMount{mountConf},
						Resources:    m.Spec.Resources,
					}},
					Volumes:  volumes,
					Affinity: m.Spec.Affinity,
				},
			},
		},
	}
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}
*/

// a service to find the federation by name internally in the cluster.
func (r *ReconcileKobeFederation) newServiceForFederation(m *kobev1alpha1.KobeFederation) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},

		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"kobeoperator_cr": m.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Port: m.Spec.Port,
					TargetPort: intstr.IntOrString{
						IntVal: m.Spec.Port,
					},
				},
			},
		},
	}
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// A job that remove the (temporary) files and create some dirs...
func (r *ReconcileKobeFederation) newJobForFederation(m *kobev1alpha1.KobeFederation) *batchv1.Job {
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
				metav1.ObjectMeta{},
				corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           "busybox",
						Name:            m.Name,
						ImagePullPolicy: corev1.PullIfNotPresent,
						VolumeMounts:    vmounts,
						Command:         []string{"sh", "-c"},
						Args: []string{"cd /kobe ;rm -r " + m.Name + " ; rm -r" + " temp-" + m.Name + " ; for d in */; do   cd $d;  mkdir " + m.Spec.FederatorName +
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
func (r *ReconcileKobeFederation) newJobForDataset(m *kobev1alpha1.KobeFederation, dataset string) *batchv1.Job {
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
func (r *ReconcileKobeFederation) newPodForFederation(m *kobev1alpha1.KobeFederation, datasets []string, endpoints []string) *corev1.Pod {
	labels := labelsForKobeFederation(m.Name)

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
		envs := []corev1.EnvVar{}
		env := corev1.EnvVar{Name: "DATASET_NAME", Value: datasetname}
		envs = append(envs, env)

		env = corev1.EnvVar{Name: "DATASET_ENDPOINT", Value: endpoints[i]}
		envs = append(envs, env)

		if m.Spec.ForceNewInit == true { //optional variable to skip creating the files if they already exist in /exports/<dataset-name>/<federator-name>.Is passed by the experiment yaml
			env = corev1.EnvVar{Name: "INITIALIZE", Value: "yes"}
			envs = append(envs, env)
		}
		volumeIn := corev1.Volume{Name: "nfs-in-" + datasetname, VolumeSource: corev1.VolumeSource{NFS: &corev1.NFSVolumeSource{Server: nfsip, Path: "/exports/" + datasetname + "/dump"}}}
		volumeOut := corev1.Volume{Name: "nfs-out-" + datasetname, VolumeSource: corev1.VolumeSource{NFS: &corev1.NFSVolumeSource{Server: nfsip, Path: "/exports/" + datasetname + "/" + m.Spec.FederatorName}}}
		volumes = append(volumes, volumeIn, volumeOut)

		vmountIn := corev1.VolumeMount{Name: "nfs-in-" + datasetname, MountPath: m.Spec.InputDumpDir}
		vmountOut := corev1.VolumeMount{Name: "nfs-out-" + datasetname, MountPath: m.Spec.OutputDumpDir}
		vmounts = append(vmounts, vmountIn, vmountOut)

		container := corev1.Container{
			Image:        m.Spec.ConfFromFileImage,
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
	for i, datasetname := range m.Spec.DatasetNames {
		env := corev1.EnvVar{Name: "DATASET_NAME_" + strconv.Itoa(i), Value: datasetname}
		envs = append(envs, env)
		env = corev1.EnvVar{Name: "DATASET_ENDPOINT_" + strconv.Itoa(i), Value: m.Spec.Endpoints[i]}
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

	vmountInFinal := corev1.VolumeMount{Name: "nfs-final-in", MountPath: m.Spec.InputDir}
	vmounts = append(vmounts, vmountInFinal)

	path = "/exports/" + m.Name

	volumeOutFinal := corev1.Volume{Name: "nfs-final-out", VolumeSource: corev1.VolumeSource{NFS: &corev1.NFSVolumeSource{Server: nfsip, Path: path}}}
	volumes = append(volumes, volumeOutFinal)

	vmountOutFinal := corev1.VolumeMount{Name: "nfs-final-out", MountPath: m.Spec.OutputDir}
	vmounts = append(vmounts, vmountOutFinal)

	container := corev1.Container{
		Image:        m.Spec.ConfImage,
		Name:         "init" + "final",
		Env:          envs,
		VolumeMounts: vmounts,
	}
	initContainers = append(initContainers, container)

	//create the deployment of the federation .
	//mount the config files to where the federator needs (for example etc/default/semagrow) -->passed by the yaml of federator
	volumeConf := corev1.Volume{Name: "volumeconf", VolumeSource: corev1.VolumeSource{NFS: &corev1.NFSVolumeSource{Server: nfsip, Path: "/exports/" + m.Name + "/"}}}
	volumes = append(volumes, volumeConf)

	mountConf := corev1.VolumeMount{Name: "volumeconf", MountPath: m.Spec.FedConfDir}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{

			InitContainers: initContainers,
			Containers: []corev1.Container{{
				Image:           m.Spec.Image,
				Name:            m.Name,
				ImagePullPolicy: m.Spec.ImagePullPolicy,
				Ports: []corev1.ContainerPort{{
					ContainerPort: m.Spec.Port,
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
