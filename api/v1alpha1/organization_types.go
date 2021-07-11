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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OrganizationSpec defines the desired state of Organization
type OrganizationSpec struct {
	Slug  *string `json:"slug,omitempty"`
	Token *string `json:"token"`
}

// OrganizationStatus defines the observed state of Organization
type OrganizationStatus struct {
	// Name is the human friendly representation of the organization
	Name *string `json:"name,omitempty"`

	Repository   *string `json:"repository,omitempty"`
	PipelinesURL *string `json:"pipelines_url,omitempty"`
	AgentsURL    *string `json:"agents_url,omitempty"`

	// ID is a UUID
	ID        *string `json:"id,omitempty"`
	URL       *string `json:"url,omitempty"`
	WebURL    *string `json:"web_url,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Organization is the Schema for the organizations API
type Organization struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OrganizationSpec   `json:"spec,omitempty"`
	Status OrganizationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OrganizationList contains a list of Organization
type OrganizationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Organization `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Organization{}, &OrganizationList{})
}
