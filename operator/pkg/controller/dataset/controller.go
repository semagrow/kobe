package dataset

import (
	"context"
	"reflect"
	"strconv"

	api "github.com/semagrow/kobe/operator/pkg/apis/kobe/v1alpha1"
	istioapi "istio.io/api/networking/v1alpha3"
	istioclient "istio.io/client-go/pkg/apis/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
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

var log = logf.Log.WithName("controller_dataset")

// Add creates a new Dataset Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDataset{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("dataset-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Dataset
	err = c.Watch(&source.Kind{Type: &api.Dataset{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &api.Dataset{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &api.Dataset{},
	})
	if err != nil {
		return err
	}
	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Dataset
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &api.Dataset{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileDataset implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileDataset{}

// ReconcileDataset reconciles a Dataset object
type ReconcileDataset struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Dataset object and makes changes based on the state read
// and what is in the Dataset.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
//-----------------------------------------reconciling-----------------------------------------
func (r *ReconcileDataset) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Dataset")

	// Fetch the Dataset instance
	instance := &api.Dataset{}
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

	instance.SetDefaults()
	// check ForceLoad

	// check  if a KobeUtil instance exists for this namespace and if not create it
	requeue, err := ensureNFS(r.client, instance.Namespace)
	if requeue {
		return reconcile.Result{Requeue: requeue}, err
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// From here do the actual work to set up the pod and service for the dataset

	created, err := r.reconcilePods(instance)
	if created {
		return reconcile.Result{Requeue: true}, err
	} else if err != nil {
		return reconcile.Result{}, err
	}

	//Service health check for the dataset
	if err := r.reconcileSvc(instance); err != nil {
		return reconcile.Result{}, err
	}

	if instance.Status.PodNames == nil && len(instance.Status.PodNames) == 0 {
		return reconcile.Result{RequeueAfter: 25}, nil
	}
	if instance.Status.PodNames != nil && len(instance.Status.PodNames) > 0 {
		instance.Status.ForceLoad = false
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Info("failed to update the dataset forcedownload flag")
			return reconcile.Result{}, err
		}
	}
	// -------------------------------finishing line everything should be fine here---------------------------------
	reqLogger.Info("Loop for a dataset went through the end for reconciling dataset\n")
	return reconcile.Result{}, nil
}

//helper function
func labelsForDataset(name string) map[string]string {
	return map[string]string{"app": "Kobe-Operator", "kobeoperator_cr": name}
}

// status update and retrieval to check actual pods besides deployment.
// Can be used to delay experiment for Experiment till everything is up
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

func (r *ReconcileDataset) reconcilePods(instance *api.Dataset) (bool, error) {
	reqLogger := log

	// health check for the pods of dataset
	for i := 0; i < int(*instance.Spec.Replicas); i++ {
		foundPod := &corev1.Pod{}
		err := r.client.Get(context.TODO(), types.NamespacedName{
			Name:      instance.Name + "-dataset-" + strconv.Itoa(i),
			Namespace: instance.Namespace},
			foundPod)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Making a new pod for dataset", "instance.Namespace", instance.Namespace, "instance.Name", instance.Name)
			pod := r.newPod(instance, instance.Name+"-dataset-"+strconv.Itoa(i))
			err = r.client.Create(context.TODO(), pod)
			if err != nil {
				reqLogger.Info("Failed to create new Pod: %v\n", err)
				return false, err
			}
			return true, nil
		} else if err != nil {
			return false, err
		}
	}

	podList := &corev1.PodList{}
	listOps := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels(labelsForDataset(instance.Name)),
	}
	err := r.client.List(context.TODO(), podList, listOps...)
	if err != nil {
		reqLogger.Info("Failed to list pods: %v", err)
		return false, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.PodNames if needed
	if !reflect.DeepEqual(podNames, instance.Status.PodNames) {
		instance.Status.PodNames = podNames
		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Info("failed to update node status: %v", err)
			return false, err
		}
	}

	podForDelete := &corev1.Pod{}
	for i, podName := range podNames {
		if i >= int(*instance.Spec.Replicas) {
			//check if we need to scale down the pods if user has changed count to lower number and delete if needed
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: podName, Namespace: instance.Namespace}, podForDelete)
			err = r.client.Delete(context.TODO(), podForDelete, client.PropagationPolicy(metav1.DeletionPropagation("Background")))
			if err != nil {
				reqLogger.Info("Failed to delete the federation job from the cluster")
				return false, err
			}
		}
	}
	return false, nil
}

func (r *ReconcileDataset) reconcileSvc(instance *api.Dataset) error {
	reqLogger := log
	foundService := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundService)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Making a new service for dataset", "instance.Namespace", instance.Namespace, "instance.Name", instance.Name)
		service := r.newSvc(instance)
		reqLogger.Info("Creating a new Service %s/%s\n", service.Namespace, service.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			reqLogger.Info("Failed to create new Service: %v\n", err)
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	foundVirtualService := &istioclient.VirtualService{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundVirtualService)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Making a new virtual service for dataset", "instance.Namespace", instance.Namespace, "instance.Name", instance.Name)
		service := r.newVirtualSvc(instance)
		reqLogger.Info("Creating a new VRITUAL Service %s/%s\n", service.Namespace, service.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			reqLogger.Info("Failed to create new Service: %v\n", err)
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	return nil
}

func (r *ReconcileDataset) newPod(m *api.Dataset, podName string) *corev1.Pod {
	labels := labelsForDataset(m.Name)

	envs := []corev1.EnvVar{
		{Name: "DOWNLOAD_URL", Value: m.Spec.DownloadFrom},
		{Name: "DATASET_NAME", Value: m.Name},
	}

	if m.Status.ForceLoad == true {
		envs = append(envs, corev1.EnvVar{Name: "FORCE_LOAD", Value: "YES"})
	}

	for _, v := range m.Spec.Env {
		envs = append(envs, corev1.EnvVar{Name: v.Name, Value: v.Value})
	}

	volume := corev1.Volume{
		Name: "nfs",
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "kobepvc"}}}

	volumes := []corev1.Volume{}
	volumes = append(volumes, volume)

	volumemount := corev1.VolumeMount{
		Name:      "nfs",
		MountPath: "/kobe/dataset"}

	volumemounts := []corev1.VolumeMount{}
	volumemounts = append(volumemounts, volumemount)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: m.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Image:           m.Spec.Image,
				Name:            m.Name,
				ImagePullPolicy: m.Spec.ImagePullPolicy,
				Ports: []corev1.ContainerPort{{
					ContainerPort: m.Spec.Port,
					Name:          m.Name,
				}},
				Env:          envs,
				VolumeMounts: volumemounts,
				Resources:    m.Spec.Resources,
			}},
			Volumes:  volumes,
			Affinity: m.Spec.Affinity,
		},
	}
	controllerutil.SetControllerReference(m, pod, r.scheme)
	return pod
}

func (r *ReconcileDataset) newSvc(m *api.Dataset) *corev1.Service {
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

func (r *ReconcileDataset) newVirtualSvc(m *api.Dataset) *istioclient.VirtualService {
	// service := &corev1.Service{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name:      m.Name,
	// 		Namespace: m.Namespace,
	// 	},

	// 	Spec: corev1.ServiceSpec{
	// 		Selector: map[string]string{
	// 			"kobeoperator_cr": m.Name,
	// 		},
	// 		Ports: []corev1.ServicePort{
	// 			{
	// 				Port: m.Spec.Port,
	// 				TargetPort: intstr.IntOrString{
	// 					IntVal: m.Spec.Port,
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	vservice := &istioclient.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: istioapi.VirtualService{
			Hosts: []string{m.Name},
		},
	}

	controllerutil.SetControllerReference(m, vservice, r.scheme)
	return vservice
}

func ensureNFS(client client.Client, ns string) (bool, error) {
	reqLogger := log

	// check  if a KobeUtil instance exists for this namespace and if not create it
	kobeUtil := &api.KobeUtil{}

	err := client.Get(context.TODO(), types.NamespacedName{Name: "kobeutil", Namespace: ns}, kobeUtil)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating the kobe utility custom resource")
		kobeutil := newKobeUtility(ns)
		err = client.Create(context.TODO(), kobeutil)
		if err != nil {
			reqLogger.Info("Failed to create the kobe utility instance: %v\n", err)
			return false, err
		}
		return true, nil
	}

	//check for  the nfs pod if it exist and wait if not
	nfsPodFound := &corev1.Pod{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: "kobenfs", Namespace: ns}, nfsPodFound)
	if err != nil && errors.IsNotFound(err) {
		return true, nil
	}
	//check if the persistent volume claim exist and wait if not
	pvcFound := &corev1.PersistentVolumeClaim{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: "kobepvc", Namespace: ns}, pvcFound)
	if err != nil && errors.IsNotFound(err) {
		return true, nil
	}
	//make sure nfs pod status is running, to take care of racing condition
	if nfsPodFound.Status.Phase != corev1.PodRunning {
		return true, nil
	}

	return false, nil
}

func newKobeUtility(ns string) *api.KobeUtil {
	kutil := &api.KobeUtil{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kobeutil",
			Namespace: ns,
		},
	}
	return kutil
}
