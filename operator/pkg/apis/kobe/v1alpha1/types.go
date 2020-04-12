package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KobeDatasetSpec defines the desired state of KobeDataset
// +k8s:openapi-gen=true
type KobeDatasetSpec struct {

	// Docker image name.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	// +optional
	Image string `json:"image"`

	// Image pull policy.
	// One of Always, Never, IfNotPresent.
	// Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/containers/images#updating-images
	// +optional
	ImagePullPolicy v1.PullPolicy `json:"imagePullPolicy"`

	// Replicas is the number of desired replicas.
	// This is a pointer to distinguish between explicit zero and unspecified.
	// Defaults to 1.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller#what-is-a-replicationcontroller
	// +optional
	Replicas *int32 `json:"replicas"`

	// Forces to download and load from dataset file
	ForceLoad bool `json:"forceLoad"`

	// A URL that points to the compressed dataset file
	DownloadFrom string `json:"downloadFrom"`

	// Number of port to expose on the host.
	// If specified, this must be a valid port number, 0 < x < 65536.
	Port int32 `json:"port"`

	// Path that the container will listen for queries
	Path string `json:"path"`

	// List of environment variables to set in the container.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	Env []v1.EnvVar `json:"env,omitempty"`

	// If specified, the pod's scheduling constraints
	// +optional
	Affinity *v1.Affinity `json:"affinity,omitempty"`

	// Resources are not allowed for ephemeral containers. Ephemeral containers use spare resources
	// already allocated to the pod.
	// +optional
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
}

// KobeDatasetStatus defines the observed state of KobeDataset
// +k8s:openapi-gen=true
type KobeDatasetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	PodNames []string `json:"podNames"`
	AppGroup string   `json:"appGroup"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeDataset is the Schema for the kobedatasets API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kobedatasets,scope=Namespaced
type KobeDataset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KobeDatasetSpec   `json:"spec,omitempty"`
	Status KobeDatasetStatus `json:"status,omitempty"`
}

func (r *KobeDataset) SetDefaults() bool {
	changed := false
	rs := &r.Spec
	if rs.Replicas == nil {
		var replicas int32 = 1
		rs.Replicas = &replicas
		changed = true
	}
	if rs.Path == "" {
		rs.Path = "/sparql"
		changed = true
	}
	return changed
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeDatasetList contains a list of KobeDataset
type KobeDatasetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KobeDataset `json:"items"`
}

// Dataset contains the dataset info
type Dataset struct {
	Name         string `json:"name"`
	Image        string `json:"image"`
	DownloadFrom string `json:"downloadFrom"`
}

//Query contains the query info
type Query struct {
	Name        string `json:"name"`
	Language    string `json:"language"`
	QueryString string `json:"queryString"`
}

//KobeBenchmarkSpec defines the components for this benchmark setup
type KobeBenchmarkSpec struct {
	Datasets []Dataset `json:"datasets"`
	Queries  []Query   `json:"queries"`
}

// KobeBenchmarkStatus defines the observed state of KobeBenchmark
// +k8s:openapi-gen=true
type KobeBenchmarkStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeBenchmark is the Schema for the kobebenchmarks API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kobebenchmarks,scope=Namespaced
type KobeBenchmark struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KobeBenchmarkSpec   `json:"spec,omitempty"`
	Status KobeBenchmarkStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeBenchmarkList contains a list of KobeBenchmark
type KobeBenchmarkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KobeBenchmark `json:"items"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Federator is helper struct
// +k8s:openapi-gen=true
type Federator struct {
	Name              string        `json:"name"`
	Image             string        `json:"image"`
	ImagePullPolicy   v1.PullPolicy `json:"imagePullPolicy"`
	Affinity          *v1.Affinity  `json:"affinity"` //choose which nodes the fed likes to run in
	Port              int32         `json:"port"`
	ConfFromFileImage string        `json:"confFromFileImage"` //image that makes init file from dump or endpoint
	InputFileDir      string        `json:"inputFileDir"`      //where the above image expects the dump to be(if from dump)
	OutputFileDir     string        `json:"outputFileDir"`     //where the above image will place its result config file
	ConfImage         string        `json:"confImage"`
}

// KobeExperimentSpec defines the desired state of KobeExperiment
// +k8s:openapi-gen=true
type KobeExperimentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Benchmark    string `json:"benchmark"`
	Federator    string `json:"federator"`
	DryRun       bool   `json:"dryRun"`
	TimesToRun   int    `json:"timesToRun"`
	ForceNewInit bool   `json:"forceNewInit"`

	EvalImage    string   `json:"evalImage"`
	EvalCommands []string `json:"evalCommands"`
}

// KobeExperimentStatus defines the observed state of KobeExperiment
// +k8s:openapi-gen=true
type KobeExperimentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeExperiment is the Schema for the kobeexperiments API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
//+kubebuilder:resource:path=kobeexperiments,scope=Namespaced
type KobeExperiment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KobeExperimentSpec   `json:"spec,omitempty"`
	Status KobeExperimentStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeExperimentList contains a list of KobeExperiment
type KobeExperimentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KobeExperiment `json:"items"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KobeFederationSpec defines the desired state of KobeFederation
// +k8s:openapi-gen=true
type KobeFederationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Image           string        `json:"image"`
	ImagePullPolicy v1.PullPolicy `json:"imagePullPolicy"`

	// If specified, the pod's scheduling constraints
	// +optional
	Affinity *v1.Affinity `json:"affinity,omitempty"`

	// Resources are not allowed for ephemeral containers. Ephemeral containers use spare resources
	// already allocated to the pod.
	// +optional
	Resources v1.ResourceRequirements `json:"resources,omitempty"`

	Port int32  `json:"port"`
	Path string `json:"path"` //suffix to be added to endpoint of federator f.e ../SemaGrow/sparql

	ConfFromFileImage string `json:"confFromFileImage"` //image that makes init file from dump or endpoint
	InputDumpDir      string `json:"inputDumpDir"`      //where the above image expects the dump to be(if from dump)
	OutputDumpDir     string `json:"outputDumpDir"`     //where the above image will place its result config file
	ConfImage         string `json:"confImage"`         //image that makes one init file from multiple init files
	InputDir          string `json:"inputDir"`
	OutputDir         string `json:"outputDir"`

	FedConfDir    string   `json:"fedConfDir"` //which directory the federator needs the metadata config files in order to find them
	ForceNewInit  bool     `json:"forceNewInit"`
	Init          bool     `json:"init"`
	FederatorName string   `json:"federatorName"`
	Endpoints     []string `json:"endpoints"`
	DatasetNames  []string `json:"datasetNames"`
}

func (r *KobeFederation) SetDefaults() bool {
	changed := false
	rs := &r.Spec;
	if rs.InputDumpDir == "" {
		rs.InputDumpDir = "/kobe/input"
		changed = true
	}
	if rs.OutputDumpDir == "" {
		rs.OutputDumpDir = "/kobe/output"
		changed = true
	}
	if rs.InputDir == "" {
		rs.InputDir = "/kobe/input"
		changed = true
	}
	if rs.OutputDir == "" {
		rs.OutputDir = "/kobe/output"
		changed = true
	}
	return changed
}

// KobeFederationStatus defines the observed state of KobeFederation
// +k8s:openapi-gen=true
type KobeFederationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	PodNames []string `json:"podNames"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeFederation is the Schema for the kobefederations API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kobefederations,scope=Namespaced
type KobeFederation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KobeFederationSpec   `json:"spec,omitempty"`
	Status KobeFederationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeFederationList contains a list of KobeFederation
type KobeFederationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KobeFederation `json:"items"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KobeFederatorSpec defines the desired state of KobeFederator
// +k8s:openapi-gen=true
type KobeFederatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	//InitContainers    []corev1.Container `json:"initContainer"` //obsolete
	Image           string        `json:"image"`
	ImagePullPolicy v1.PullPolicy `json:"imagePullPolicy"`

	// If specified, the pod's scheduling constraints
	// +optional
	Affinity *v1.Affinity `json:"affinity,omitempty"`

	// Resources are not allowed for ephemeral containers. Ephemeral containers use spare resources
	// already allocated to the pod.
	// +optional
	Resources v1.ResourceRequirements `json:"resources,omitempty"`

	Port int32  `json:"port"`
	Path string `json:"path"` //suffix to be added to endpoint of federator f.e ../SemaGrow/sparql

	ConfFromFileImage string `json:"confFromFileImage"` //image that makes init file from dump or endpoint
	InputDumpDir      string `json:"inputDumpDir"`      //where the above image expects the dump to be(if from dump)
	OutputDumpDir     string `json:"outputDumpDir"`     //where the above image will place its result config file
	ConfImage         string `json:"confImage"`         //image that makes one init file from multiple init files
	InputDir          string `json:"inputDir"`
	OutputDir         string `json:"outputDir"`
	FedConfDir        string `json:"fedConfDir"` //which directory the federator needs the metadata config files in order to find them
}

// KobeFederatorStatus defines the observed state of KobeFederator
// +k8s:openapi-gen=true
type KobeFederatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	PodNames []string `json:"podNames"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeFederator is the Schema for the kobefederators API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kobefederators,scope=Namespaced
type KobeFederator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KobeFederatorSpec   `json:"spec,omitempty"`
	Status KobeFederatorStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeFederatorList contains a list of KobeFederator
type KobeFederatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KobeFederator `json:"items"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KobeUtilSpec defines the desired state of KobeUtil
// +k8s:openapi-gen=true
type KobeUtilSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// KobeUtilStatus defines the observed state of KobeUtil
// +k8s:openapi-gen=true
type KobeUtilStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeUtil is the Schema for the kobeutils API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kobeutils,scope=Namespaced
type KobeUtil struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KobeUtilSpec   `json:"spec,omitempty"`
	Status KobeUtilStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeUtilList contains a list of KobeUtil
type KobeUtilList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KobeUtil `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KobeDataset{}, &KobeDatasetList{})
	SchemeBuilder.Register(&KobeBenchmark{}, &KobeBenchmarkList{})
	SchemeBuilder.Register(&KobeExperiment{}, &KobeExperimentList{})
	SchemeBuilder.Register(&KobeFederation{}, &KobeFederationList{})
	SchemeBuilder.Register(&KobeFederator{}, &KobeFederatorList{})
	SchemeBuilder.Register(&KobeUtil{}, &KobeUtilList{})
}
