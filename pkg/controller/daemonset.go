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
	Cloner *imagecloner.ImageCloner
}

var _ reconcile.Reconciler = &DaemonSet{}

func (r *DaemonSet) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	logger := log.FromContext(ctx)

	// TODO: remove this
	// reconcile only from image-clone-controller namespace
	if !strings.HasPrefix(request.NamespacedName.String(), "image-clone-controller/test") {
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

	daemonSetUpdated := false

	// containers
	for i, container := range daemonSet.Spec.Template.Spec.Containers {
		// check if image is already backup image
		if r.Cloner.IsBackupImage(container.Image) {
			continue
		}
		// if image is not backup then clone it and replace
		clonedImage, err := r.Cloner.Clone(container.Image)
		if err != nil {
			return reconcile.Result{}, nil
		}
		daemonSet.Spec.Template.Spec.Containers[i].Image = clonedImage

		daemonSetUpdated = true
	}

	// init containers
	for i, container := range daemonSet.Spec.Template.Spec.InitContainers {
		// check if image is already backup image
		if r.Cloner.IsBackupImage(container.Image) {
			continue
		}
		// if image is not backup then clone it and replace
		clonedImage, err := r.Cloner.Clone(container.Image)
		if err != nil {
			return reconcile.Result{}, nil
		}
		daemonSet.Spec.Template.Spec.InitContainers[i].Image = clonedImage

		daemonSetUpdated = true
	}

	// ephemeral containers
	for i, container := range daemonSet.Spec.Template.Spec.EphemeralContainers {
		// check if image is already backup image
		if r.Cloner.IsBackupImage(container.Image) {
			continue
		}
		// if image is not backup then clone it and replace
		clonedImage, err := r.Cloner.Clone(container.Image)
		if err != nil {
			return reconcile.Result{}, nil
		}
		daemonSet.Spec.Template.Spec.EphemeralContainers[i].Image = clonedImage

		daemonSetUpdated = true
	}

	if daemonSetUpdated {
		err := r.Client.Update(ctx, daemonSet)
		if err != nil {
			logger.Error(err, "failed to update daemonSet")
			return reconcile.Result{Requeue: true}, nil
		}
	}
	return reconcile.Result{}, nil
}

func isDaemonSetReady(d *v1.DaemonSet) bool {
	if d.Status.DesiredNumberScheduled == d.Status.NumberReady && d.Status.DesiredNumberScheduled > 0 {
		return true
	}
	return false
}
