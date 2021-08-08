package components

import (
	//networking "k8s.io/api/networking/v1"
	networking "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MIngressSpec struct {
	Scheme string
	Host   string
	Path   string
	Port   uint16
}

type MIngress []*MIngressSpec

func (mi *MIngress) Convert(meta *metav1.ObjectMeta, labels, annotations map[string]string) client.Object {

	ingress := new(networking.Ingress)
	ingress.ObjectMeta = *meta
	ingress.SetLabels(labels)
	ingress.SetAnnotations(annotations)

	for _, i := range *mi {
		rule := networking.IngressRule{}
		rule.Host = i.Host

		path := networking.HTTPIngressPath{}
		path.Path = i.Path
		path.Backend.ServiceName = ingress.Name
		path.Backend.ServicePort = intstr.FromInt(int(i.Port))
		rule.IngressRuleValue.HTTP.Paths = []networking.HTTPIngressPath{path}

		ingress.Spec.Rules = append(ingress.Spec.Rules, rule)
	}
	return ingress
}

//func (mi MIngress) String() string {
//	if mi.Scheme == "" {
//		mi.Scheme = "http"
//	}
//	if mi.Port == 0 {
//		mi.Port = 80
//	}
//
//	return (&url.URL{
//		Scheme: mi.Scheme,
//		Host:   fmt.Sprintf("%s:%d", mi.Host, mi.Port),
//		Path:   mi.Path,
//	}).String()
//}
//
//func ParseUrlToIngress(urlString string) (*MIngress, error) {
//	if urlString == "" {
//		return nil, fmt.Errorf("url string is empty")
//	}
//
//	u, err := url.Parse(urlString)
//	if err != nil {
//		return nil, err
//	}
//
//	mi := &MIngress{
//		Scheme: u.Scheme,
//		Host:   u.Hostname(),
//		Path:   u.Path,
//	}
//
//	if mi.Scheme == "" {
//		mi.Scheme = "http"
//	}
//
//	port := u.Port()
//	if port == "" {
//		mi.Port = 80
//		return mi, nil
//	}
//
//	p, err := strconv.ParseUint(port, 10, 16)
//	if err != nil {
//		return nil, err
//	}
//
//	mi.Port = uint16(p)
//	return mi, nil
//}