package components

import (
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ISelfResource interface {
	Convert(meta *metav1.ObjectMeta, labels, annotations map[string]string) (client.Object, schema.GroupVersionKind)
}

func GetCurrentObject(co client.Object) (client.Object, error) {
	copyCO := co.DeepCopyObject()
	current, err := meta.Accessor(copyCO)
	if err != nil {
		return nil, err
	}
	return current.(client.Object), nil
}
