package v1alpha1

import (
	kobefederatorv1alpha1 "github.com/semagrow/kobe/operator/pkg/apis/kobefederator/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KobeExperimentSpec defines the desired state of KobeExperiment
// +k8s:openapi-gen=true
type KobeExperimentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Benchmark     string                                `json:"benchmark"`
	Federators    []kobefederatorv1alpha1.KobeFederator `json:"federators"`
	RunFlag       bool                                  `json:"runFlag"`
	ClientImage   string                                `json:"clientImage"`
	ClientCommand []string                              `json:"clientCommands"`
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
