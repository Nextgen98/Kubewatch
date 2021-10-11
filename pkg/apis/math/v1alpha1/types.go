package v1alpha1

import meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// MyResource describes a MyResource resource

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Math struct {
	// TypeMeta is the metadata for the resource, like kind and apiversion
	meta_v1.TypeMeta `json:",inline"`
	// ObjectMeta contains the metadata for the particular object, including
	// things like...
	//  - name
	//  - namespace
	//  - self link
	//  - labels
	//  - ... etc ...
	meta_v1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the custom resource spec
	Spec MathSpec `json:"spec"`
}

// MyResourceSpec is the spec for a MyResource resource
type MathSpec struct {
	// Message and SomeValue are example custom spec fields
	//
	// this is where you would put your custom resource data
	FirstNum  *int32  `json:"firstNum"`
	SecondNum *int32  `json:"secondNum"`
	Operation string `json:"operation"`
}

// MyResourceList is a list of MyResource resources
type MathList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`

	Items []Math `json:"items"`
}
