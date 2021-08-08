package components

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ISelfResource interface {
	Convert(meta *metav1.ObjectMeta, labels, annotations map[string]string) client.Object
}
