package v1alpha1

import (
	types "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KobeDatasetSpec defines the desired state of KobeDataset
// +k8s:openapi-gen=true
type KobeDatasetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Image           string           `json:"image"`
	ToDownload      bool             `json:"toDownload"`
	DumpDirectory   string           `json:"dumbDirectory"`
	DownloadFrom    string           `json:"downloadFrom"`
	ImagePullPolicy types.PullPolicy `json:"imagePullPolicy"`
	Count           int32            `json:"count"`
	Group           string           `json:"group"`
	Port            int32            `json:"port"`
}

// KobeDatasetStatus defines the observed state of KobeDataset
// +k8s:openapi-gen=true
type KobeDatasetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	PodNames []string `json:"podnames"`
	AppGroup string   `json:"appgroup"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeDataset is the Schema for the kobedatasets API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type KobeDataset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KobeDatasetSpec   `json:"spec,omitempty"`
	Status KobeDatasetStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeDatasetList contains a list of KobeDataset
type KobeDatasetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KobeDataset `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KobeDataset{}, &KobeDatasetList{})
}
