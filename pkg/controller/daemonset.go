package controller

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/sm43/image-clone-controller/pkg/imagecloner"
)

type DaemonSet struct {
	Client client.Client
	Cloner imagecloner.ImageCloner
}

var _ reconcile.Reconciler = &DaemonSet{}

func (r *DaemonSet) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	logger := log.FromContext(ctx).WithValues("daemonSet", request.NamespacedName)

	// TODO: remove this
	// reconcile only from default namespace
	if !strings.HasPrefix(request.NamespacedName.String(), "default") {
		return reconcile.Result{}, nil
	}

	logger.Info(fmt.Sprint("Reconciling daemonSet ", request.NamespacedName))

	// fetch daemonSet from the cluster
	daemonSet := &v1.DaemonSet{}
	err := r.Client.Get(ctx, request.NamespacedName, daemonSet)
	if err != nil {
		logger.Error(err, fmt.Sprint("failed to get daemonSet: ", request.NamespacedName))
		return reconcile.Result{}, err
	}

	// check if daemonSet is ready
	if !isDaemonSetReady(daemonSet) {
		logger.Info("daemonSet not in ready state, will reconcile once ready")
		return reconcile.Result{}, nil
	}

	daemonSetUpdated, err := podImageCloner(&logger, r.Cloner, &daemonSet.Spec.Template.Spec)
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	if daemonSetUpdated {
		err := r.Client.Update(ctx, daemonSet)
		if err != nil {
			logger.Error(err, "failed to update daemonSet")
			return reconcile.Result{Requeue: true}, nil
		}
		logger.Info("updated images in daemonSet")
	}
	return reconcile.Result{}, nil
}

func isDaemonSetReady(d *v1.DaemonSet) bool {
	if d.Status.DesiredNumberScheduled == d.Status.NumberReady && d.Status.DesiredNumberScheduled > 0 {
		return true
	}
	return false
}
