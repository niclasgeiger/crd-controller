package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Foo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec FooSpec `json:"spec"`
}

type FooSpec struct {
	Bar string `json:"bar"`
}
