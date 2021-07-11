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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	pipelinev1alpha1 "github.com/alam0rt/buildkite-operator/api/v1alpha1"
	"github.com/alam0rt/go-buildkite/v2/buildkite"
)

// OrganizationReconciler reconciles a Organization object
type OrganizationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=pipeline.buildkite.alam0rt.io,resources=organizations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pipeline.buildkite.alam0rt.io,resources=organizations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pipeline.buildkite.alam0rt.io,resources=organizations/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Organization object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *OrganizationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	var organization pipelinev1alpha1.Organization
	if err := r.Get(ctx, req.NamespacedName, &organization); err != nil {
		log.Log.Error(err, "unable to fetch Organization")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var token string
	token = *organization.Spec.Token

	config, err := buildkite.NewTokenConfig(token, false)
	if err != nil {
		log.Log.Error(err, "unable to authenticate to buildkite using supplied token")
		return ctrl.Result{}, err
	}

	log.Log.Info("loaded Organization")

	client := buildkite.NewClient(config.Client())

	var slug string
	slug = *organization.Spec.Slug

	remoteOrganization, _, err := client.Organizations.Get(slug)
	if err != nil {
		log.Log.Error(err, "unable to fetch remote organization using Buildkite's API")
		return ctrl.Result{}, err
	}

	// update the status of the resource
	organization.Status.Name = remoteOrganization.Name
	organization.Status.ID = remoteOrganization.ID
	organization.Status.URL = remoteOrganization.URL
	organization.Status.AgentsURL = remoteOrganization.AgentsURL
	organization.Status.WebURL = remoteOrganization.WebURL

	var time string
	time = remoteOrganization.CreatedAt.Time.String()
	organization.Status.CreatedAt = &time

	if err := r.Status().Update(ctx, &organization); err != nil {
		log.Log.Error(err, "unable to update Organization status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OrganizationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pipelinev1alpha1.Organization{}).
		Complete(r)
}
