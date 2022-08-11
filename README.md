# Image Clone Controller

Image Clone Controller is a kubernetes controller to safeguard against the risk of public container images disappearing from the registry while
we use them, breaking our deployments or daemonsets.

This controller make sure all the deployments and daemonsets running on the cluster are running using the images from our own registry which we
have configured with it. However, it will ignore kube-system namespace and the namespace where the controller is running.

### Prerequisites

* a kubernetes cluster
* an image registry account where we would back up the images e.g. docker.io or quay.io

### How to run?

Before installing the controller, we need to create a secret with our image registry credentials

Create `image-clone-controller` namespace
```shell
  kubectl create namespace image-clone-controller
```

Now, you can use the following command to create the secret
```shell
  kubectl create secret generic -n image-clone-controller backup-registry-credentials \
    --from-literal=REGISTRY_NAME=<> \
    --from-literal=REPOSITORY_NAME=<> \
    --from-literal=REGISTRY_USERNAME=<> \
    --from-literal=REGISTRY_PASSWORD=<> 
```
NOTE: If `REPOSITORY_NAME` is passed empty, controller will use `REGISTRY_USERNAME` as `REPOSITORY_NAME`.

Now, you can use the dev-release.yaml to create rest of the resources. You can build an images using Dockerfile or use the
pre-built image.

```shell
  kubectl apply -f https://raw.githubusercontent.com/sm43/image-clone-controller/main/dev-release.yaml
```

The dev-release file uses the `ghcr.io/sm43/image-clone-controller:main` built by GitHub Workflows on every new commit pushed, so
it will always be updated to the latest code.

#### Running from the code

You can run the controller by using [`ko`](https://github.com/google/ko) too.
```shell
  ko apply -f config/
```

### Improvements Ideas

* Controller clones images one after another, this could be improved by doing them parallely using go routines. So if we have more than
one image to be cloned in the deployment or daemonset, we can start one go routine for each image and clone them.

* This works with docker.io as image pushed to it will be public by default but in case of quay.io, the image would be private so, need to find a way
to make the image public or adding a secret for the deployment/daemonset and attaching as imagepullsecret.

* Image cloning can be decoupled from the controller and could be its own service which controller talks and request for cloning images. It would be async 
execution, so controller reconcile register a request and ask the clone service to copy the image and requeue the reconcile request. The clone service can also
keeps track of cloned images and return the name of image if already cloned. 
