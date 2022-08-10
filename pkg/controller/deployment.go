package controller

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/sm43/image-clone-controller/pkg/imagecloner"
)

type Deployment struct {
	Client client.Client
	Cloner *imagecloner.ImageCloner
}

var _ reconcile.Reconciler = &Deployment{}

func (r *Deployment) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	logger := log.FromContext(ctx).WithValues("deployment", request.NamespacedName)

	// TODO: remove this
	// reconcile only from image-clone-controller namespace
	if !strings.HasPrefix(request.NamespacedName.String(), "image-clone-controller/test") {
		return reconcile.Result{}, nil
	}

	logger.Info(fmt.Sprint("Reconciling deployment ", request.NamespacedName))

	// fetch deployment from the cluster
	deployment := &v1.Deployment{}
	err := r.Client.Get(ctx, request.NamespacedName, deployment)
	if err != nil {
		logger.Error(err, fmt.Sprint("failed to get deployment: ", request.NamespacedName))
		return reconcile.Result{}, err
	}

	// check if deployment is ready
	if !isDeploymentReady(deployment) {
		logger.Info("deployment not in ready state, will reconcile once ready")
		return reconcile.Result{}, nil
	}

	deploymentUpdated, err := podImageCloner(r.Cloner, &deployment.Spec.Template.Spec)
	if deploymentUpdated {
		err := r.Client.Update(ctx, deployment)
		if err != nil {
			logger.Error(err, "failed to update deployment")
			return reconcile.Result{Requeue: true}, nil
		}
	}
	return reconcile.Result{}, nil
}

func isDeploymentReady(d *v1.Deployment) bool {
	if d.Status.Replicas == d.Status.ReadyReplicas && d.Status.Replicas > 0 {
		return true
	}
	return false
}
