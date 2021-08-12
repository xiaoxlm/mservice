package components

import (
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MDeployment struct {
	Replicas *int32 `json:"replicas,omitempty"`
	Strategy appv1.DeploymentStrategy `json:"strategy,omitempty"`
	corev1.PodSpec `json:",inline"`
}

func (md *MDeployment) Convert(meta *metav1.ObjectMeta, labels, annotations map[string]string) (client.Object, schema.GroupVersionKind) {
	deployment := new(appv1.Deployment)
	deployment.ObjectMeta = *meta
	deployment.SetLabels(labels)
	deployment.SetAnnotations(annotations)

	deployment.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: map[string]string{
			MATCH_LABEL_KEY: deployment.GetName(),
		},
	}

	//deployment.Spec.Template.SetLabels(labels)
	//deployment.Spec.Template.SetAnnotations(annotations)
	deployment.Spec.Template.Spec = md.PodSpec
	deployment.Spec.Strategy = md.Strategy

	return deployment, deployment.GroupVersionKind()
}
