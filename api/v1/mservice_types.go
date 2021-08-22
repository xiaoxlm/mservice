/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"r.kubebuilder.io/pkg/components"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MServiceSpec defines the desired state of MService
type MServiceSpec struct {

	//components.MDeployment `json:",inline"`

	Ingress      *components.MIngress `json:"ingress,omitempty"`
	Ports        *components.MPorts    `json:"ports,omitempty"`
	Secret 		 *components.MSecret `json:"secret,omitempty"`
}

// MServiceStatus defines the observed state of MService
type MServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Stage string `json:"stage,omitempty"`
	//DeploymentStage         string `json:"deploymentStage,omitempty"`
	//DeploymentComments      string `json:"deploymentComments,omitempty"`
	//appsv1.DeploymentStatus `json:",inline"`
	//appsv1.DeploymentSpec
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MService is the Schema for the mservices API
type MService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MServiceSpec   `json:"spec,omitempty"`
	Status MServiceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MServiceList contains a list of MService
type MServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MService{}, &MServiceList{})
}
