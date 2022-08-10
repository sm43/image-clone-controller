package controller

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/sm43/image-clone-controller/pkg/imagecloner"
)

type DaemonSet struct {
	Client client.Client
	Cloner *imagecloner.ImageCloner
}

var _ reconcile.Reconciler = &DaemonSet{}

func (r *DaemonSet) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info(fmt.Sprint("Reconciling daemonSet ", request.NamespacedName))

	return reconcile.Result{}, nil
}
