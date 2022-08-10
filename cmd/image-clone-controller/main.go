package main

import (
	"log"

	v1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	klog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"

	ctrl "github.com/sm43/image-clone-controller/pkg/controller"
	"github.com/sm43/image-clone-controller/pkg/imagecloner"
)

func init() {
	klog.SetLogger(zap.New())
}

func main() {
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		log.Fatal("failed to create new manager", err)
	}

	cloner, err := imagecloner.NewCloner()
	if err != nil {
		log.Fatal("failed to init cloner: ", err)
	}

	// deployment controller
	deploymentCtrl, err := controller.New("deployment-controller", mgr, controller.Options{
		Reconciler: &ctrl.Deployment{Client: mgr.GetClient(), Cloner: cloner},
	})
	if err != nil {
		log.Fatal("failed to add controller: ", err)
	}

	if err := deploymentCtrl.Watch(&source.Kind{Type: &v1.Deployment{}}, &handler.EnqueueRequestForObject{}); err != nil {
		log.Fatal("failed to watch deployment: ", err)
	}

	// daemonSet controller
	daemonSetCtrl, err := controller.New("daemonset-controller", mgr, controller.Options{
		Reconciler: &ctrl.DaemonSet{Client: mgr.GetClient(), Cloner: cloner},
	})
	if err != nil {
		log.Fatal("failed to add controller: ", err)
	}

	if err := daemonSetCtrl.Watch(&source.Kind{Type: &v1.DaemonSet{}}, &handler.EnqueueRequestForObject{}); err != nil {
		log.Fatal("failed to watch daemonset: ", err)
	}

	// starting manager
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Fatal("failed to run manager: ", err)
	}
}
