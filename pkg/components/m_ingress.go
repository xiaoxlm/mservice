package components

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MIngress struct {
	Scheme string
	Host   string
	Path   string
	Port   uint16
}

func (*MIngress) Convert(meta *metav1.ObjectMeta) client.Object {
	return new(networkingv1.Ingress)
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