package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

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
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Datasets []Dataset `json:"datasets"`
	Queries  []Query   `json:"queries"`
}

// KobeBenchmarkStatus defines the observed state of KobeBenchmark
// +k8s:openapi-gen=true
type KobeBenchmarkStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeBenchmark is the Schema for the kobebenchmarks API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
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

func init() {
	SchemeBuilder.Register(&KobeBenchmark{}, &KobeBenchmarkList{})
}
