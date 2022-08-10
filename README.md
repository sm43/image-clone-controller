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