package kobeutil

import (
	"context"

	kobeutilv1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobeutil/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_kobeutil")

/**
Utility controller for setting up things that make no sense belonging to any other resource
*/

// Add creates a new KobeUtil Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileKobeUtil{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("kobeutil-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource KobeUtil
	err = c.Watch(&source.Kind{Type: &kobeutilv1alpha1.KobeUtil{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobeutilv1alpha1.KobeUtil{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.PersistentVolume{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobeutilv1alpha1.KobeUtil{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.PersistentVolumeClaim{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobeutilv1alpha1.KobeUtil{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileKobeUtil implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileKobeUtil{}

// ReconcileKobeUtil reconciles a KobeUtil object
type ReconcileKobeUtil struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a KobeUtil object and makes changes based on the state read
// and what is in the KobeUtil.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileKobeUtil) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling KobeUtil")

	// Fetch the KobeUtil instance
	instance := &kobeutilv1alpha1.KobeUtil{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	//-------------------------------------------------checking for nfs config map health----------------------------------
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
		return reconcile.Result{Requeue: true}, nil

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
		return reconcile.Result{RequeueAfter: 10}, nil

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
		return reconcile.Result{RequeueAfter: 5}, nil
	}
	return reconcile.Result{}, nil
}

//-------------------functions that create native kubernetes that are controlled by utility controller-----------------------

//NFS SERVICE (its actually useless cause nfs service dns bug)
func (r *ReconcileKobeUtil) newServiceForNfs(m *kobeutilv1alpha1.KobeUtil) *corev1.Service {
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
func (r *ReconcileKobeUtil) newNfsConfig(m *kobeutilv1alpha1.KobeUtil) *corev1.ConfigMap {
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
func (r *ReconcileKobeUtil) newPodForNfs(m *kobeutilv1alpha1.KobeUtil) *corev1.Pod {
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

	volume1 := corev1.Volume{Name: "nfs-server-config",
		VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: "nfsconfig"}}}}
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
func (r *ReconcileKobeUtil) newPvForKobe(m *kobeutilv1alpha1.KobeUtil, ip string) *corev1.PersistentVolume {

	capacity := resource.MustParse("50Gi")
	rmap := corev1.ResourceList{}
	rmap["storage"] = capacity

	accessmodes := []corev1.PersistentVolumeAccessMode{"ReadWriteMany"}

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
func (r *ReconcileKobeUtil) newPvcForKobe(m *kobeutilv1alpha1.KobeUtil) *corev1.PersistentVolumeClaim {
	s := ""
	accessmodes := []corev1.PersistentVolumeAccessMode{"ReadWriteMany"}
	capacity := resource.MustParse("49Gi")
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
