package controller

import (
	"fmt"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"

	"github.com/sm43/image-clone-controller/pkg/imagecloner"
)

func podImageCloner(logger *logr.Logger, cloner imagecloner.ImageCloner, podSpec *v1.PodSpec) (bool, error) {
	resourceUpdated := false

	// containers
	for i, container := range podSpec.Containers {
		// check if image is already backup image
		if cloner.IsBackupImage(container.Image) {
			continue
		}
		logger.Info("cloning image : " + container.Image)
		// if image is not backup then clone it and replace
		clonedImage, err := cloner.Clone(container.Image)
		if err != nil {
			return resourceUpdated, err
		}
		podSpec.Containers[i].Image = clonedImage
		logger.Info(fmt.Sprintf("cloned image for %s => %s", container.Image, clonedImage))

		resourceUpdated = true
	}

	// init containers
	for i, container := range podSpec.InitContainers {
		// check if image is already backup image
		if cloner.IsBackupImage(container.Image) {
			continue
		}
		logger.Info("cloning image : " + container.Image)
		// if image is not backup then clone it and replace
		clonedImage, err := cloner.Clone(container.Image)
		if err != nil {
			return resourceUpdated, err
		}
		podSpec.InitContainers[i].Image = clonedImage
		logger.Info(fmt.Sprintf("cloned image for %s => %s", container.Image, clonedImage))

		resourceUpdated = true
	}

	// ephemeral containers
	for i, container := range podSpec.EphemeralContainers {
		// check if image is already backup image
		if cloner.IsBackupImage(container.Image) {
			continue
		}
		logger.Info("cloning image : " + container.Image)
		// if image is not backup then clone it and replace
		clonedImage, err := cloner.Clone(container.Image)
		if err != nil {
			return resourceUpdated, err
		}
		podSpec.EphemeralContainers[i].Image = clonedImage
		logger.Info(fmt.Sprintf("cloned image for %s => %s", container.Image, clonedImage))

		resourceUpdated = true
	}
	return resourceUpdated, nil
}
