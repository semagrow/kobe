package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DatasetFile struct {
	URL      string `json:"url"`
	Checksum string `json:"checksum,omitempty"`
}

type Dataset struct {
	Name        string             `json:"name"`
	Files       []DatasetFile      `json:"files"`
	SystemSpec  *SystemDatasetSpec `json:"systemspec,omitempty"`
	TemplateRef string             `json:"templateRef,omitempty"` //  reference
	// If specified, the pod's scheduling constraints
	// +optional
	Affinity *v1.Affinity `json:"affinity,omitempty"`

	// Resources are not allowed for ephemeral containers. Ephemeral containers use spare resources
	// already allocated to the pod.
	// +optional
	Resources v1.ResourceRequirements `json:"resources,omitempty"`

	// network delays
	NetworkTopology     []NetworkConnection `json:"topology,omitempty"`
	FederatorConnection *NetworkConnection  `json:"FederatorDelay,omitempty"`
}

type Delay struct {
	// Add a fixed delay before forwarding the request. Format: 1h/1m/1s/1ms. MUST be >=1ms.
	FixedDelaySec *int64 `json:"fixedDelaySec"`

	// Add a fixed delay before forwarding the request. Format: 1h/1m/1s/1ms. MUST be >=1ms.
	FixedDelayMSec *int32 `json:"fixedDelayMSec"`

	// +optional
	Percentage *int32 `json:"percentage,omitempty"` // `protobuf:"fixed64,1,opt,name=value,proto3" json:"percentage,omitempty"`

	// +optional
	Percent *int32 `json:"percent,omitempty"`
}

type NetworkConnection struct {
	Source         *string `json:"datasetSource,omitempty"`
	DelayInjection Delay   `json:"delayInjection,omitempty"`
}

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
	Datasets []Dataset `json:"datasets"`
	Queries  []Query   `json:"queries"`
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
	Env             []v1.EnvVar   `json:"env"`
}

// ExperimentSpec defines the desired state of Experiment
// +k8s:openapi-gen=true
type ExperimentSpec struct {
	Benchmark            string         `json:"benchmark"`
	FederatorName        string         `json:"federatorName"`
	FederatorSpec        *FederatorSpec `json:"federatorSpec"`
	FederatorTemplateRef string         `json:"federatorTemplateRef"`
	Evaluator            Evaluator      `json:"evaluator"`
	TimesToRun           int            `json:"timesToRun"`
	RestartPolicy        RestartPolicy  `json:"restartPolicy,omitempty"`
	DryRun               bool           `json:"dryRun"`
	ForceNewInit         bool           `json:"forceNewInit"`
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

type FederationPhase string

const (
	FederationInitializing FederationPhase = "Initializing"
	FederationRunning      FederationPhase = "Running"
)

type InitializationPolicy string

const (
	ForceInit InitializationPolicy = "ForceInit"
)

type DatasetEndpoint struct {
	Host      string `json:"host"`
	Namespace string `json:"namespace"`
	Port      uint32 `json:"port"`
	Path      string `json:"path"`
}

// FederationSpec defines the desired state of Federation
// +k8s:openapi-gen=true
type FederationSpec struct {
	FederatorName   string               `json:"federatorName"`
	Template        FederatorSpec        `json:"spec"`
	Datasets        []DatasetEndpoint    `json:"datasets"` // use v1.LocalObjectReference ?
	NetworkTopology []NetworkConnection  `json:"topology,omitempty"`
	InitPolicy      InitializationPolicy `json:"initPolicy,omitempty"`
}

// FederationStatus defines the observed state of KobeFederation
// +k8s:openapi-gen=true
type FederationStatus struct {
	PodNames []string        `json:"podNames"`
	Phase    FederationPhase `json:"phase"`
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

// FederatorTemplate defines a federator and its components that it needs to be installed.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:path=federatortemplates,scope=Namespaced
// +k8s:openapi-gen=true
type FederatorTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              FederatorSpec `json:"spec,inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type FederationTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FederatorTemplate `json:"items"`
}

// FederatorSpec contains all necessary information for a federator
type FederatorSpec struct {
	InitContainers []v1.Container `json:"initContainers,omitempty"`
	Containers     []v1.Container `json:"containers"`

	// Number of port to expose on the host.
	// If specified, this must be a valid port number, 0 < x < 65536.
	Port int32 `json:"port"`

	//suffix to be added to endpoint of federator f.e ../SemaGrow/sparql
	// +optional
	Path string `json:"path"`

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
	Spec              FederatorSpec `json:"spec,omitempty"`
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:resource:path=datasettemplates,scope=Namespaced
// +k8s:openapi-gen=true
type DatasetTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SystemDatasetSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FederatorList contains a list of Federator
type DatasetTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DatasetTemplate `json:"items"`
}

// type DatasetTemplateSpec struct {
// 	metav1.ObjectMeta `json:"metadata,omitempty"`
// 	Spec              DatasetSpec `json:"spec,omitempty"`
// }

// DatasetSpec defines the desired state of Dataset
type SystemDatasetSpec struct {
	ImportContainers []v1.Container `json:"importContainers,omitempty"`

	//List of initialization containers belonging to the pod. Init containers
	//are executed in order prior to containers being started. If any init
	//container fails, the pod is considered to have failed and is handled
	//according to its restartPolicy. The name for an init container or normal
	//container must be unique among all containers. Init containers may not
	//have Lifecycle actions, Readiness probes, or Liveness probes. The
	//resourceRequirements of an init container are taken into account during
	//scheduling by finding the highest request/limit for each resource type,
	//and then using the max of of that value or the sum of the normal
	//containers. Limits are applied to init containers in a similar fashion.
	//Init containers cannot currently be added or removed. Cannot be updated.
	//More info:
	//https://kubernetes.io/docs/concepts/workloads/pods/init-containers/
	InitContainers []v1.Container `json:"initContainers,omitempty"`

	// List of containers belonging to the pod. Containers cannot currently be
	// added or removed. There must be at least one container in a Pod. Cannot
	// be updated.
	Containers []v1.Container `json:"containers"`

	// Forces to download and load from dataset file
	InitPolicy InitializationPolicy `json:"initPolicy"`

	// Number of port to expose on the host.
	// If specified, this must be a valid port number, 0 < x < 65536.
	Port uint32 `json:"port"`

	// Path that the container will listen for queries
	Path string `json:"path"`
}

// DatasetStatus defines the observed state of Dataset
// +k8s:openapi-gen=true
type EphemeralDatasetStatus struct {
	PodNames  []string `json:"podNames"`
	Phase     string   `json:"phase"`
	ForceLoad bool     `json:"forceLoad"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EphemeralDataset is the Schema for the kobedatasets API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=ephemeraldatasets,scope=Namespaced
type EphemeralDataset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Dataset                `json:"spec,omitempty"`
	Status EphemeralDatasetStatus `json:"status,omitempty"`
}

// SetDefaults sets the defaults of the KobeDatasetSpec
func (r *EphemeralDataset) SetDefaults() bool {
	changed := false
	// rs := r.Spec.Template.TemplateSpec.Path
	// if rs == "" {
	// 	r.Spec.Template.TemplateSpec.Path = "/sparql"
	// 	changed = true
	// }
	return changed
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EphemeralDatasetList contains a list of EphemeralDataset
type EphemeralDatasetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EphemeralDataset `json:"items"`
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
	SchemeBuilder.Register(&Benchmark{}, &BenchmarkList{})
	SchemeBuilder.Register(&Experiment{}, &ExperimentList{})
	SchemeBuilder.Register(&Federation{}, &FederationList{})
	SchemeBuilder.Register(&Federator{}, &FederatorList{})
	SchemeBuilder.Register(&EphemeralDataset{}, &EphemeralDatasetList{})
	SchemeBuilder.Register(&KobeUtil{}, &KobeUtilList{})
	SchemeBuilder.Register(&DatasetTemplate{}, &DatasetTemplateList{})
	SchemeBuilder.Register(&FederatorTemplate{}, &FederationTemplateList{})
}
