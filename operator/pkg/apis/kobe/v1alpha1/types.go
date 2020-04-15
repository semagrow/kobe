package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Benchmark is the Schema for the benchmarks API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=benchmarks,scope=Namespaced
type Benchmark struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              BenchmarkSpec   `json:"spec,omitempty"`
	Status            BenchmarkStatus `json:"status,omitempty"`
}

// BenchmarkSpec defines the components for this benchmark setup
type BenchmarkSpec struct {
	Datasets []string `json:"datasets"`
	Queries  []Query  `json:"queries"`
}

//Query contains the query info
type Query struct {
	Name        string `json:"name"`
	Language    string `json:"language"`
	QueryString string `json:"queryString"`
}

// BenchmarkStatus defines the observed state of Benchmark
// +k8s:openapi-gen=true
type BenchmarkStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BenchmarkList contains a list of Benchmark
type BenchmarkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Benchmark `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Experiment is the Schema for the experiments API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=experiments,scope=Namespaced
type Experiment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ExperimentSpec   `json:"spec,omitempty"`
	Status            ExperimentStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExperimentList contains a list of Experiment
type ExperimentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Experiment `json:"items"`
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
	Command         []string      `json:"command"`
	Parallelism     int32         `json:"parallelism"`
}

// ExperimentSpec defines the desired state of Experiment
// +k8s:openapi-gen=true
type ExperimentSpec struct {
	Benchmark     string        `json:"benchmark"`
	Federator     string        `json:"federator"`
	Evaluator     Evaluator     `json:"evaluator"`
	TimesToRun    int           `json:"timesToRun"`
	RestartPolicy RestartPolicy `json:"restartPolicy,omitempty"`
	DryRun        bool          `json:"dryRun"`
	ForceNewInit  bool          `json:"forceNewInit"`
}

// ExperimentStatus defines the observed state of Experiment
// +k8s:openapi-gen=true
type ExperimentStatus struct {

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

// FederationSpec defines the desired state of Federation
// +k8s:openapi-gen=true
type FederationSpec struct {
	FederatorName string            `json:"federatorName"`
	Template      FederatorTemplate `json:"template"`
	Endpoints     []string          `json:"endpoints"`
	Datasets      []string          `json:"datasets"` // use v1.LocalObjectReference ?
	ForceNewInit  bool              `json:"forceNewInit"`
	Init          bool              `json:"init"`
}

// FederationStatus defines the observed state of KobeFederation
// +k8s:openapi-gen=true
type FederationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	PodNames []string `json:"podNames"`
	Phase    string   `json:"phase"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Federation is the Schema for the federations API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=federations,scope=Namespaced
type Federation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              FederationSpec   `json:"spec,omitempty"`
	Status            FederationStatus `json:"status,omitempty"`
}

// SetDefaults set the defaults of a federation
func (r *Federation) SetDefaults() bool {
	changed := false
	return changed
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FederationList contains a list of Federation
type FederationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Federation `json:"items"`
}

// FederatorSpec defines the desired state of Federator
// +k8s:openapi-gen=true
type FederatorSpec struct {
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

// Federator is the Schema for the federators API
// +k8s:openapi-gen=true
// +kubebuilder:resource:path=federators,scope=Namespaced
type Federator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec FederatorSpec `json:"spec,omitempty"`
}

// SetDefaults set the defaults of a federation
func (r *Federator) SetDefaults() bool {
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

// FederatorList contains a list of Federator
type FederatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Federator `json:"items"`
}

type DatasetInitializationPolicy string

const (
	ForceLoad DatasetInitializationPolicy = "ForceLoad"
)

// DatasetSpec defines the desired state of Dataset
// +k8s:openapi-gen=true
type DatasetSpec struct {

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

// DatasetStatus defines the observed state of Dataset
// +k8s:openapi-gen=true
type DatasetStatus struct {
	PodNames []string `json:"podNames"`

	Phase string `json:"phase"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Dataset is the Schema for the kobedatasets API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=datasets,scope=Namespaced
type Dataset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatasetSpec   `json:"spec,omitempty"`
	Status DatasetStatus `json:"status,omitempty"`
}

// SetDefaults sets the defaults of the KobeDatasetSpec
func (r *Dataset) SetDefaults() bool {
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

// DatasetList contains a list of KobeDataset
type DatasetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Dataset `json:"items"`
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
	SchemeBuilder.Register(&Dataset{}, &DatasetList{})
	SchemeBuilder.Register(&Benchmark{}, &BenchmarkList{})
	SchemeBuilder.Register(&Experiment{}, &ExperimentList{})
	SchemeBuilder.Register(&Federation{}, &FederationList{})
	SchemeBuilder.Register(&Federator{}, &FederatorList{})
	SchemeBuilder.Register(&KobeUtil{}, &KobeUtilList{})
}
