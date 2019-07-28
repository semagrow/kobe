package kobedataset

import (
	"context"
	"reflect"
	"strconv"

	kobedatasetv1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobedataset/v1alpha1"
	kobeutilv1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobeutil/v1alpha1"
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

var log = logf.Log.WithName("controller_kobedataset")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new KobeDataset Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileKobeDataset{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("kobedataset-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource KobeDataset
	err = c.Watch(&source.Kind{Type: &kobedatasetv1alpha1.KobeDataset{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobedatasetv1alpha1.KobeDataset{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobedatasetv1alpha1.KobeDataset{},
	})
	if err != nil {
		return err
	}
	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner KobeDataset
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobedatasetv1alpha1.KobeDataset{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileKobeDataset implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileKobeDataset{}

// ReconcileKobeDataset reconciles a KobeDataset object
type ReconcileKobeDataset struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a KobeDataset object and makes changes based on the state read
// and what is in the KobeDataset.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
//-----------------------------------------reconciling-----------------------------------------
func (r *ReconcileKobeDataset) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling KobeDataset")

	// Fetch the KobeDataset instance
	instance := &kobedatasetv1alpha1.KobeDataset{}
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
	// if following fields are not present in yaml then set the defaults
	if instance.Spec.Group == "" {
		instance.Spec.Group = "kobe"
	}
	if instance.Spec.Count < 1 {
		instance.Spec.Count = 1
	}
	if instance.Spec.Image == "" {
		instance.Spec.Image = "kostbabis/virtuoso"
	}
	//check  if a KobeUtil instance exists for this namespace and if not create it
	kobeUtil := &kobeutilv1alpha1.KobeUtil{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "kobeutil", Namespace: instance.Namespace}, kobeUtil)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating the kobe utility custom resource")
		kobeutil := r.newKobeUtility(instance)
		err = r.client.Create(context.TODO(), kobeutil)
		if err != nil {
			reqLogger.Info("Failed to create the kobe utility instance: %v\n", err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}

	//check for  the nfs pod if it exist and wait if not
	nfsPodFound := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "kobenfs", Namespace: instance.Namespace}, nfsPodFound)
	if err != nil && errors.IsNotFound(err) {
		return reconcile.Result{Requeue: true}, nil
	}
	//check if the persistent volume claim exist and wait if not
	pvcFound := &corev1.PersistentVolumeClaim{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "kobepvc", Namespace: instance.Namespace}, pvcFound)
	if err != nil && errors.IsNotFound(err) {
		return reconcile.Result{Requeue: true}, nil
	}
	//make sure nfs pod status is running, to take care of racing condition
	if nfsPodFound.Status.Phase != "Running" {
		return reconcile.Result{Requeue: true}, nil
	}

	//----------------------------From here do the actual work to set up the pod and service for the dataset--------------
	// health check for the pods of dataset
	for i := 0; i < int(instance.Spec.Count); i++ {
		foundPod := &corev1.Pod{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name + "-kobedataset-" + strconv.Itoa(i), Namespace: instance.Namespace}, foundPod)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Making a new pod for kobedataset", "instance.Namespace", instance.Namespace, "instance.Name", instance.Name)
			pod := r.newPodForKobeDataset(instance, instance.Name+"-kobedataset-"+strconv.Itoa(i))
			err = r.client.Create(context.TODO(), pod)
			if err != nil {
				reqLogger.Info("Failed to create new Pod: %v\n", err)
				return reconcile.Result{}, err
			}
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			return reconcile.Result{}, err
		}
	}

	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(labelsForKobeDataset(instance.Name))
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

	// Update AppGroup status
	if instance.Spec.Group != instance.Status.AppGroup {
		instance.Status.AppGroup = instance.Spec.Group
		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Info("failed to update group status: %v", err)
			return reconcile.Result{}, err
		}
	}
	podForDelete := &corev1.Pod{}
	for i, podName := range podNames {
		if i >= int(instance.Spec.Count) { //check if we need to scale down the pods if user has changed count to lower number and delete if needed
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: podName, Namespace: instance.Namespace}, podForDelete)
			err = r.client.Delete(context.TODO(), podForDelete, client.PropagationPolicy(metav1.DeletionPropagation("Background")))
			if err != nil {
				reqLogger.Info("Failed to delete the federation job from the cluster")
				return reconcile.Result{}, err
			}
		}
	}
	//Service health check for the dataset
	foundService := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Making a new service for kobedataset", "instance.Namespace", instance.Namespace, "instance.Name", instance.Name)
		service := r.newServiceForDataset(instance)
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

	if podNames == nil && len(podNames) == 0 {
		return reconcile.Result{RequeueAfter: 15}, nil
	}
	if podNames != nil && len(podNames) > 0 {
		instance.Spec.ForceLoad = false
		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Info("failed to update the dataset forcedownload flag ")
			return reconcile.Result{}, err
		}
	}
	// -------------------------------finishing line everything should be fine here---------------------------------
	reqLogger.Info("Loop for a kobedataset went through the end for reconciling kobedataset\n")
	return reconcile.Result{}, nil

}

//------------------status update and retrieval to check actual pods besides deployment.Can be used to delay experiment for kobeexperiment till everything is up
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

//helper function
func labelsForKobeDataset(name string) map[string]string {
	return map[string]string{"app": "Kobe-Operator", "kobeoperator_cr": name}
}

//------------------Functions that define native kubernetes object to create that are all controlled by the kobedataset custom resource-----------------------
//--------tied to a dataset------
func (r *ReconcileKobeDataset) newDeploymentForKobeDataset(m *kobedatasetv1alpha1.KobeDataset) *appsv1.Deployment {
	labels := labelsForKobeDataset(m.Name)
	replicas := m.Spec.Count

	envs := []corev1.EnvVar{}

	env := corev1.EnvVar{Name: "DOWNLOAD_URL", Value: m.Spec.DownloadFrom}
	envs = append(envs, env)

	env = corev1.EnvVar{Name: "DATASET_NAME", Value: m.Name}
	envs = append(envs, env)

	if m.Spec.ForceLoad == true {
		env = corev1.EnvVar{Name: "FORCE_LOAD", Value: "YES"}
		envs = append(envs, env)
	}

	volume := corev1.Volume{Name: "nfs", VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "kobepvc"}}}
	volumes := []corev1.Volume{}
	volumes = append(volumes, volume)

	volumemount := corev1.VolumeMount{Name: "nfs", MountPath: "/kobe/dataset"}
	volumemounts := []corev1.VolumeMount{}
	volumemounts = append(volumemounts, volumemount)

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
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           m.Spec.Image,
						Name:            m.Name,
						ImagePullPolicy: corev1.PullPolicy("Always"),
						Ports: []corev1.ContainerPort{{
							ContainerPort: m.Spec.Port,
							Name:          m.Name,
						}},
						Env:          envs,
						VolumeMounts: volumemounts,
					}},
					Volumes: volumes,
				},
			},
		},
	}
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep

}

func (r *ReconcileKobeDataset) newPodForKobeDataset(m *kobedatasetv1alpha1.KobeDataset, podName string) *corev1.Pod {
	labels := labelsForKobeDataset(m.Name)

	envs := []corev1.EnvVar{}

	env := corev1.EnvVar{Name: "DOWNLOAD_URL", Value: m.Spec.DownloadFrom}
	envs = append(envs, env)

	env = corev1.EnvVar{Name: "DATASET_NAME", Value: m.Name}
	envs = append(envs, env)

	if m.Spec.ForceLoad == true {
		env = corev1.EnvVar{Name: "FORCE_LOAD", Value: "YES"}
		envs = append(envs, env)
	}

	volume := corev1.Volume{Name: "nfs", VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "kobepvc"}}}
	volumes := []corev1.Volume{}
	volumes = append(volumes, volume)

	volumemount := corev1.VolumeMount{Name: "nfs", MountPath: "/kobe/dataset"}
	volumemounts := []corev1.VolumeMount{}
	volumemounts = append(volumemounts, volumemount)

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: m.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{

			Containers: []corev1.Container{{
				Image:           m.Spec.Image,
				Name:            m.Name,
				ImagePullPolicy: corev1.PullPolicy("Always"),
				Ports: []corev1.ContainerPort{{
					ContainerPort: m.Spec.Port,
					Name:          m.Name,
				}},
				Env:          envs,
				VolumeMounts: volumemounts,
			}},
			Volumes: volumes,
		},
	}
	controllerutil.SetControllerReference(m, pod, r.scheme)
	return pod

}

func (r *ReconcileKobeDataset) newServiceForDataset(m *kobedatasetv1alpha1.KobeDataset) *corev1.Service {
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

func (r *ReconcileKobeDataset) newKobeUtility(m *kobedatasetv1alpha1.KobeDataset) *kobeutilv1alpha1.KobeUtil {
	kutil := &kobeutilv1alpha1.KobeUtil{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KobeUtil",
			APIVersion: "kobeutil.kobe.com/v1alpha1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      "kobeutil",
			Namespace: m.Namespace,
		},

		Spec: kobeutilv1alpha1.KobeUtilSpec{},
	}
	return kutil
}
