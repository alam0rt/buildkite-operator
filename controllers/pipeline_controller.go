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

package controllers

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/alam0rt/buildkite-operator/api/v1alpha1"
	pipelinev1alpha1 "github.com/alam0rt/buildkite-operator/api/v1alpha1"
	"github.com/alam0rt/go-buildkite/v2/buildkite"
)

// PipelineReconciler reconciles a Pipeline object
type PipelineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func nameToSlug(s string) string {
	s = strings.ToLower(s)
	split := strings.Split(s, " ")
	slug := strings.Join(split, "-")
	return slug
}

//+kubebuilder:rbac:groups=pipeline.buildkite.alam0rt.io,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pipeline.buildkite.alam0rt.io,resources=pipelines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pipeline.buildkite.alam0rt.io,resources=pipelines/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Pipeline object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *PipelineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	var pipeline pipelinev1alpha1.Pipeline
	if err := r.Get(ctx, req.NamespacedName, &pipeline); err != nil {
		log.Log.Error(err, "unable to fetch Pipeline")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var accessTokenResource pipelinev1alpha1.AccessToken
	if err := r.Get(ctx, client.ObjectKey{Name: pipeline.Spec.AccessTokenRef, Namespace: req.Namespace}, &accessTokenResource); err != nil {
		log.Log.Error(err, "unable to fetch referenced AccessToken")
		return ctrl.Result{}, err
	}

	apiToken := accessTokenResource.Status.Token
	config, err := buildkite.NewTokenConfig(apiToken, false)
	if err != nil {
		log.Log.Error(err, "unable to authenticate to buildkite using supplied token")
		// this error is bad and thus we exit
		return ctrl.Result{}, err
	}

	client := buildkite.NewClient(config.Client())

	organization := pipeline.Spec.Organization

	// check if exists
	var resp *buildkite.Pipeline
	if pipeline.Status.Slug != nil {
		guessSlug := nameToSlug(pipeline.ObjectMeta.Name)
		resp, _, err = client.Pipelines.Get(organization, guessSlug)
		if err != nil {
			log.Log.Error(err, "there was an problem when retrieving the pipeline from the Buildkite API")
		} else {

			// this is the most important thing to set
			pipeline.Status.Slug = resp.Slug

			// update the status
			pipeline.Status.CreatedAt = (*v1.Time)(resp.ArchivedAt)
			pipeline.Status.BuildState = pipelinev1alpha1.PassedBuildState
			pipeline.Status.RunningBuildsCount = resp.RunningBuildsCount
			pipeline.Status.ScheduledBuildsCount = resp.ScheduledBuildsCount
			pipeline.Status.ScheduledJobsCount = resp.ScheduledJobsCount
			pipeline.Status.RunningJobsCount = resp.RunningJobsCount
			pipeline.Status.BuildsURL = resp.BuildsURL
			pipeline.Status.Provider = &v1alpha1.Provider{
				ID:         resp.Provider.ID,
				WebhookURL: resp.Provider.WebhookURL,
			}
		}
	} else {
		// if we can't retrieve the pipeline we assume it may not exist yet
		input := &buildkite.CreatePipeline{
			Name:          pipeline.Spec.PipelineName,
			Repository:    pipeline.Spec.Repository,
			Configuration: pipeline.Spec.Configuration,
		}

		jsonInput, _ := json.Marshal(input)
		log.Log.Info(string(jsonInput))

		resp, apiResp, err := client.Pipelines.Create(organization, input)
		if err != nil {
			if apiResp.Response.StatusCode == 422 {
				log.Log.Error(err, "there was an issue validating the pipeline input")
				return ctrl.Result{}, err
			}
			log.Log.Error(err, "there was an unknown exception")
		}

		// this is the most important thing to set
		pipeline.Status.Slug = resp.Slug

		// update the status
		pipeline.Status.CreatedAt = (*v1.Time)(resp.ArchivedAt)
		pipeline.Status.BuildState = pipelinev1alpha1.PassedBuildState
		pipeline.Status.RunningBuildsCount = resp.RunningBuildsCount
		pipeline.Status.ScheduledBuildsCount = resp.ScheduledBuildsCount
		pipeline.Status.ScheduledJobsCount = resp.ScheduledJobsCount
		pipeline.Status.RunningJobsCount = resp.RunningJobsCount
		pipeline.Status.BuildsURL = resp.BuildsURL
		pipeline.Status.Provider = &v1alpha1.Provider{
			ID:         resp.Provider.ID,
			WebhookURL: resp.Provider.WebhookURL,
		}

	}

	if err := r.Status().Update(ctx, &pipeline); err != nil {
		log.Log.Error(err, "unable to update Pipeline status")
		return ctrl.Result{}, err
	}

	time.Sleep(time.Duration(5) * time.Second)

	return ctrl.Result{}, nil
}

type PipelineErrorResponse struct {
	Message string          `json:"message"`
	Errors  []PipelineError `json:"errors"`
}

type PipelineError struct {
	Field string `json:"field,omitempty"`
	Code  string `json:"code,omitemtpy"`
	Value string `json:"value,omitempty"`
}

// SetupWithManager sets up the controller with the Manager.
func (r *PipelineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pipelinev1alpha1.Pipeline{}).
		Complete(r)
}
