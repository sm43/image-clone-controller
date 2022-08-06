package controller

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type DaemonSet struct {
	Client client.Client
}

var _ reconcile.Reconciler = &DaemonSet{}

func (r *DaemonSet) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info(fmt.Sprint("Reconciling daemonSet ", request.NamespacedName))

	return reconcile.Result{}, nil
}
