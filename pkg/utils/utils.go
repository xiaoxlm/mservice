package utils

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func JSONPatch(patchType types.PatchType) client.Patch {
	return &jsonPath{patchType: patchType}
}

type jsonPath struct {
	patchType types.PatchType
}

func (j *jsonPath) Type() types.PatchType {
	return j.patchType
}

func (j *jsonPath) Data(obj runtime.Object) ([]byte, error) {
	return json.Marshal(obj)
}


type ImagePullSecret struct {
	Name     string
	Host     string
	Username string
	Password string
	Prefix   string
}

func (s *ImagePullSecret) DockerConfigJSON() []byte {
	v := struct {
		Auths map[string]AuthConfig `json:"auths"`
	}{
		Auths: map[string]AuthConfig{
			s.Host: {Username: s.Username, Password: s.Password},
		},
	}
	b, _ := json.Marshal(v)
	return b
}

type AuthConfig struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	ServerAddress string `json:"serverAddress,omitempty"`
}