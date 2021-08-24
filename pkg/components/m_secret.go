package components

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MSecret struct {
	Name     string `json:"name,omitempty"`
	Host     string `json:"host,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Prefix   string `json:"prefix,omitempty"`
}

type AuthConfig struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	ServerAddress string `json:"serverAddress,omitempty"`
}

func (ms *MSecret) DockerConfigJSON() []byte {
	v := struct {
		Auths map[string]AuthConfig `json:"auths"`
	}{
		Auths: map[string]AuthConfig{
			ms.Host: {Username: ms.Username, Password: ms.Password},
		},
	}
	b, _ := json.Marshal(v)
	return b
}

func (ms *MSecret) Convert(meta *metav1.ObjectMeta, labels, annotations map[string]string) (client.Object, schema.GroupVersionKind) {
	secret := new(corev1.Secret)
	secret.SetName(meta.GetName())
	secret.SetNamespace(meta.GetNamespace())
	//secret.ObjectMeta = *meta
	secret.SetLabels(labels)
	secret.SetAnnotations(annotations)
	secret.Type = corev1.SecretTypeDockerConfigJson
	secret.Data = map[string][]byte{
		".dockerconfigjson": ms.DockerConfigJSON(),
	}

	return secret, schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Secret",
	}
}
