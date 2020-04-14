package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

//KobeBenchmarkSpec defines the components for this benchmark setup
type KobeBenchmarkSpec struct {
	Datasets []Dataset `json:"datasets"`
	Queries  []Query   `json:"queries"`
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

// KobeBenchmarkStatus defines the observed state of KobeBenchmark
// +k8s:openapi-gen=true
type KobeBenchmarkStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeBenchmarkList contains a list of KobeBenchmark
type KobeBenchmarkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KobeBenchmark `json:"items"`
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

type ExperimentPhase string

// possible experiment phases
const (
	ExperimentNotStarted   ExperimentPhase = "NotStarted"
	ExperimentInitializing ExperimentPhase = "Initializing"
	ExperimentRunning      ExperimentPhase = "Running"
	ExperimentCompleted    ExperimentPhase = "Completed"
	ExperimentFailed       ExperimentPhase = "Failed"
)

type RestartPolicy string

// possible restart policies
const (
	RestartNever  RestartPolicy = "Never"
	RestartAlways RestartPolicy = "Always"
)

// Evaluator defines the
type Evaluator struct {
	Image           string        `json:"image"`
	ImagePullPolicy v1.PullPolicy `json:"imagePullPolicy"`

	Command []string `json:"command"`

	Parallelism int32 `json:"parallelism"`
}

// KobeExperimentSpec defines the desired state of KobeExperiment
// +k8s:openapi-gen=true
type KobeExperimentSpec struct {
	Benchmark string `json:"benchmark"`

	Federator string `json:"federator"`

	DryRun bool `json:"dryRun"`

	TimesToRun int `json:"timesToRun"`

	ForceNewInit bool `json:"forceNewInit"`

	RestartPolicy RestartPolicy `json:"restartPolicy,omitempty"`

	Evaluator Evaluator `json:"evaluator"`
}

// KobeExperimentStatus defines the observed state of KobeExperiment
// +k8s:openapi-gen=true
type KobeExperimentStatus struct {

	// Time at which this workflow started
	StartTime metav1.Time `json:"startTime,omitempty"`

	// Time at which this workflow completed
	CompletionTime metav1.Time `json:"completionTime,omitempty"`

	// The current iteration of the experiment
	// It should be zero if not started yet
	CurrentRun int `json:"run,omitempty"`

	// The phase of the experiment
	Phase ExperimentPhase `json:"phase"`
}

// KobeFederationSpec defines the desired state of KobeFederation
// +k8s:openapi-gen=true
type KobeFederationSpec struct {
	FederatorName string `json:"federatorName"`

	Template FederatorTemplate `json:"template"`

	Endpoints []string `json:"endpoints"`
	Datasets  []string `json:"datasets"`

	ForceNewInit bool `json:"forceNewInit"`
	Init         bool `json:"init"`
}

// SetDefaults set the defaults of a federation
func (r *KobeFederation) SetDefaults() bool {
	changed := false
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

// KobeFederatorSpec defines the desired state of KobeFederator
// +k8s:openapi-gen=true
type KobeFederatorSpec struct {
	FederatorTemplate `json:",inline"`
}

// FederatorTemplate defines
type FederatorTemplate struct {
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

	// Number of port to expose on the host.
	// If specified, this must be a valid port number, 0 < x < 65536.
	Port int32 `json:"port"`

	//suffix to be added to endpoint of federator f.e ../SemaGrow/sparql
	// +optional
	Path string `json:"path"`

	// If specified, the pod's scheduling constraints
	// +optional
	Affinity *v1.Affinity `json:"affinity,omitempty"`

	// Resources are not allowed for ephemeral containers. Ephemeral containers use spare resources
	// already allocated to the pod.
	// +optional
	Resources v1.ResourceRequirements `json:"resources,omitempty"`

	// The Docker image that receives a compressed dataset and may produce
	// configuration needed from the federator to federate this specific dataset
	// This container will run one time for each dataset in the federation
	ConfFromFileImage string `json:"confFromFileImage"`

	// where the above image expects the dump to be (if from dump)
	InputDumpDir string `json:"inputDumpDir"`

	// where the above image will place its result config file
	OutputDumpDir string `json:"outputDumpDir"`

	// The Docker image that initializes the federator (equivalent to initContainers)
	ConfImage string `json:"confImage"`

	InputDir string `json:"inputDir"`

	OutputDir string `json:"outputDir"`

	//which directory the federator needs the metadata config files in order to find them
	FedConfDir string `json:"fedConfDir"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeFederator is the Schema for the kobefederators API
// +k8s:openapi-gen=true
// +kubebuilder:resource:path=kobefederators,scope=Namespaced
type KobeFederator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec KobeFederatorSpec `json:"spec,omitempty"`
}

// SetDefaults set the defaults of a federation
func (r *KobeFederator) SetDefaults() bool {
	changed := false
	rs := &r.Spec
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeFederatorList contains a list of KobeFederator
type KobeFederatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KobeFederator `json:"items"`
}

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

// SetDefaults sets the defaults of the KobeDatasetSpec
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeUtil is the Schema for the kobeutils API
// +k8s:openapi-gen=true
// +kubebuilder:resource:path=kobeutils,scope=Namespaced
type KobeUtil struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
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
