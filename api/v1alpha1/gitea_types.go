/*
Copyright 2023 Soumyadip Chowdhury.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GiteaSpec defines the desired state of Gitea
type GiteaSpec struct {

	Name int `json:"name"`
	Image int32 `json:"image"`
	Versions int32 `json:"versions"`
	Port int32 `json:"port"`
	NodePort int32 `json:"nodeport"`
	TargetPort int32 `json:"targetport"`
	Replicas int32 `json:"replicas"`
	ContainerPort int32 `json:"containerport"`
	
}

// GiteaStatus defines the observed state of Gitea
type GiteaStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Gitea is the Schema for the gitea API
type Gitea struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GiteaSpec   `json:"spec,omitempty"`
	Status GiteaStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GiteaList contains a list of Gitea
type GiteaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Gitea `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Gitea{}, &GiteaList{})
}
