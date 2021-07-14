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
	"errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	pipelinev1alpha1 "github.com/alam0rt/buildkite-operator/api/v1alpha1"
	"github.com/alam0rt/go-buildkite/v2/buildkite"
)

// AccessTokenReconciler reconciles a AccessToken object
type AccessTokenReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=pipeline.buildkite.alam0rt.io,resources=accesstokens,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pipeline.buildkite.alam0rt.io,resources=accesstokens/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pipeline.buildkite.alam0rt.io,resources=accesstokens/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;watch;list

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AccessToken object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *AccessTokenReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	_ = log.FromContext(ctx)
	var accessToken pipelinev1alpha1.AccessToken
	if err := r.Get(ctx, req.NamespacedName, &accessToken); err != nil {
		log.Log.Error(err, "unable to fetch AccessToken")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	tokenSecret := &corev1.Secret{}
	if err := r.Get(
		ctx,
		client.ObjectKey{
			Name:      accessToken.Spec.SecretRef,
			Namespace: req.Namespace,
		},
		tokenSecret); err != nil {
		log.Log.Error(err, "unable to get secret")
		return ctrl.Result{}, err
	}

	apiKey := string(tokenSecret.Data["token"])
	if len(apiKey) == 0 {
		err := errors.New("supplied token is empty")
		log.Log.Error(err, "unable to autoenticate to Buildkite")
		return ctrl.Result{}, err
	}

	config, err := buildkite.NewTokenConfig(apiKey, false)
	if err != nil {
		log.Log.Error(err, "unable to authenticate to buildkite using supplied token")
		// this error is bad and thus we exit
		return ctrl.Result{}, err
	}

	implemented := false
	if implemented {
		client := buildkite.NewClient(config.Client())

		var token *buildkite.AccessToken
		token, _, err = client.GetToken()
		if err != nil {
			log.Log.Error(err, "there was a problem getting the current access token")
			// this error is bad and thus we exit
			return ctrl.Result{}, err
		}
		log.Log.Info(*token.UUID)
	}

	accessToken.Status.Token = apiKey

	if err := r.Status().Update(ctx, &accessToken); err != nil {
		log.Log.Error(err, "unable to update AccessToken status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AccessTokenReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pipelinev1alpha1.AccessToken{}).
		Complete(r)
}
