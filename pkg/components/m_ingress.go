package components

import (
	//networking "k8s.io/api/networking/v1"
	networking "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MIngressSpec struct {
	Scheme string `json:"scheme,omitempty"`
	Host   string `json:"host,omitempty"`
	Path   string `json:"path,omitempty"`
	Port   uint16 `json:"port,omitempty"`
}

type MIngress []MIngressSpec

func (mi *MIngress) Convert(meta *metav1.ObjectMeta, labels, annotations map[string]string) (client.Object, schema.GroupVersionKind) {

	ingress := new(networking.Ingress)
	ingress.SetName(meta.GetName())
	ingress.SetNamespace(meta.GetNamespace())
	//ingress.ObjectMeta = *meta
	ingress.SetLabels(labels)
	ingress.SetAnnotations(annotations)

	for _, i := range *mi {
		path := networking.HTTPIngressPath{}
		path.Path = i.Path
		path.Backend.ServiceName = ingress.GetName()
		if i.Port < 1 {
			i.Port = 80
		}
		path.Backend.ServicePort = intstr.FromInt(int(i.Port))

		rule := networking.IngressRule{}
		rule.Host = i.Host
		rule.IngressRuleValue.HTTP = new(networking.HTTPIngressRuleValue)
		rule.IngressRuleValue.HTTP.Paths = []networking.HTTPIngressPath{path}

		ingress.Spec.Rules = append(ingress.Spec.Rules, rule)
	}
	return ingress, schema.GroupVersionKind{
		Group:   "extensions",
		Version: "v1beta1",
		Kind:    "Ingress",
	}
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
