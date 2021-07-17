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

// +kubebuilder:validation:Enum=read_agents;write_agents;read_teams;read_artifacts;write_artifacts;read_builds;write_builds;read_job_env;read_build_logs;write_build_logs;read_organizations;read_pipelines;write_pipelines;read_user
type Scope string

const (
	// Permission to list and retrieve details of agents
	ReadAgentsScope Scope = "read_agents"

	// Permission to create, update and delete agents
	WriteAgentsScope Scope = "write_agents"

	// Permission to list teams
	ReadTeamsScope Scope = "read_teams"

	// Permission to retrieve build artifacts
	ReadArtifactsScope Scope = "read_artifacts"

	// Permission to delete build artifacts
	WriteArtifactsScope Scope = "write_artifacts"

	// Permission to list and retrieve details of builds
	ReadBuildsScope Scope = "read_builds"

	// Permission to create new builds
	WriteBuildsScope Scope = "write_builds"

	// Permission to retrieve job environment variables
	ReadJobEnvScope Scope = "read_job_env"

	// Permission to retrieve build logs
	ReadBuildLogsScope Scope = "read_build_logs"

	// Permission to delete build logs
	WriteBuildLogsScope Scope = "write_build_logs"

	// Permission to list and retrieve details of organizations
	ReadOrganizationsScope Scope = "read_organizations"

	// Permission to list and retrieve details of pipelines
	ReadPipelinesScope Scope = "read_pipelines"

	// Permission to create, update and delete pipelines
	WritePipelinesScope Scope = "write_pipelines"

	// Permission to retrieve basic details of the user
	ReadUserScope Scope = "read_user"
)

type Scopes []Scope

// AccessTokenSpec defines the desired state of AccessToken
type AccessTokenSpec struct {
	// SecretRef takes  a secret name and uses the `token` key to authenticate
	// all requests.
	SecretRef string `json:"secretRef"`
}

// AccessTokenStatus defines the observed state of AccessToken
type AccessTokenStatus struct {
	// UUID is the unique identifier of the supplied access token
	UUID string `json:"uuid,omitempty"`

	// Scopes are a list of the access permissions of the token
	// https://buildkite.com/docs/apis/managing-api-tokens#token-scopes
	Scopes []Scope `json:"scopes,omitempty"`

	// Token is the Buildkite API token retrieved from the AccessToken resource
	Token string `json:"token,omitempty"`

	// LastTimeAuthenticated is the time when the token was used last to successfully
	// authenticate to the Buildkite API
	// TOOD: use *metav1.Time instead of string
	LastTimeAuthenticated string `json:"lastTimeAuthenticated,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="UUID",type=string,JSONPath=`.status.uuid`
//+kubebuilder:printcolumn:name="Scopes",type=string,JSONPath=`.status.scopes`

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
