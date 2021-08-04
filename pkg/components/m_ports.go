package components

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MPorts []MPort

type MService struct {
	Ports []MPort `json:"ports,omitempty"`
}

type MPort struct {
	AppProtocol   string
	Port          uint16
	IsNodePort    bool
	ContainerPort uint16
	Protocol      string
}


func (mp *MPort) Convert(meta *metav1.ObjectMeta) client.Object {
	s := new(corev1.Service)
	s.ObjectMeta = *meta

	return s
}

func (*MPort) toServiceSpec() *corev1.ServiceSpec {
	serviceSpec := new(corev1.ServiceSpec)
	serviceSpec.Type = corev1.ServiceTypeClusterIP


	for _, port := range s.Spec.Ports {
		servicePort := v1.ServicePort{}

		appProtocol := port.AppProtocol

		if appProtocol == "" {
			appProtocol = "http"
		}

		servicePort.Name = fmt.Sprintf("%s-%d", appProtocol, port.Port)
		servicePort.Port = int32(port.Port)
		servicePort.TargetPort = intstr.FromInt(int(port.ContainerPort))

		if port.IsNodePort {
			serviceSpec.Type = v1.ServiceTypeNodePort
			servicePort.Name = "np-" + servicePort.Name
			servicePort.NodePort = int32(port.Port)
		}

		servicePort.Protocol = toProtocol(port.Protocol)
		serviceSpec.Ports = append(serviceSpec.Ports, servicePort)
	}

	return serviceSpec
}