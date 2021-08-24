package components

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// service
type MPort struct {
	AppProtocol string          `json:"appProtocol,omitempty"`
	Port        int32           `json:"port,omitempty"`
	IsNodePort  bool            `json:"isNodePort,omitempty"`
	TargetPort  uint16          `json:"targetPort,omitempty"`
	Protocol    corev1.Protocol `json:"protocol,omitempty"`
}

type MPorts []MPort

func (mp *MPorts) Convert(meta *metav1.ObjectMeta, labels, annotations map[string]string) (client.Object, schema.GroupVersionKind) {
	s := new(corev1.Service)
	s.SetName(meta.GetName())
	s.SetNamespace(meta.GetNamespace())
	//s.ObjectMeta = *meta

	s.Spec = *mp.toServiceSpec()
	s.Spec.Selector = map[string]string{
		"app": s.ObjectMeta.Name,
	}
	return s, schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Service",
	}
}

func (mp *MPorts) toServiceSpec() *corev1.ServiceSpec {
	serviceSpec := new(corev1.ServiceSpec)
	serviceSpec.Type = corev1.ServiceTypeClusterIP

	for _, port := range *mp {
		servicePort := corev1.ServicePort{}

		appProtocol := port.AppProtocol
		if appProtocol == "" {
			appProtocol = "http"
		}

		servicePort.Name = fmt.Sprintf("%s-%d", appProtocol, port.Port)
		servicePort.Port = port.Port
		servicePort.TargetPort = intstr.FromInt(int(port.TargetPort))

		if port.IsNodePort {
			serviceSpec.Type = corev1.ServiceTypeNodePort
			servicePort.Name = "np-" + servicePort.Name
			servicePort.NodePort = int32(port.Port)
		}

		servicePort.Protocol = port.Protocol
		serviceSpec.Ports = append(serviceSpec.Ports, servicePort)
	}

	return serviceSpec
}
