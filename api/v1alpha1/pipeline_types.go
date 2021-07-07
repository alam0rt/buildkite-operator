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

// PipelineSpec defines the desired state of Pipeline
type PipelineSpec struct {
	// Steps for the pipeline to take
	// https://buildkite.com/docs/pipelines/defining-steps#step-defaults
	Steps []Step `json:"step"`
}

// Step represents a build step in buildkites build pipeline
type Step struct {
	Type                *string           `json:"type,omitempty"`
	Name                *string           `json:"name,omitempty"`
	Command             *string           `json:"command,omitempty"`
	ArtifactPaths       *string           `json:"artifact_paths,omitempty"`
	BranchConfiguration *string           `json:"branch_configuration,omitempty"`
	Env                 map[string]string `json:"env,omitempty"`
	TimeoutInMinutes    *int              `json:"timeout_in_minutes,omitempty"`
	AgentQueryRules     []string          `json:"agent_query_rules,omitempty"`
}

// PipelineStatus defines the observed state of Pipeline
type PipelineStatus struct {
	// Last observed state of the build job
	BuildState BuildState `json:"build"`
}

// BuildState describes the state of the overall pipeline.
// https://buildkite.com/docs/pipelines/defining-steps#build-states
// +kubebuilder:validation:Enum=Scheduled;Running;Passed;Failed;Blocked;Canceling;Canceled;Skipped;NotRun;Finished
type BuildState string

const (
	// ScheduledBuildState indicates a build has been scheduled by Buildkite
	ScheduledBuildState BuildState = "Scheduled"

	// RunningBuildState means the build is running on an agent
	RunningBuildState BuildState = "Running"

	// PassedBuildState states that a build has successfully completed running
	PassedBuildState BuildState = "Passed"

	// FailedBuildState indicates the build has exited on a non-zero value
	FailedBuildState BuildState = "Failed"

	// BlockedBuildState means the build is waiting for manual interventing to be unblocked
	BlockedBuildState BuildState = "Blocked"

	// CancelingBuildState indicates a build has been instructed to cancel and is doing so
	CancelingBuildState BuildState = "Canceling"

	// CanceledBuildState means the build has successfully canceled
	CanceledBuildState BuildState = "Canceled"

	// SkippedBuildState means the build has skipped due to some logic in the pipeline
	SkippedBuildState BuildState = "Skipped"

	// NotRunBuildState indicates the build has not run and will not be re-run
	NotRunBuildState BuildState = "NotRun"

	// FinishedBuildState means the build has passed and is an alias for `passed`,
	// `failed`, `blocked` and `cancelled`.
	FinishedBuildState BuildState = "Finished"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Pipeline is the Schema for the pipelines API
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineSpec   `json:"spec,omitempty"`
	Status PipelineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PipelineList contains a list of Pipeline
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pipeline `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pipeline{}, &PipelineList{})
}
