package kobefederation

import (
	"context"
	"reflect"
	"strconv"

	kobefederationv1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobefederation/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
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
	err = c.Watch(&source.Kind{Type: &kobefederationv1alpha1.KobeFederation{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner KobeFederator
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobefederationv1alpha1.KobeFederation{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobefederationv1alpha1.KobeFederation{},
	})
	if err != nil {
		return err
	}
	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner KobeDataset
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobefederationv1alpha1.KobeFederation{},
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

// Reconcile reads that state of the cluster for a KobeFederation object and makes changes based on the state read
// and what is in the KobeFederation.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileKobeFederation) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling KobeFederation")

	// Fetch the KobeFederation instance
	instance := &kobefederationv1alpha1.KobeFederation{}
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

	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new deployment for kobefederation", "instance.Namespace", instance.Namespace, "instance.Name", instance.Name)
		dep := r.newDeploymentForFederation(instance)
		reqLogger.Info("Creating a new Deployment %s/%s\n", dep.Namespace, dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Info("Failed to create new Deployment: %v\n", err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	//check for status changes
	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(labelsForKobeFederation(instance.Name))
	listOps := &client.ListOptions{Namespace: instance.Namespace, LabelSelector: labelSelector}
	err = r.client.List(context.TODO(), listOps, podList)
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

	//update kobe federation affinity
	if instance.Spec.Affinity.NodeAffinity != nil || instance.Spec.Affinity.PodAffinity != nil || instance.Spec.Affinity.PodAntiAffinity != nil {
		affinity := instance.Spec.Affinity
		if *found.Spec.Template.Spec.Affinity != affinity {
			found.Spec.Template.Spec.Affinity = &affinity
			err = r.client.Update(context.TODO(), found)
			if err != nil {
				reqLogger.Info("Failed to update Deployment: %v\n", err)
				return reconcile.Result{}, err
			}
			// Spec updated return and reque .Affinity fixed possible other fixes like this here later
			return reconcile.Result{Requeue: true}, nil

		}
	}
	//check the healthiness of the federation service
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

func (r *ReconcileKobeFederation) newDeploymentForFederation(m *kobefederationv1alpha1.KobeFederation) *appsv1.Deployment {
	labels := labelsForKobeFederation(m.Name)

	//-------------------------------------------------------crap--------------------------------------------
	nfsPodFound := &corev1.Pod{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: "kobenfs", Namespace: m.Namespace}, nfsPodFound)
	if err != nil && errors.IsNotFound(err) {
		return nil
	}
	nfsip := nfsPodFound.Status.PodIP //it seems we need this cause dns for service of the nfs doesnt work in kubernetes
	//-------------------------------------------------------/crap------------------------------------

	//create init containers definitions that make one config file for federation per dataset
	initContainers := []corev1.Container{}
	volumes := []corev1.Volume{}
	//initContainers = append(initContainers, corev1.Container{Image: "busybox", Name: "firstcont", Command: []string{"mkdir"}, Args: []string{"/" + m.Name}})
	for i, datasetname := range m.Spec.DatasetNames {
		//each init container is given DATASET_NAME and DATASET_ENDPOINT environment variables to work with)
		//also inputfiledir and outputfiledir both point to exports/<datasetname>/dumps/ and exports/dataset/<federation>/ respectively to nfs server
		vmounts := []corev1.VolumeMount{}
		envs := []corev1.EnvVar{}
		env := corev1.EnvVar{Name: "DATASET_NAME", Value: datasetname}
		envs = append(envs, env)
		env = corev1.EnvVar{Name: "DATASET_ENDPOINT", Value: "http://" + m.Spec.Endpoints[i] + ":" + "8890"}
		envs = append(envs, env)

		volumeIn := corev1.Volume{Name: "nfs-in-" + datasetname, VolumeSource: corev1.VolumeSource{NFS: &corev1.NFSVolumeSource{Server: nfsip, Path: "/exports/" + datasetname + "/dump"}}}
		volumeOut := corev1.Volume{Name: "nfs-out-" + datasetname, VolumeSource: corev1.VolumeSource{NFS: &corev1.NFSVolumeSource{Server: nfsip, Path: "/exports/" + datasetname + "/"}}}
		volumes = append(volumes, volumeIn, volumeOut)

		vmountIn := corev1.VolumeMount{Name: "nfs-in-" + datasetname, MountPath: m.Spec.InputFileDir}
		vmountOut := corev1.VolumeMount{Name: "nfs-out-" + datasetname, MountPath: m.Spec.OutputFileDir}
		vmounts = append(vmounts, vmountIn, vmountOut)

		container := corev1.Container{
			Image:        m.Spec.ConfFromFileImage,
			Name:         "initcontainer" + strconv.Itoa(i),
			Env:          envs,
			VolumeMounts: vmounts,
		}
		initContainers = append(initContainers, container)
	}

	//create the initcontainer that will run the image that combines many configs (1 per dataset) to one config for the experiment
	envs := []corev1.EnvVar{}
	vmounts := []corev1.VolumeMount{}
	for i, datasetname := range m.Spec.DatasetNames {
		env := corev1.EnvVar{Name: "DATASET_NAME_" + strconv.Itoa(i), Value: datasetname}
		envs = append(envs, env)
		env = corev1.EnvVar{Name: "DATASET_NAME_" + strconv.Itoa(i), Value: "http://" + m.Spec.Endpoints[i] + ":" + "8890"}
		envs = append(envs, env)
	}
	volumeFinal := corev1.Volume{Name: "nfs-final", VolumeSource: corev1.VolumeSource{NFS: &corev1.NFSVolumeSource{Server: nfsip, Path: "/exports/"}}}
	volumes = append(volumes, volumeFinal)

	vmountFinal := corev1.VolumeMount{Name: "nfs-final", MountPath: "/"}
	vmounts = append(vmounts, vmountFinal)

	container := corev1.Container{
		Image:        m.Spec.ConfImage,
		Name:         "init" + "final",
		Env:          envs,
		VolumeMounts: vmounts,
	}
	initContainers = append(initContainers, container)

	//create the deployment of the federation .
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
					}},
					Volumes: volumes,
				},
			},
		},
	}
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep

}

func (r *ReconcileKobeFederation) newServiceForFederation(m *kobefederationv1alpha1.KobeFederation) *corev1.Service {
	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},

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
