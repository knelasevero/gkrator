/*
Copyright 2022.

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
	v1 "k8s.io/api/core/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GoogleKubernetesEngineSpec defines the desired state of GoogleKubernetesEngine
type GoogleKubernetesEngineSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ProjectID is the project where the cluster needs to be created.
	ProjectID string `json:"projectID,omitempty"`

	// Auth has fields that define one of the methods for authentication
	Auth Auth `json:"auth,omitempty"`
}

type Auth struct {
	SecretRef *GCPAuthSecretRef `json:"secretRef,omitempty"`

	// TODO
	// WorkloadID
}

type GCPAuthSecretRef struct {
	// The SecretAccessKey is used for authentication
	// +optional
	SecretAccessKey v1.SecretKeySelector `json:"secretAccessKeySecretRef,omitempty"`
}

// GoogleKubernetesEngineStatus defines the observed state of GoogleKubernetesEngine
type GoogleKubernetesEngineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GoogleKubernetesEngine is the Schema for the googlekubernetesengines API
type GoogleKubernetesEngine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GoogleKubernetesEngineSpec   `json:"spec,omitempty"`
	Status GoogleKubernetesEngineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GoogleKubernetesEngineList contains a list of GoogleKubernetesEngine
type GoogleKubernetesEngineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GoogleKubernetesEngine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GoogleKubernetesEngine{}, &GoogleKubernetesEngineList{})
}
