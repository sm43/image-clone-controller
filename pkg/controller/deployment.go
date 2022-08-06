package controller

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Deployment struct {
	Client client.Client
}

var _ reconcile.Reconciler = &Deployment{}

func (r *Deployment) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info(fmt.Sprint("Reconciling deployment ", request.NamespacedName))

	return reconcile.Result{}, nil
}
