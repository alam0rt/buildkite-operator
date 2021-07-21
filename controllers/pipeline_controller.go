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
	"errors"
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

//+kubebuilder:rbac:groups=pipeline.buildkite.alam0rt.io,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pipeline.buildkite.alam0rt.io,resources=pipelines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pipeline.buildkite.alam0rt.io,resources=pipelines/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
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

	// check that supplied access token has the correct permissions
	for _, scope := range accessTokenResource.Status.Scopes {
		if scope == pipelinev1alpha1.WritePipelinesScope {
			log.Log.Info("provided token has write pipeline permissions")
		}
	}

	apiToken := accessTokenResource.Status.Token
	config, err := buildkite.NewTokenConfig(apiToken, false)
	if err != nil {
		log.Log.Error(err, "unable to authenticate to buildkite using supplied token")
		// this error is bad and thus we exit
		return ctrl.Result{}, err
	}

	// check that the remote pipeline exists
	var resp *buildkite.Pipeline
	var httpResp *buildkite.Response

	nameSlug := pipeline.Spec.PipelineName // we make the opinion that slug needs to equal the pipeline name - alas, nameSlug :)
	organization := pipeline.Spec.Organization

	client := buildkite.NewClient(config.Client())
	resp, httpResp, err = client.Pipelines.Get(organization, nameSlug)
	if httpResp.Response.StatusCode != 200 {
		log.Log.Info("could not find pipeline remotely, a pipeline will be created")
		// if we can't retrieve the pipeline we assume it may not exist yet
		input := &buildkite.CreatePipeline{
			Name:          nameSlug,
			Repository:    pipeline.Spec.Repository,
			Configuration: pipeline.Spec.Configuration,
		}

		jsonInput, _ := json.Marshal(input)
		log.Log.Info(string(jsonInput))
		var apiResp *buildkite.Response

		resp, apiResp, err = client.Pipelines.Create(organization, input)
		if err != nil {
			if apiResp.Response.StatusCode == 422 {
				log.Log.Error(err, "there was an issue validating the pipeline input")
				return ctrl.Result{}, err
			}
			log.Log.Error(err, "there was an unknown exception")
		}
		log.Log.Info("successfully created pipeline")

	} else if httpResp.StatusCode == 200 {
		pipeline.Status.Slug = resp.Slug
		if nameSlug != *resp.Name {
			log.Log.Info("remote pipeline was found but does not match expected name - will update to make it match")
			resp.Name = &nameSlug
		}

		if nameSlug != *resp.Slug {
			// this check is just to confirm the slug used to retrieve the pipeline matches said pipelines slug
			// which it always should...
			err = errors.New("provided slug does not match remote slug and should never occur")
			log.Log.Error(err, "something is very wrong")
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Log.Error(err, "there was a problem creating the pipeline")
		return ctrl.Result{}, err

	}

	resp.Name = &nameSlug // the name must equal the slug and vice versa
	resp.BranchConfiguration = &pipeline.Spec.BranchConfiguration
	resp.CancelRunningBranchBuilds = &pipeline.Spec.CancelRunningBranchBuilds
	resp.CancelRunningBranchBuildsFilter = &pipeline.Spec.CancelRunningBranchBuildsFilter
	resp.DefaultBranch = &pipeline.Spec.DefaultBranch
	resp.Description = &pipeline.Spec.Description
	resp.Repository = &pipeline.Spec.Repository
	resp.SkipQueuedBranchBuilds = &pipeline.Spec.SkipQueuedBranchBuilds
	resp.SkipQueuedBranchBuildsFilter = &pipeline.Spec.SkipQueuedBranchBuildsFilter
	// TODO: implement Visibility when its in go-buildkite
	// resp.Visibility = &pipeline.Spec.Visibility

	// TODO: implement ProviderSettings

	updateResp, err := client.Pipelines.Update(organization, resp)
	if err != nil {
		log.Log.Error(err, "unable to update Pipeline API")
		return ctrl.Result{}, err
	} else if updateResp.Response.StatusCode == 200 {
		log.Log.Info("successfully reconciled the pipeline")
	}

	pipeline.Status = v1alpha1.PipelineStatus{
		Slug: resp.Slug,
		Provider: &pipelinev1alpha1.Provider{
			ID:         resp.Provider.ID,
			WebhookURL: resp.Provider.WebhookURL,
		},
		URL:       resp.URL,
		BuildsURL: resp.BuildsURL,
		BadgeURL:  resp.BadgeURL,

		CreatedAt:            (*v1.Time)(resp.CreatedAt),
		ArchivedAt:           (*v1.Time)(resp.ArchivedAt),
		RunningBuildsCount:   resp.RunningBuildsCount,
		ScheduledBuildsCount: resp.ScheduledBuildsCount,

		RunningJobsCount:   resp.RunningJobsCount,
		ScheduledJobsCount: resp.ScheduledJobsCount,
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
