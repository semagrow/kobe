package v1alpha1

import (
	types "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Federator is helper struct
// +k8s:openapi-gen=true
type Federator struct {
	Name              string           `json:"name"`
	Image             string           `json:"image"`
	ImagePullPolicy   types.PullPolicy `json:"imagePullPolicy"`
	Affinity          types.Affinity   `json:"affinity"` //choose which nodes the fed likes to run in
	Port              int32            `json:"port"`
	ConfFromFileImage string           `json:"confFromFileImage"` //image that makes init file from dump or endpoint
	InputFileDir      string           `json:"inputFileDir"`      //where the above image expects the dump to be(if from dump)
	OutputFileDir     string           `json:"outputFileDir"`     //where the above image will place its result config file
	ConfImage         string           `json:"confImage"`
}

//Fed is something
type Fed struct {
	Name string `json:"name"`
}

// KobeExperimentSpec defines the desired state of KobeExperiment
// +k8s:openapi-gen=true
type KobeExperimentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Benchmark    string   `json:"benchmark"`
	Federator    []Fed    `jsons:"federator"`
	RunFlag      bool     `json:"runFlag"`
	TimesToRun   int      `json:"timesToRun"`
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

func init() {
	SchemeBuilder.Register(&KobeExperiment{}, &KobeExperimentList{})
}
