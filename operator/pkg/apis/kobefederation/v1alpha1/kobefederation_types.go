package v1alpha1

import (
	types "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KobeFederationSpec defines the desired state of KobeFederation
// +k8s:openapi-gen=true
type KobeFederationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Image             string           `json:"image"`
	ImagePullPolicy   types.PullPolicy `json:"imagePullPolicy"`
	Affinity          types.Affinity   `json:"affinity"` //choose which nodes the fed likes to run in
	Port              int32            `json:"port"`
	ConfFromFileImage string           `json:"confFromFileImage"` //image that makes init file from dump or endpoint
	InputFileDir      string           `json:"inputFileDir"`      //where the above image expects the dump to be(if from dump)
	OutputFileDir     string           `json:"outputFileDir"`     //where the above image will place its result config file
	ConfImage         string           `json:"confImage"`         //image that makes one init file from multiple init files
	InputDir          string           `json:"inputDir"`
	OutputDir         string           `json:"outputDir"`

	Endpoints    []string `json:"endpoints"`
	DatasetNames []string `json:"datasetNames"`
}

// KobeFederationStatus defines the observed state of KobeFederation
// +k8s:openapi-gen=true
type KobeFederationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	PodNames []string `json:"podnames"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeFederation is the Schema for the kobefederations API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
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

func init() {
	SchemeBuilder.Register(&KobeFederation{}, &KobeFederationList{})
}
