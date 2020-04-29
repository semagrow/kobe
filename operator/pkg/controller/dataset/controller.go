package dataset

import (
	"context"
	"reflect"

	api "github.com/semagrow/kobe/operator/pkg/apis/kobe/v1alpha1"
	istioapi "istio.io/api/networking/v1alpha3"
	istioclient "istio.io/client-go/pkg/apis/networking/v1alpha3"
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
	err = c.Watch(&source.Kind{Type: &api.EphemeralDataset{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &api.EphemeralDataset{},
	})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Dataset
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &api.EphemeralDataset{},
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
	instance := &api.EphemeralDataset{}
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

	//check if template in the template field exists. If not try to find a reference template and set that to the template fields .
	if instance.Spec.Template == nil {
		foundTemplate := &api.DatasetTemplate{}
		reqLogger.Info("Finding the template reference specified for " + instance.Name + " %v\n")
		err := r.client.Get(context.TODO(), types.NamespacedName{
			Name:      instance.Spec.TemplateRef,
			Namespace: corev1.NamespaceDefault},
			foundTemplate)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Failed to create new Service: %v\n", err)
			return reconcile.Result{}, err
		}

		instance.Spec.Template = foundTemplate
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
func labelsForDataset(m *api.EphemeralDataset) map[string]string {
	return map[string]string{"app": "Kobe-Operator", "datasetName": m.Name, "benchmark": m.Namespace}
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

func (r *ReconcileDataset) reconcilePods(instance *api.EphemeralDataset) (bool, error) {
	reqLogger := log

	// health check for the pods of dataset
	foundPod := &corev1.Pod{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      instance.Name,
		Namespace: instance.Namespace},
		foundPod)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Making a new pod for dataset", "instance.Namespace", instance.Namespace, "instance.Name", instance.Name)
		pod := r.newPod(instance)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			reqLogger.Info("Failed to create new Pod: %v\n", err)
			return false, err
		}
		return true, nil
	} else if err != nil {
		return false, err
	}

	podList := &corev1.PodList{}
	listOps := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels(labelsForDataset(instance)),
	}
	err = r.client.List(context.TODO(), podList, listOps...)
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
		if i >= 1 {
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

func (r *ReconcileDataset) reconcileSvc(instance *api.EphemeralDataset) error {
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

func (r *ReconcileDataset) newPod(m *api.EphemeralDataset) *corev1.Pod {
	labels := labelsForDataset(m)
	//reqLogger := log

	envs := []corev1.EnvVar{
		{Name: "DOWNLOAD_URL", Value: m.Spec.Files[1].Checksum},
		{Name: "DATASET_NAME", Value: m.Name},
	}

	if m.Status.ForceLoad == true {
		envs = append(envs, corev1.EnvVar{Name: "FORCE_LOAD", Value: "YES"})
	}

	volume := corev1.Volume{
		Name: "nfs",
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "kobepvc"}}}

	volumes := []corev1.Volume{}
	volumes = append(volumes, volume)

	initContainers := m.Spec.Template.TemplateSpec.InitContainers

	if m.Status.ForceLoad == true {
		// add volumemounts to importcontainers
		volumemount := corev1.VolumeMount{
			Name:      "nfs",
			MountPath: "/kobe/dataset"}

		volumemounts := []corev1.VolumeMount{}
		volumemounts = append(volumemounts, volumemount)

		initContainers = append(initContainers, m.Spec.Template.TemplateSpec.ImportContainers...)
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-pod",
			Namespace: m.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			InitContainers: initContainers,
			Containers:     m.Spec.Template.TemplateSpec.Containers,
			Volumes:        volumes,
			Affinity:       m.Spec.Affinity,
		},
	}
	controllerutil.SetControllerReference(m, pod, r.scheme)
	return pod
}

func (r *ReconcileDataset) newSvc(m *api.EphemeralDataset) *corev1.Service {
	servicePorts := []corev1.ServicePort{}
	for _, container := range m.Spec.Template.TemplateSpec.Containers {
		for _, port := range container.Ports {
			newPort := corev1.ServicePort{
				Port: port.ContainerPort,
				TargetPort: intstr.IntOrString{
					IntVal: port.ContainerPort,
				},
			}
			servicePorts = append(servicePorts, newPort)
		}
	}
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},

		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"kobeoperator_cr": m.Name,
				"benchmark":       m.Namespace,
			},
			Ports: servicePorts,
		},
	}
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

func (r *ReconcileDataset) newVirtualSvc(m *api.EphemeralDataset) *istioclient.VirtualService {
	vsvc := &istioclient.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: istioapi.VirtualService{
			Hosts: []string{m.Name},
		},
	}

	controllerutil.SetControllerReference(m, vsvc, r.scheme)
	return vsvc
}
