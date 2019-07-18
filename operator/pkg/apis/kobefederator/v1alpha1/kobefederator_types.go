package v1alpha1

import (
	types "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KobeFederatorSpec defines the desired state of KobeFederator
// +k8s:openapi-gen=true
type KobeFederatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	//InitContainers    []corev1.Container `json:"initContainer"` //obsolete
	Image             string           `json:"image"`
	ImagePullPolicy   types.PullPolicy `json:"imagePullPolicy"`
	Affinity          types.Affinity   `json:"affinity"` //choose which nodes the fed likes to run in
	Port              int32            `json:"port"`
	ConfFromFileImage string           `json:"confFromFileImage"` //image that makes init file from dump or endpoint
	InputFileDir      string           `json:"inputFileDir"`      //where the above image expects the dump to be(if from dump)
	OutputFileDir     string           `json:"outputFileDir"`     //where the above image will place its result config file
	ConfImage         string           `json:"confImage"`         //image that makes one init file from multiple init files
}

// KobeFederatorStatus defines the observed state of KobeFederator
// +k8s:openapi-gen=true
type KobeFederatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	PodNames []string `json:"podnames"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KobeFederator is the Schema for the kobefederators API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
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

func init() {
	SchemeBuilder.Register(&KobeFederator{}, &KobeFederatorList{})
}
