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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Scope string

const (
	// TODO: implement access token API
	ReadAgentsScope Scope = "read_agents"
)

type Scopes []Scope

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AccessTokenSpec defines the desired state of AccessToken
type AccessTokenSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// SecretRef takes  a secret name and uses the `token` key to authenticate
	// all requests.
	SecretRef string `json:"secretRef"`
}

// AccessTokenStatus defines the observed state of AccessToken
type AccessTokenStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	UUID   string `json:"uuid,omitempty"`
	Scopes Scopes `json:"scopes,omitempty"`

	// Token is the Buildkite API token retrieved from the AccessToken resource
	Token string `json:"token,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AccessToken is the Schema for the accesstokens API
type AccessToken struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AccessTokenSpec   `json:"spec,omitempty"`
	Status AccessTokenStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AccessTokenList contains a list of AccessToken
type AccessTokenList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AccessToken `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AccessToken{}, &AccessTokenList{})
}
