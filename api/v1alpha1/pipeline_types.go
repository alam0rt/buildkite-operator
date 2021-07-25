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
	Organization string `json:"organization"`
	Repository   string `json:"repository"`

	// Required for authentication
	AccessTokenRef string `json:"accessTokenRef"`

	// Only support YAML pipelines
	// https://buildkite.com/docs/apis/rest-api/pipelines#create-a-yaml-pipeline
	Configuration string `json:"configuration,omitempty"`

	// Optional fields
	DefaultBranch       string            `json:"defaultBranch,omitempty"`
	Description         string            `json:"description,omitempty"`
	Env                 map[string]string `json:"env,omitempty"`
	ProviderSettings    ProviderSettings  `json:"providerSettings,omitempty"`
	BranchConfiguration string            `json:"branchConfiguration,omitempty"`

	// +kubebuilder:validation:Default=false
	SkipQueuedBranchBuilds       bool   `json:"skipQueuedBranchBuilds,omitempty"`
	SkipQueuedBranchBuildsFilter string `json:"skipQueuedBranchBuildsFilter,omitempty"`

	// +kubebuilder:validation:Default=false
	CancelRunningBranchBuilds bool `json:"cancelRunningBranchBuilds,omitempty"`

	CancelRunningBranchBuildsFilter string   `json:"cancelRunningBranchBuildsFilter,omitempty"`
	TeamUuids                       []string `json:"teamUUIDS,omitempty"`
	ClusterID                       string   `json:"clusterID,omitempty"`
}

// PipelineStatus defines the observed state of Pipeline
type PipelineStatus struct {
	// Last observed state of the build job
	Slug       *string      `json:"slug,omitempty"`
	URL        *string      `json:"URL,omitempty"`
	WebURL     *string      `json:"webURL,omitempty"`
	BuildsURL  *string      `json:"buildsURL,omitempty"`
	BadgeURL   *string      `json:"badgeURL,omitempty"`
	CreatedAt  *metav1.Time `json:"createdAt,omitempty"`
	ArchivedAt *metav1.Time `json:"archivedAt,omitempty"`

	ScheduledBuildsCount *int `json:"scheduledBuildsCount,omitempty"`
	RunningBuildsCount   *int `json:"runningBuildsCount,omitempty"`
	ScheduledJobsCount   *int `json:"scheduledJobsCount,omitempty"`
	RunningJobsCount     *int `json:"runningJobsCount,omitempty"`
	WaitingJobsCount     *int `json:"waitingJobsCount,omitempty"`

	// the provider of sources
	Provider *Provider `json:"provider,omitempty"`
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
//+kubebuilder:printcolumn:name="Organization",type=string,JSONPath=`.spec.organization`
//+kubebuilder:printcolumn:name="RunningBuilds",type=integer,JSONPath=`.status.runningBuildsCount`
//+kubebuilder:printcolumn:name="RunningJobs",type=integer,JSONPath=`.status.runningJobsCount`

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

// +kubebuilder:validation:Enum=github;github_enterprise;bitbucket;gitlab
type ProviderID string

const (
	GitHubProvider           ProviderID = "github"
	GitHubEnterpriseProvider ProviderID = "github_enterprise"
	BitbucketProvider        ProviderID = "bitbucket"
	GitLabProvider           ProviderID = "gitlab"
)

type Provider struct {
	ID         *string `json:"ID"`
	WebhookURL *string `json:"webhookURL"`
}

// ProviderSettings represents a source code provider. It is read-only, but settings may be written using Pipeline.ProviderSettings.
type ProviderSettings struct {
	GitHubSettings           GitHubSettings           `json:"githubSettings,omitempty"`
	GitHubEnterpriseSettings GitHubEnterpriseSettings `json:"githubEnterpriseSettings,omitempty"`
	BitbucketSettings        BitbucketSettings        `json:"bitbucketSettings,omitempty"`
	GitLabSettings           GitLabSettings           `json:"gitlabSettings,omitempty"`
}

// BitbucketSettings are settings for pipelines building from Bitbucket repositories.
type BitbucketSettings struct {
	// +kubebuilder:validation:Default=true
	BuildPullRequests *bool `json:"buildPullRequests,omitempty"`

	// +kubebuilder:validation:Default=false
	PullRequestBranchFilterEnabled       *bool   `json:"pullRequestBranchFilterEnabled,omitempty"`
	PullRequestBranchFilterConfiguration *string `json:"pullRequestBranchFilterConfiguration,omitempty"`
	// +kubebuilder:validation:Default=true
	SkipPullRequestBuildsForExistingCommits *bool `json:"skipPullRequestBuildsForExistingCommits,omitempty"`
	// +kubebuilder:validation:Default=false
	BuildTags *bool `json:"buildTags,omitempty"`

	// +kubebuilder:validation:Default=true
	PublishCommitStatus *bool `json:"publishCommitStatus,omitempty"`

	// +kubebuilder:validation:Default=false
	PublishCommitStatusPerStep *bool `json:"publishCommitStatusPerStep,omitempty"`

	// Read-only
	Repository *string `json:"repository,omitempty"`
}

// GitHubSettings are settings for pipelines building from GitHub repositories.
type GitHubSettings struct {
	// +kubebuilder:validation:Default="code"
	TriggerMode *string `json:"triggerMode,omitempty"`

	// +kubebuilder:validation:Default=true
	BuildPullRequests *bool `json:"buildPullRequests,omitempty"`

	// +kubebuilder:validation:Default=false
	PullRequestBranchFilterEnabled       *bool   `json:"pullRequestBranchFilterEnabled,omitempty"`
	PullRequestBranchFilterConfiguration *string `json:"pullRequestBranchFilterConfiguration,omitempty"`
	// +kubebuilder:validation:Default=true
	SkipPullRequestBuildsForExistingCommits *bool `json:"skipPullRequestBuildsForExistingCommits,omitempty"`

	// +kubebuilder:validation:Default=false
	BuildPullRequestForks *bool `json:"buildPullRequestForks,omitempty"`

	// +kubebuilder:validation:Default=true
	PrefixPullRequestForkBranchNames *bool `json:"prefixPullRequestForkBranchNames,omitempty"`

	// +kubebuilder:validation:Default=false
	BuildTags *bool `json:"buildTags,omitempty"`

	// +kubebuilder:validation:Default=true
	PublishCommitStatus *bool `json:"publishCommitStatus,omitempty"`

	// +kubebuilder:validation:Default=false
	PublishCommitStatusPerStep *bool `json:"publishCommitStatusPerStep,omitempty"`

	// +kubebuilder:validation:Default=true
	FilterEnabled   *bool   `json:"filterEnabled,omitempty"`
	FilterCondition *string `json:"filterCondition,omitempty"`

	// +kubebuilder:validation:Default=false
	SeparatePullRequestStatuses *bool `json:"separatePullRequestStatuses,omitempty"`

	// +kubebuilder:validation:Default=false
	PublishBlockedAsPending *bool `json:"publishBlockedAsPending,omitempty"`

	// Read-only
	Repository *string `json:"repository,omitempty"`
}

// GitHubEnterpriseSettings are settings for pipelines building from GitHub Enterprise repositories.
type GitHubEnterpriseSettings struct {
	// +kubebuilder:validation:Default=true
	BuildPullRequests *bool `json:"buildPullRequests,omitempty"`

	// +kubebuilder:validation:Default=false
	PullRequestBranchFilterEnabled *bool `json:"pullRequestBranchFilterEnabled,omitempty"`

	// +kubebuilder:validation:Default=true
	PullRequestBranchFilterConfiguration *string `json:"pullRequestBranchFilterConfiguration,omitempty"`

	// +kubebuilder:validation:Default=true
	SkipPullRequestBuildsForExistingCommits *bool `json:"skipPullRequestBuildsForExistingCommits,omitempty"`

	// +kubebuilder:validation:Default=false
	BuildTags *bool `json:"buildTags,omitempty"`

	// +kubebuilder:validation:Default=true
	PublishCommitStatus *bool `json:"publishCommitStatus,omitempty"`

	// +kubebuilder:validation:Default=false
	PublishCommitStatusPerStep *bool `json:"publishCommitStatusPerStep,omitempty"`

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
