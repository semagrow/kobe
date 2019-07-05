package kobedataset

import (
	"context"
	"reflect"

	kobedatasetv1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobedataset/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
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
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
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

	//why are these here  in this controller????????????????

	//-------------------------------------------------checking if nfs config map health----------------------------------
	nfsconfigFound := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "nfsconfig", Namespace: instance.Namespace}, nfsconfigFound)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating config map for nfs\n") //used by the dataset pods to store the db file and regain it if that dataset is restarted
		conf := r.newNfsConfig(instance)
		err = r.client.Create(context.TODO(), conf)
		if err != nil {
			reqLogger.Info("Failed to create the nfs config: %v\n", err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}
	//------------------------------------------------checking for nfs server health---------------------------------------------------
	nfsServiceFound := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "kobenfs", Namespace: instance.Namespace}, nfsServiceFound)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Setting up the nfs service and server\n") //used by the dataset pods to store the db file and regain it if that dataset is restarted
		serv := r.newServiceForNfs(instance)
		err = r.client.Create(context.TODO(), serv)
		if err != nil {
			reqLogger.Info("Failed to create the kobe nfs Service: %v\n", err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil

	}
	nfsPodFound := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "kobenfs", Namespace: instance.Namespace}, nfsPodFound)
	if err != nil && errors.IsNotFound(err) {
		pod := r.newPodForNfs(instance)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			reqLogger.Info("Failed to create the kobe nfs Pod: %v\n", err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{RequeueAfter: 5}, nil

	}
	nfsip := nfsPodFound.Status.PodIP //it seems we need this cause dns for service of the nfs doesnt work in kubernetes

	//--------------------------------------------Persistent volume health check-----------------------------------------------
	pvFound := &corev1.PersistentVolume{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "kobepv"}, pvFound)
	if err != nil && errors.IsNotFound(err) {
		pv := r.newPvForKobe(instance, nfsip)
		err = r.client.Create(context.TODO(), pv)
		if err != nil {
			reqLogger.Info("Failed to create the persistent volume that the datasets will use to retain their data if they shutdown and restarted")
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil

	}
	//--------------------------------------------Persistent volume claim  health check---------------------------------------------
	pvcFound := &corev1.PersistentVolumeClaim{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "kobepvc", Namespace: instance.Namespace}, pvcFound)
	if err != nil && errors.IsNotFound(err) {
		pvc := r.newPvcForKobe(instance)
		err := r.client.Create(context.TODO(), pvc)
		if err != nil {
			reqLogger.Info("Failed to create the single persistent volume claim that all the datasets gonna use to mount their data directories to nfs")
			return reconcile.Result{}, err
		}
		return reconcile.Result{RequeueAfter: 1}, nil
	}

	//----------------------------From here do the actual work to set up the deployment and service for the dataset--------------

	// deployment health check for dataset
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
	reqLogger.Info("HEY THERE")
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Making a new deployment for kobedataset", "instance.Namespace", instance.Namespace, "instance.Name", instance.Name)
		dep := r.newDeploymentForKobeDataset(instance)
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

	count := instance.Spec.Count
	if *found.Spec.Replicas != count {
		found.Spec.Replicas = &count

		err = r.client.Update(context.TODO(), found)
		if err != nil {
			reqLogger.Info("Failed to update Deployment: %v\n", err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
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

	// -------------------------------finishingline everything should be fine here---------------------------------
	reqLogger.Info("Loop for a kobedataset went through the end for reconciling kobedataset\n")
	return reconcile.Result{}, err

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
						ImagePullPolicy: m.Spec.ImagePullPolicy,
						Ports: []corev1.ContainerPort{{
							ContainerPort: m.Spec.Port,
							Name:          m.Name,
						}},
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

//-------------------functions that create native kubernetes objects thatare checked when there is a kobedataset change by this controller -----------------------
//-------------------they are not tied to a kobedataset resource (todo in future)----------------

//NFS SERVICE (its actually useless cause nfs service dns bug)
func (r *ReconcileKobeDataset) newServiceForNfs(m *kobedatasetv1alpha1.KobeDataset) *corev1.Service {
	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      "kobenfs",
			Namespace: m.Namespace,
		},

		Spec: corev1.ServiceSpec{
			ClusterIP: "10.96.0.5",
			Selector: map[string]string{
				"role": "kobe-nfs-pod",
			},
			Ports: []corev1.ServicePort{
				{
					Name: "nfs",
					Port: int32(2049),
				},
				{
					Name: "mountd",
					Port: int32(20048),
				},
				{
					Name: "rpcbind",
					Port: int32(111),
				},
			},
		},
	}

	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

//NFS CONFIG
func (r *ReconcileKobeDataset) newNfsConfig(m *kobedatasetv1alpha1.KobeDataset) *corev1.ConfigMap {
	cmap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nfsconfig",
			Namespace: m.Namespace,
		},
		Data: map[string]string{"share": "/exports *(rw,fsid=0,insecure,no_root_squash)"},
	}
	controllerutil.SetControllerReference(m, cmap, r.scheme)
	return cmap
}

//NFS POD
func (r *ReconcileKobeDataset) newPodForNfs(m *kobedatasetv1alpha1.KobeDataset) *corev1.Pod {
	var priv bool
	priv = true

	env := corev1.EnvVar{Name: "SHARED_DIRECTORY", Value: "/kobe"}
	envs := []corev1.EnvVar{}
	envs = append(envs, env)

	volumemount := corev1.VolumeMount{Name: "nfs-disk", MountPath: "/exports"}
	volumemounts := []corev1.VolumeMount{}
	volumemounts = append(volumemounts, volumemount)

	volumemount1 := corev1.VolumeMount{Name: "nfs-server-config", MountPath: "etc/exports.d/"}
	volumemounts = append(volumemounts, volumemount1)

	volume := corev1.Volume{Name: "nfs-disk", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}}
	volumes := []corev1.Volume{}
	volumes = append(volumes, volume)

	volume1 := corev1.Volume{Name: "nfs-server-config", VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: "nfsconfig"}}}}
	volumes = append(volumes, volume1)

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      "kobenfs",
			Namespace: m.Namespace,
			Labels:    map[string]string{"role": "kobe-nfs-pod"},
		},
		Spec: corev1.PodSpec{

			Containers: []corev1.Container{{
				Image:           "alphayax/docker-volume-nfs:latest",
				Name:            "nfs-server-container",
				ImagePullPolicy: corev1.PullIfNotPresent,
				SecurityContext: &corev1.SecurityContext{
					Privileged: &priv,
				},
				Env: envs,
				Ports: []corev1.ContainerPort{
					{Name: "nfs", ContainerPort: int32(2049)},
					{Name: "mountd", ContainerPort: int32(20048)},
					{Name: "rpcbind", ContainerPort: int32(111)},
				},
				VolumeMounts: volumemounts,
			},
			},
			Volumes: volumes,
		},
	}
	controllerutil.SetControllerReference(m, pod, r.scheme)
	return pod
}

//PERSISTENT VOLUME
func (r *ReconcileKobeDataset) newPvForKobe(m *kobedatasetv1alpha1.KobeDataset, ip string) *corev1.PersistentVolume {
	//POSO MALAKAS EIMAI POU PREPEI NA PSAKSW TA API DEFINITIONS GIA AYTO...

	capacity := resource.MustParse("5Gi")
	rmap := corev1.ResourceList{}
	rmap["storage"] = capacity

	accessmodes := []corev1.PersistentVolumeAccessMode{"ReadWriteMany"}

	//nfs := &corev1.NFSVolumeSource{Server: "kobenfs." + m.Namespace + ".svc.cluster.local", Path: "/" + m.Name}
	//nfs := &corev1.NFSVolumeSource{Server: ip, Path: "/dumps/" + m.Name, ReadOnly: false}
	nfs := &corev1.NFSVolumeSource{Server: ip, Path: "/"}

	pv := &corev1.PersistentVolume{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolume",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "kobepv",
		},
		Spec: corev1.PersistentVolumeSpec{
			Capacity:                      rmap,
			AccessModes:                   accessmodes,
			PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimRetain,
			PersistentVolumeSource:        corev1.PersistentVolumeSource{NFS: nfs},
		},
	}
	controllerutil.SetControllerReference(m, pv, r.scheme)
	return pv
}

//PERSISTENT VOLUME CLAIM
func (r *ReconcileKobeDataset) newPvcForKobe(m *kobedatasetv1alpha1.KobeDataset) *corev1.PersistentVolumeClaim {
	s := ""
	accessmodes := []corev1.PersistentVolumeAccessMode{"ReadWriteMany"}
	capacity := resource.MustParse("4Gi")
	rmap := corev1.ResourceList{}
	rmap["storage"] = capacity
	pvc := &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kobepvc",
			Namespace: m.Namespace,
			//Annotations: map[string]string{"volume.beta.kubernetes.io/storage-class": ""}, -->EITHER THIS OR StorageClassNames : "" else it will create dynamic pv
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      accessmodes,
			Resources:        corev1.ResourceRequirements{Requests: rmap},
			StorageClassName: &s,
		},
	}

	controllerutil.SetControllerReference(m, pvc, r.scheme)
	return pvc
}
