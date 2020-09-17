package kobeutil

import (
	"context"
	//"strings"

	api "github.com/semagrow/kobe/operator/pkg/apis/kobe/v1alpha1"
	kobev1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobe/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	//"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
	err = c.Watch(&source.Kind{Type: &kobev1alpha1.KobeUtil{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobev1alpha1.KobeUtil{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.PersistentVolume{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobev1alpha1.KobeUtil{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.PersistentVolumeClaim{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kobev1alpha1.KobeUtil{},
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
	instance := &kobev1alpha1.KobeUtil{}
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
		err := r.client.Create(context.TODO(), pod)
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
		err := r.client.Create(context.TODO(), pv)
		if err != nil {
			reqLogger.Info("Failed to create the persistent volume that the datasets will use to retain their data if they shutdown and restarted")
			return reconcile.Result{}, err
		}
		return reconcile.Result{RequeueAfter: 1000000000}, nil
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
		return reconcile.Result{RequeueAfter: 1000000000}, nil
	}
	//-------------------------------------------------checking for ftp config map health----------------------------------
	ftpconfigFound := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "sftp-keys", Namespace: instance.Namespace}, ftpconfigFound)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating config map for nfs\n") //used by the dataset pods to store the db file and regain it if that dataset is restarted
		conf := r.newFtpConfig(instance)
		err = r.client.Create(context.TODO(), conf)
		if err != nil {
			reqLogger.Info("Failed to create the nfs config: %v\n", err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}
	//------------------------------------------------checking for ftp service health---------------------------------------------------
	ftpServiceFound := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "kobe-ftp", Namespace: corev1.NamespaceDefault}, ftpServiceFound)
	if err != nil && errors.IsNotFound(err) {
		service := r.newServiceForFtp(instance)
		err := r.client.Create(context.TODO(), service)
		if err != nil {
			reqLogger.Info("Failed to create the kobe ftp service: %v\n", err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{RequeueAfter: 1000000000}, nil
	}
	//------------------------------------------------checking for ftp deployment health---------------------------------------------------
	ftpDeploymentFound := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "kobe-ftp", Namespace: corev1.NamespaceDefault}, ftpDeploymentFound)
	if err != nil && errors.IsNotFound(err) {
		deployment := r.newDeploymentForFtp(instance)
		err := r.client.Create(context.TODO(), deployment)
		if err != nil {
			reqLogger.Info("Failed to create the kobe ftp deployment: %v\n", err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{RequeueAfter: 1000000000}, nil
	}
	return reconcile.Result{}, nil
}

//-------------------functions that create native kubernetes that are controlled by utility controller-----------------------

//NFS SERVICE (its actually useless cause nfs service dns bug)
func (r *ReconcileKobeUtil) newServiceForNfs(m *kobev1alpha1.KobeUtil) *corev1.Service {
	service := &corev1.Service{
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
func (r *ReconcileKobeUtil) newNfsConfig(m *kobev1alpha1.KobeUtil) *corev1.ConfigMap {
	cmap := &corev1.ConfigMap{
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
func (r *ReconcileKobeUtil) newPodForNfs(m *kobev1alpha1.KobeUtil) *corev1.Pod {
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
func (r *ReconcileKobeUtil) newPvForKobe(m *kobev1alpha1.KobeUtil, ip string) *corev1.PersistentVolume {

	capacity := resource.MustParse("50Gi")
	rmap := corev1.ResourceList{}
	rmap["storage"] = capacity

	accessmodes := []corev1.PersistentVolumeAccessMode{"ReadWriteMany"}

	nfs := &corev1.NFSVolumeSource{Server: ip, Path: "/"}

	pv := &corev1.PersistentVolume{
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
func (r *ReconcileKobeUtil) newPvcForKobe(m *kobev1alpha1.KobeUtil) *corev1.PersistentVolumeClaim {
	s := ""
	accessmodes := []corev1.PersistentVolumeAccessMode{"ReadWriteMany"}
	capacity := resource.MustParse("49Gi")
	rmap := corev1.ResourceList{}
	rmap["storage"] = capacity
	pvc := &corev1.PersistentVolumeClaim{
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

//ftp key configmap
func (r *ReconcileKobeUtil) newFtpConfig(m *kobev1alpha1.KobeUtil) *corev1.ConfigMap {
	cmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sftp-keys",
			Namespace: m.Namespace,
		},
		Data: map[string]string{"ssh_host_ed25519_key": "-----BEGIN OPENSSH PRIVATE KEY-----\n" +
			"b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW\n" +
			"QyNTUxOQAAACD82+V8lySp8+tu6aVvu+wPJxgPoKFO2nVjuOzxPQNU4AAAAJgNUkqWDVJK\n" +
			"lgAAAAtzc2gtZWQyNTUxOQAAACD82+V8lySp8+tu6aVvu+wPJxgPoKFO2nVjuOzxPQNU4A\n" +
			"AAAEAZE2tJIq4FQPKasQH7llqjdr36b04WtppMq71G3FZWXfzb5XyXJKnz627ppW+77A8n\n" +
			"GA+goU7adWO47PE9A1TgAAAAEW5jc3JAbmNzci1NUy03ODIzAQIDBA==\n" +
			"-----END OPENSSH PRIVATE KEY-----",

			"ssh_host_ed25519_key.pub": "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIPzb5XyXJKnz627ppW+77A8nGA+goU7adWO47PE9A1Tg ncsr@ncsr-MS-7823",

			"ssh_host_rsa_key": "-----BEGIN RSA PRIVATE KEY-----\n" +
				"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\n" +
				"-----END RSA PRIVATE KEY-----",

			"ssh_host_rsa_key.pub": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDGKTHRYuce7X047pPLKg6Lf48IZ7uXTD32yUWrArWQdCxjjTSsqcEGhNUtwXP/YMNrOjL8yhiQUexq3f8kCJvYN0NevvhbMMAvm8o6P3WTmDCEHPB2YcPoKIhDQXqQF+SUJpTzekMh3P5HdE3ft6BC3a28XH12pQvLAk+XFQK75plD5cFyXuThUmfDuL27IQPkMqPoTkDCCOwCUrhTZ2nXrc59zOe9IGfnOGYcIdZ0Kwd1ESZLX1hGhkVMDJQzn17ftmvO9ConDHlxJZo4WGrElYaZ/PCUS2/WCAC2zaZy0oCDCoj5nhxb14DoUF68xx7ADhNKCqwJPvUxJhKa70CiXV5ySSfviqXz8oUaRSqeoJ81I+HI6/OU++cuDCTeUmJ5H9Krh5bNKosimZ2M3vDHcYejbYCBhL32YVWPk5ZHA45fsns0bkKndVefwuCXhEBr383sEacIA6ld8dCSi5RboiWfTtgtfIMFAgGepyhS9ntcUyjBYlXhOKMPtEJd4KEoeXOGCqZ19S3TkU9U+z8r5MAJCMtrkN54DJl1ftUtVaVwRgUvOcerZeFZaNyjd1ivoGz7UaU3R3uS/N81/eRD/mdyYaFWII5L8bfUoF5Oy35TIlhi2TTqwUays2JWeX73ZNus79gZ2Iuscw5V0+DJZrQYcUFa5qRaBNgfl3IcIw== ncsr@ncsr-MS-7823",
		},
	}
	controllerutil.SetControllerReference(m, cmap, r.scheme)
	return cmap
}

//create new deployment for ftp server
func (r *ReconcileKobeUtil) newDeploymentForFtp(m *api.KobeUtil) *appsv1.Deployment {
	replicas := int32(1)

	volume := corev1.Volume{Name: "data", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}}
	volumes := []corev1.Volume{}
	volumes = append(volumes, volume)

	volume1 := corev1.Volume{Name: "sftp-keys",
		VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: "sftp-public-keys"}}}}
	volumes = append(volumes, volume1)

	// volumemount := corev1.VolumeMount{Name: "data", MountPath: "/kobe/dataset"}
	// volumemounts := []corev1.VolumeMount{}
	// volumemounts = append(volumemounts, volumemount)
	volumemount := corev1.VolumeMount{Name: "sftp-public-keys", MountPath: "/home/myUser/.ssh/keys"}
	volumemounts := []corev1.VolumeMount{}
	volumemounts = append(volumemounts, volumemount)
	//maybe create a secret
	// env := corev1.EnvVar{Name: "PASSWORD",
	// 	ValueFrom: &corev1.EnvVarSource{
	// 		SecretKeyRef: &corev1.SecretKeySelector{
	// 			LocalObjectReference: corev1.LocalObjectReference{Name: "sftp-server-sec"},
	// 			Key:                  "password"},
	// 	},
	// }
	env := corev1.EnvVar{Name: "PASSWORD", Value: "kobe"}
	envs := []corev1.EnvVar{}
	envs = append(envs, env)

	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kobe-ftp",
			Namespace: corev1.NamespaceDefault,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "Kobe-Operator", "utility": "ftp-server"}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "Kobe-Operator", "utility": "ftp-server"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           "atmoz/sftp:latest",
						Name:            "kobe-ftp",
						ImagePullPolicy: corev1.PullAlways,
						Ports: []corev1.ContainerPort{{
							ContainerPort: int32(22),
							Name:          m.Name,
						}},
						VolumeMounts: volumemounts,
						Env:          envs,
						Args:         []string{"myUser::1001:100:incoming,outgoing"},
						SecurityContext: &corev1.SecurityContext{
							Capabilities: &corev1.Capabilities{Add: []corev1.Capability{"SYS_ADMIN"}},
						},
					}},
					Volumes: volumes,
				},
			},
		},
	}
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep

}

func (r *ReconcileKobeUtil) newServiceForFtp(m *kobev1alpha1.KobeUtil) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kobe-ftp",
			Namespace: corev1.NamespaceDefault,
		},

		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":     "Kobe-Operator",
				"utility": "ftp-server",
			},
			Ports: []corev1.ServicePort{
				{
					Name: "ssh",
					Port: int32(22),
					//TargetPort: intstr.FromString("22"),
				},
			},
		},
	}

	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}
