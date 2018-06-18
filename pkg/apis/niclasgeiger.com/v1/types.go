package v1

import (
	"fmt"
	"reflect"

	"github.com/niclasgeiger/crd-controller/pkg/apis/niclasgeiger.com"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	CRDSingular = "user"
	CRDPlural   = "users"
	CRDVersion  = "v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// User CRD
type User struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec UserSpec `json:"spec"`
}

type UserSpec struct {
	Port     int32 `json:"port"`
	NodePort int32 `json:"nodePort"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// List of User CRDs
type UserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []User `json:"items"`
}

func CreateCRD(clientset apiextcs.Interface) error {
	crd := &apiextv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Foo",
			APIVersion: CRDVersion,
		},
		ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("%s.%s", CRDPlural, niclasgeiger_com.GroupName)},
		Spec: apiextv1beta1.CustomResourceDefinitionSpec{
			Group:   niclasgeiger_com.GroupName,
			Version: CRDVersion,
			Scope:   apiextv1beta1.ClusterScoped,
			Names: apiextv1beta1.CustomResourceDefinitionNames{
				Kind:     reflect.TypeOf(User{}).Name(),
				Singular: CRDSingular,
				Plural:   CRDPlural,
				ShortNames: []string{
					"user",
				},
			},
		},
	}

	_, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil && apierrors.IsAlreadyExists(err) {
		return nil
	}
	return err
}
