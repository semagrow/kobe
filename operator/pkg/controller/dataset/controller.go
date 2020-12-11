package dataset

import (
	"context"
	"reflect"
	"strconv"

	findtypes "github.com/gogo/protobuf/types"
	api "github.com/semagrow/kobe/operator/pkg/apis/kobe/v1alpha1"
	istioapi "istio.io/api/networking/v1alpha3"
	istioclient "istio.io/client-go/pkg/apis/networking/v1alpha3"
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

	//check if template in the template field exists. If not try to find a reference template and set that to the template fields .
	if instance.Spec.SystemSpec == nil {
		foundTemplate := &api.DatasetTemplate{}
		reqLogger.Info("Finding the template reference specified for " + instance.Name + "\n")
		err := r.client.Get(context.TODO(), types.NamespacedName{
			Name:      instance.Spec.TemplateRef,
			Namespace: corev1.NamespaceDefault},
			foundTemplate)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Failed to find the requested dataset template: ", err)
			return reconcile.Result{}, err
		}
		instance.Spec.SystemSpec = &foundTemplate.Spec
		err = r.client.Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Info("failed to update the template of the dataset"+instance.Spec.Name+": %v", instance.Spec.Name)
			return reconcile.Result{}, err
		}
	}

	instance.SetDefaults()
	// check ForceLoad

	//Service health check for the dataset
	if requeue, err := r.reconcileSvc(instance); requeue == true {
		return reconcile.Result{RequeueAfter: 5000000000}, err
	}

	// From here do the actual work to set up the pod and service for the dataset
	created, err := r.reconcilePods(instance)
	if created {
		return reconcile.Result{Requeue: true}, err
	} else if err != nil {
		return reconcile.Result{RequeueAfter: 10000000000}, err
	}

	if instance.Status.PodNames == nil && len(instance.Status.PodNames) == 0 {
		return reconcile.Result{RequeueAfter: 10000000000}, nil
	}
	// if instance.Status.PodNames != nil && len(instance.Status.PodNames) > 0 {
	// 	instance.Status.ForceLoad = false //????????????/
	// 	err := r.client.Status().Update(context.TODO(), instance)
	// 	if err != nil {
	// 		reqLogger.Info("failed to update the dataset forcedownload flag")
	// 		return reconcile.Result{}, err
	// 	}
	// }
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
		Name:      instance.Name + "-pod",
		Namespace: instance.Namespace},
		foundPod)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Making a new pod for dataset", "instance.Namespace", instance.Namespace, "instance.Name", instance.Name)
		pod := r.newPod(instance)
		if pod == nil {
			reqLogger.Info("Pod was not created. Requeue after waiting for 10 seconds\n")
			return true, nil
		}
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			reqLogger.Info("Failed to create new Pod : %v\n", err)
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
			reqLogger.Info("Failed to update the Dataset status: %v", err)
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
				reqLogger.Info("Failed to delete the extra pod from the cluster\n")
				return false, err
			}
		}
	}
	return false, nil
}

func (r *ReconcileDataset) newPod(m *api.EphemeralDataset) *corev1.Pod {
	labels := labelsForDataset(m)
	reqLogger := log

	envs := []corev1.EnvVar{
		{Name: "DOWNLOAD_URL", Value: m.Spec.Files[0].URL},
		{Name: "DATASET_NAME", Value: m.Name},
	}
	// fetch the benchmark to check for istio usage
	foundBenchmark := &api.Benchmark{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: m.Namespace, Namespace: corev1.NamespaceDefault}, foundBenchmark)
	if err != nil {
		reqLogger.Info("Failed to find the Benchmark this Dataset belongs to: %v\n", err)
		return nil
	}
	if foundBenchmark.Status.Istio == api.IstioUse {
		envs = append(envs, corev1.EnvVar{Name: "USE_ISTIO", Value: "YES"})
	} else if foundBenchmark.Status.Istio == api.IstioNotUse {
		envs = append(envs, corev1.EnvVar{Name: "USE_ISTIO", Value: "NO"})
	} else {
		reqLogger.Info("Status field for Istio Usage not properly set.\n")
		return nil
	}

	if m.Status.ForceLoad == true {
		envs = append(envs, corev1.EnvVar{Name: "FORCE_LOAD", Value: "YES"})
	}

	//add the env vars from above taken at the dataset level to the existing variables for each container (for now)
	for i, container := range m.Spec.SystemSpec.Containers {
		m.Spec.SystemSpec.Containers[i].Env = append(container.Env, envs...)
	}

	for i, container := range m.Spec.SystemSpec.InitContainers {
		m.Spec.SystemSpec.InitContainers[i].Env = append(container.Env, envs...)
	}

	for i, container := range m.Spec.SystemSpec.ImportContainers {
		m.Spec.SystemSpec.InitContainers[i].Env = append(container.Env, envs...)
	}

	nfsPodFound := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "kobenfs", Namespace: corev1.NamespaceDefault}, nfsPodFound)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Nfs server is not online yet: %v\n", err)
		return nil
	}
	nfsip := nfsPodFound.Status.PodIP //it seems we need this cause dns for service of the nfs doesnt work in kubernetes
	volume := corev1.Volume{Name: "nfs", VolumeSource: corev1.VolumeSource{NFS: &corev1.NFSVolumeSource{Server: nfsip, Path: "/exports/"}}}

	volumes := []corev1.Volume{}
	volumes = append(volumes, volume)

	initContainers := m.Spec.SystemSpec.InitContainers
	volumemount := corev1.VolumeMount{
		Name:      "nfs",
		MountPath: "/kobe/dataset"}

	volumemounts := []corev1.VolumeMount{}
	volumemounts = append(volumemounts, volumemount)

	initContainers = append(initContainers, m.Spec.SystemSpec.ImportContainers...)

	//supply all containers with their mount to the nfs
	//currently we give same mounts to every container in the pod
	for i := range m.Spec.SystemSpec.Containers {
		m.Spec.SystemSpec.Containers[i].VolumeMounts = append(m.Spec.SystemSpec.Containers[i].VolumeMounts, volumemounts...)
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-pod",
			Namespace: m.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			InitContainers: initContainers,
			Containers:     m.Spec.SystemSpec.Containers,
			Volumes:        volumes,
			Affinity:       m.Spec.Affinity,
		},
	}
	controllerutil.SetControllerReference(m, pod, r.scheme)
	return pod
}

func (r *ReconcileDataset) reconcileSvc(instance *api.EphemeralDataset) (bool, error) {
	reqLogger := log

	//find the associated benchmark and check the status field for istio
	foundBenchmark := &api.Benchmark{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Namespace, Namespace: corev1.NamespaceDefault}, foundBenchmark)
	if err != nil {
		reqLogger.Info("Failed to find the Benchmark this Dataset belongs to: %v\n", err)
		return true, err
	}

	foundService := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundService)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Making a new service for dataset", "instance.Namespace", instance.Namespace, "instance.Name", instance.Name)
		service := r.newSvc(instance)
		reqLogger.Info("Creating a new Service %s/%s\n", service.Namespace, service.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			reqLogger.Info("Failed to create new Service: %v\n", err)
			return true, err
		}
		return true, nil
	} else if err != nil {
		return true, err
	}

	if foundBenchmark.Status.Istio == api.IstioNotUse {
		return false, nil
	} else if foundBenchmark.Status.Istio == api.IstioUse {

		foundVirtualService := &istioclient.VirtualService{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundVirtualService)

		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Making a new virtual service for dataset", "instance.Namespace", instance.Namespace, "instance.Name", instance.Name)
			service := r.newVirtualSvc(instance)
			reqLogger.Info("Creating a new VRITUAL Service %s/%s\n", service.Namespace, service.Name)
			err = r.client.Create(context.TODO(), service)
			if err != nil {
				reqLogger.Info("Failed to create new Virtual Service: %v\n", err)
				return true, err
			}
			return true, nil
		} else if err != nil {
			return true, err
		}
	} else {
		reqLogger.Info("Flag to use Istio not set")
		return true, nil
	}
	return false, nil
}

func (r *ReconcileDataset) newSvc(m *api.EphemeralDataset) *corev1.Service {
	servicePorts := []corev1.ServicePort{{
		Port: int32(m.Spec.SystemSpec.Port),
		Name: "http",
	}}

	for i, container := range m.Spec.SystemSpec.Containers {
		for j, port := range container.Ports {
			if port.ContainerPort != int32(m.Spec.SystemSpec.Port) {

				newPort := corev1.ServicePort{
					Name: "http-" + strconv.Itoa(i) + strconv.Itoa(j),
					Port: port.ContainerPort,
				}
				servicePorts = append(servicePorts, newPort)
			}
		}
	}
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    labelsForDataset(m),
		},

		Spec: corev1.ServiceSpec{
			Selector: labelsForDataset(m),
			Ports:    servicePorts,
		},
	}
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

func (r *ReconcileDataset) newVirtualSvc(m *api.EphemeralDataset) *istioclient.VirtualService {
	http := []*istioapi.HTTPRoute{}
	if m.Spec.FederatorConnection != nil {
		m.Spec.FederatorConnection.Source = nil
		m.Spec.NetworkTopology = append(m.Spec.NetworkTopology, *m.Spec.FederatorConnection)
	}

	//add the fault injection based on network topology in the dataset virtual service
	for _, incoming := range m.Spec.NetworkTopology {

		httpMatchRequests := []*istioapi.HTTPMatchRequest{}
		match := istioapi.HTTPMatchRequest{}
		if incoming.Source == nil {
			match = istioapi.HTTPMatchRequest{
				SourceLabels: map[string]string{"kobe-resource": "federation"},
				Port:         m.Spec.SystemSpec.Port,
			}
		} else {
			match = istioapi.HTTPMatchRequest{
				SourceLabels: map[string]string{"app": "Kobe-Operator", "datasetName": *incoming.Source, "benchmark": m.Namespace},
				Port:         m.Spec.SystemSpec.Port,
			}
		}
		httpMatchRequests = append(httpMatchRequests, &match)

		route := []*istioapi.HTTPRouteDestination{
			{
				Destination: &istioapi.Destination{
					Host: m.Name, //+ "." + m.Namespace,
					Port: &istioapi.PortSelector{
						Number: m.Spec.SystemSpec.Port,
					},
				},
			},
		}
		fixedDelay := &findtypes.Duration{}
		if incoming.DelayInjection.FixedDelaySec != nil {
			fixedDelay.Seconds = int64(*incoming.DelayInjection.FixedDelaySec)
		}
		if incoming.DelayInjection.FixedDelayMSec != nil {
			fixedDelay.Nanos = int32(*incoming.DelayInjection.FixedDelayMSec * 1000000)
		}
		fault := &istioapi.HTTPFaultInjection{
			Delay: &istioapi.HTTPFaultInjection_Delay{
				HttpDelayType: &istioapi.HTTPFaultInjection_Delay_FixedDelay{
					FixedDelay: fixedDelay,
				},
				Percentage: &istioapi.Percent{Value: float64(*incoming.DelayInjection.Percentage)},
			},
		}

		httpRoute := istioapi.HTTPRoute{
			Match: httpMatchRequests,
			Route: route,
			Fault: fault,
		}

		http = append(http, &httpRoute)
	}

	//append dummy http so its not empty. Do not remove this!
	// httpMatchRequests := []*istioapi.HTTPMatchRequest{{
	// 	Name:         "dummy",
	// 	SourceLabels: map[string]string{"datasetName": "dummy", "benchmark": m.Namespace},
	// 	Port:         m.Spec.SystemSpec.Port,
	// }}
	route := []*istioapi.HTTPRouteDestination{
		{
			Destination: &istioapi.Destination{
				Host: m.Name, //+ "." + m.Namespace,
			},
		},
	}
	// fixedDelay := &findtypes.Duration{}
	// fixedDelay.Seconds = int64(10)
	// fault := &istioapi.HTTPFaultInjection{
	// 	Delay: &istioapi.HTTPFaultInjection_Delay{
	// 		HttpDelayType: &istioapi.HTTPFaultInjection_Delay_FixedDelay{
	// 			FixedDelay: fixedDelay,
	// 		},
	// 		Percentage: &istioapi.Percent{Value: float64(95)},
	// 	},
	// }
	http = append(http, &istioapi.HTTPRoute{
		Route: route,
		//Fault: fault,
	})

	vsvc := &istioclient.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: istioapi.VirtualService{
			Hosts: []string{m.Name}, //+ "." + m.Namespace},
			Http:  http,
		},
	}

	controllerutil.SetControllerReference(m, vsvc, r.scheme)
	return vsvc
}
