## Development Guide
This document is intended to be the canonical source of truth for things like supported toolchain versions for building KubeDB.
If you find a requirement that this doc does not capture, please submit an issue on github.

This document is intended to be relative to the branch in which it is found. It is guaranteed that requirements will change over time
for the development branch, but release branches of KubeDB should not change.

### Build KubeDB
Some of the KubeDB development helper scripts rely on a fairly up-to-date GNU tools environment, so most recent Linux distros should
work just fine out-of-the-box.

#### Setup GO
KubeDB is written in Google's GO programming language. Currently, KubeDB is developed and tested on **go 1.8.3**. If you haven't set up a GO
development environment, please follow [these instructions](https://golang.org/doc/code.html) to install GO.

#### Code Organization
KubeDB codebase is across various repositories under github.com/k8sdb organization. There are 5 categories of git repositories:

| Repository                            | Description                                                                                             |
|---------------------------------------|---------------------------------------------------------------------------------------------------------|
| https://github.com/k8sdb/apimachinery | Contains api types, clientset and KubeDB framework interfaces.                                          |
| https://github.com/k8sdb/db           | This repository contains operator for `db`, eg, https://github.com/k8sdb/postgres                       |
| https://github.com/k8sdb/db_exporter  | This repository contains Prometheus exporter for `db`, eg, https://github.com/k8sdb/postgres_exporter . |
| https://github.com/k8sdb/operator     | This repository contains the combined operator for all databased supported by KubeDB.                   |
| https://github.com/k8sdb/cli          | This repository contains CLI for KubeDB.                                                                |





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
