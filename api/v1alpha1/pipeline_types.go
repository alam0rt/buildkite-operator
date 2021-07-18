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

// PipelineSpec defines the desired state of Pipeline
type PipelineSpec struct {
	// Steps for the pipeline to take
	// https://buildkite.com/docs/pipelines/defining-steps#step-defaults
	PipelineName string `json:"pipelineName"`
	Organization string `json:"organization"`
	Repository   string `json:"repository"`

	// Required for authentication
	AccessTokenRef string `json:"accessTokenRef"`

	// Only support YAML pipelines
	// https://buildkite.com/docs/apis/rest-api/pipelines#create-a-yaml-pipeline
	Configuration string `json:"configuration,omitempty"`

	// Optional fields
	DefaultBranch                   string            `json:"default_branch,omitempty"`
	Description                     string            `json:"description,omitempty"`
	Env                             map[string]string `json:"env,omitempty"`
	ProviderSettings                ProviderSettings  `json:"provider_settings,omitempty"`
	BranchConfiguration             string            `json:"branch_configuration,omitempty"`
	SkipQueuedBranchBuilds          bool              `json:"skip_queued_branch_builds,omitempty"`
	SkipQueuedBranchBuildsFilter    string            `json:"skip_queued_branch_builds_filter,omitempty"`
	CancelRunningBranchBuilds       bool              `json:"cancel_running_branch_builds,omitempty"`
	CancelRunningBranchBuildsFilter string            `json:"cancel_running_branch_builds_filter,omitempty"`
	TeamUuids                       []string          `json:"team_uuids,omitempty"`
	ClusterID                       string            `json:"cluster_id,omitempty"`
}

// PipelineStatus defines the observed state of Pipeline
type PipelineStatus struct {
	// Last observed state of the build job
	Slug       *string      `json:"slug,omitempty"`
	URL        *string      `json:"url,omitempty"`
	WebURL     *string      `json:"web_url,omitempty"`
	BuildsURL  *string      `json:"builds_url,omitempty"`
	BadgeURL   *string      `json:"badge_url,omitempty"`
	CreatedAt  *metav1.Time `json:"created_at,omitempty"`
	ArchivedAt *metav1.Time `json:"archived_at,omitempty"`

	ScheduledBuildsCount *int `json:"scheduled_builds_count,omitempty"`
	RunningBuildsCount   *int `json:"running_builds_count,omitempty"`
	ScheduledJobsCount   *int `json:"scheduled_jobs_count,omitempty"`
	RunningJobsCount     *int `json:"running_jobs_count,omitempty"`
	WaitingJobsCount     *int `json:"waiting_jobs_count,omitempty"`

	// the provider of sources
	Provider *Provider `json:"provider,omitempty" yaml:"provider,omitempty"`

	// build state
	//	BuildState BuildState `json:"build"`
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
//+kubebuilder:printcolumn:name="Slug",type=string,JSONPath=`.status.slug`
//+kubebuilder:printcolumn:name="Organization",type=string,JSONPath=`.spec.organization`
//+kubebuilder:printcolumn:name="Url",type=string,JSONPath=`.status.builds_url`

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

// Provider represents a source code provider. It is read-only, but settings may be written using Pipeline.ProviderSettings.
type Provider struct {
	ID         string  `json:"id"`
	WebhookURL *string `json:"webhook_url"`
}

// ProviderSettings are settings for the pipeline provider.
type ProviderSettings struct {
	GitHubSettings           GitHubSettings           `json:"github_settings,omitempty"`
	GitHubEnterpriseSettings GitHubEnterpriseSettings `json:"github_enterprise_settings,omitempty"`
	BitbucketSettings        BitbucketSettings        `json:"bitbucket_settings,omitempty"`
	GitLabSettings           GitLabSettings           `json:"gitlab_settings,omitempty"`
}

// BitbucketSettings are settings for pipelines building from Bitbucket repositories.
type BitbucketSettings struct {
	BuildPullRequests                       *bool   `json:"build_pull_requests,omitempty"`
	PullRequestBranchFilterEnabled          *bool   `json:"pull_request_branch_filter_enabled,omitempty"`
	PullRequestBranchFilterConfiguration    *string `json:"pull_request_branch_filter_configuration,omitempty"`
	SkipPullRequestBuildsForExistingCommits *bool   `json:"skip_pull_request_builds_for_existing_commits,omitempty"`
	BuildTags                               *bool   `json:"build_tags,omitempty"`
	PublishCommitStatus                     *bool   `json:"publish_commit_status,omitempty"`
	PublishCommitStatusPerStep              *bool   `json:"publish_commit_status_per_step,omitempty"`

	// Read-only
	Repository *string `json:"repository,omitempty"`
}

// GitHubSettings are settings for pipelines building from GitHub repositories.
type GitHubSettings struct {
	TriggerMode                             *string `json:"trigger_mode,omitempty"`
	BuildPullRequests                       *bool   `json:"build_pull_requests,omitempty"`
	PullRequestBranchFilterEnabled          *bool   `json:"pull_request_branch_filter_enabled,omitempty"`
	PullRequestBranchFilterConfiguration    *string `json:"pull_request_branch_filter_configuration,omitempty"`
	SkipPullRequestBuildsForExistingCommits *bool   `json:"skip_pull_request_builds_for_existing_commits,omitempty"`
	BuildPullRequestForks                   *bool   `json:"build_pull_request_forks,omitempty"`
	PrefixPullRequestForkBranchNames        *bool   `json:"prefix_pull_request_fork_branch_names,omitempty"`
	BuildTags                               *bool   `json:"build_tags,omitempty"`
	PublishCommitStatus                     *bool   `json:"publish_commit_status,omitempty"`
	PublishCommitStatusPerStep              *bool   `json:"publish_commit_status_per_step,omitempty"`
	FilterEnabled                           *bool   `json:"filter_enabled,omitempty"`
	FilterCondition                         *string `json:"filter_condition,omitempty"`
	SeparatePullRequestStatuses             *bool   `json:"separate_pull_request_statuses,omitempty"`
	PublishBlockedAsPending                 *bool   `json:"publish_blocked_as_pending,omitempty"`

	// Read-only
	Repository *string `json:"repository,omitempty"`
}

// GitHubEnterpriseSettings are settings for pipelines building from GitHub Enterprise repositories.
type GitHubEnterpriseSettings struct {
	BuildPullRequests                       *bool   `json:"build_pull_requests,omitempty"`
	PullRequestBranchFilterEnabled          *bool   `json:"pull_request_branch_filter_enabled,omitempty"`
	PullRequestBranchFilterConfiguration    *string `json:"pull_request_branch_filter_configuration,omitempty"`
	SkipPullRequestBuildsForExistingCommits *bool   `json:"skip_pull_request_builds_for_existing_commits,omitempty"`
	BuildTags                               *bool   `json:"build_tags,omitempty"`
	PublishCommitStatus                     *bool   `json:"publish_commit_status,omitempty"`
	PublishCommitStatusPerStep              *bool   `json:"publish_commit_status_per_step,omitempty"`

	// Read-only
	Repository *string `json:"repository,omitempty"`
}

// GitLabSettings are settings for pipelines building from GitLab repositories.
type GitLabSettings struct {
	// Read-only
	Repository *string `json:"repository,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Pipeline{}, &PipelineList{})
}
