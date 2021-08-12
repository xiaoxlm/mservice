package apply

import (
	"context"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"r.kubebuilder.io/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func setNamespace(o metav1.Object, namespace string) {
	o.SetNamespace(namespace)
}

func GetCurrentObject(co client.Object) (client.Object, error) {
	copyCO := co.DeepCopyObject()
	current, err := meta.Accessor(copyCO)
	if err != nil {
		return nil, err
	}
	return current.(client.Object), nil
}

func Action(ctx context.Context, c client.Client, namespace string, applyObject client.Object) error {
	applyObject.SetNamespace(namespace)
	currentObject, err := GetCurrentObject(applyObject)
	if err != nil {
		return err
	}

	if err := c.Get(ctx, client.ObjectKeyFromObject(applyObject), currentObject); err != nil {
		if apierrors.IsNotFound(err) {
			return c.Create(ctx, applyObject)
		}
		return err
	}

	return c.Patch(ctx, applyObject, utils.JSONPatch(types.MergePatchType))
}

//func applyIngress(ctx context.Context, c client.Client, namespace string, ingress *networkingv1.Ingress) error {
//	setNamespace(ingress, namespace)
//
//	current := new(networkingv1.Ingress)
//
//	if err := c.Get(ctx, client.ObjectKeyFromObject(ingress), current); err != nil {
//		if apierrors.IsNotFound(err) {
//			return c.Create(ctx, ingress)
//		}
//		return err
//	}
//
//	return c.Patch(ctx, ingress, utils.JSONPatch(types.MergePatchType))
//}


//func applySecret(ctx context.Context, c client.Client, name, namespace string) error {
//	imagePullSecret := new(utils.ImagePullSecret)
//
//	secret := new(corev1.Secret)
//	secret.Type = "kubernetes.io/dockerconfigjson"
//	secret.Name = name
//	secret.Data = map[string][]byte{
//		".dockerconfigjson": imagePullSecret.DockerConfigJSON(),
//	}
//	setNamespace(secret, namespace)
//
//	current := new(corev1.Secret)
//
//	err := c.Get(ctx, types.NamespacedName{Name: secret.Name, Namespace: namespace}, current)
//	if err != nil {
//		if apierrors.IsNotFound(err) {
//			return c.Create(ctx, secret)
//		}
//		return err
//	}
//
//	return c.Patch(ctx, secret, utils.JSONPatch(types.MergePatchType))
//}
//
//func applyService() {
//
//}
//
//func applyDeployment() {
//
//}