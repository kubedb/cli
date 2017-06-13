# Build Instructions

## Requirements
- go1.8+
- glide

## Build Binary
```sh
# Install/Update dependency (needs glide)
$ glide slow

# Build
$ ./hack/make.py build
```

## Build Docker
```sh
# Build Docker image
$ ./hack/docker/operator/setup.sh
```

#### Push Docker Image
```sh
# This will push docker image to other repositories

# Add docker tag for your repository
$ docker tag kubedb/operator:<tag> <image>:<tag>

# Push Image
$ docker push <image>:<tag>

# Example:
$ docker tag kubedb/operator:default aerokite/operator:default
$ docker push aerokite/operator:default
```
