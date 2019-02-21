---
title: Release
menu:
  docs_0.10.0:
    identifier: developer-guide-release
    name: Release Process
    parent: developer-guide
    weight: 20
menu_name: docs_0.10.0
section_menu_id: setup
---

# Release Process

## apimachinery

- Make sure `apimachinery` master builds and create a `release-*` branch.

```console
apimachinery (master) $ ./hack/make.py
Downloading:  https://raw.githubusercontent.com/appscode/libbuild/master/libbuild.py
Using existing version:  github.appscode.libbuild.libbuild
Ungrouping imports of dir: api
Ungrouping imports of dir: client
Ungrouping imports of dir: pkg
goimports -w api client pkg
gofmt -s -w api client pkg
go build ./pkg/... ./api/... ./client/...
apimachinery (master) $ git checkout release-0.1
Switched to branch 'release-0.1'
apimachinery (release-0.1) $ git merge master
```

## postgres

- Change the apimachinery dependency to `release-*` branch
- Revendor `glide slow`
- Make sure master branch compiles, commit any changes due to vendoring & push to origin master
- Create a matching `release-*` branch
- Apply a matching `X.Y.Z` tag
- Push to origin `release-*` branch
- Build and release docker image for pg operator.
- Go back to master branch and change the dependency on `apimachinery` back to master branch.

```console
postgres (master) $ glide slow
postgres (master) $ ./hack/make.py
postgres (master) $ git commit -a -m 'Vendor apimachinery release-0.1'
postgres (master) $ git push origin master

postgres (master) $ git checkout release-0.1
postgres (release-0.1) $ git merge master
postgres (release-0.1) $ git tag -fa 0.1.0
postgres (release-0.1) $ git push origin release-0.1 --tags
postgres (release-0.1) $ rm -rf dist
postgres (release-0.1) $ ./hack/docker/pg-operator/setup.sh; env APPSCODE_ENV=prod ./hack/docker/pg-operator/setup.sh release

postgres (release-0.1) $ git checkout master
postgres (master) $ git commit -a -m 'Start 0.2 development cycle'
postgres (master) $ git push origin master
```

## elasticsearch

- Change the apimachinery dependency to `release-*` branch
- Revendor `glide slow`
- Make sure master branch compiles, commit any changes due to vendoring & push to origin master
- Create a matching `release-*` branch
- Apply a matching `X.Y.Z` tag
- Push to origin `release-*` branch
- Build and release docker image for pg operator.
- Go back to master branch and change the dependency on `apimachinery` back to master branch.

```console
elasticsearch (master) $ glide slow
elasticsearch (master) $ ./hack/make.py
elasticsearch (master) $ git commit -a -m 'Vendor apimachinery release-0.1'
elasticsearch (master) $ git push origin master

elasticsearch (master) $ git checkout release-0.1
elasticsearch (release-0.1) $ git merge master
elasticsearch (release-0.1) $ git tag -fa 0.1.0
elasticsearch (release-0.1) $ git push origin release-0.1 --tags
elasticsearch (release-0.1) $ rm -rf dist
elasticsearch (release-0.1) $ ./hack/docker/es-operator/setup.sh; env APPSCODE_ENV=prod ./hack/docker/es-operator/setup.sh release

elasticsearch (release-0.1) $ git checkout master
elasticsearch (master) $ git commit -a -m 'Start 0.2 development cycle'
elasticsearch (master) $ git push origin master
```

## operator

- Change the `apimachinery` & specific db operator dependency to `release-*` branch
- Revendor `glide slow`
- Make sure master branch compiles, commit any changes due to vendoring & push to origin master
- Create a matching `release-*` branch
- Apply a matching `X.Y.Z` tag
- Push to origin `release-*` branch
- Build and release docker image for pg operator.
- Go back to master branch and change the dependency on `apimachinery` & specific db operator dependency back to master branch.

```console
operator (master) $ glide slow
operator (master) $ ./hack/make.py
operator (master) $ git commit -a -m 'Vendor apimachinery & db operator release-0.1'
operator (master) $ git push origin master

operator (master) $ git checkout release-0.1
operator (release-0.1) $ git merge master
operator (release-0.1) $ git tag -fa 0.1.0
operator (release-0.1) $ git push origin release-0.1 --tags
operator (release-0.1) $ rm -rf dist
operator (release-0.1) $ ./hack/docker/setup.sh; env APPSCODE_ENV=prod ./hack/docker/setup.sh release

operator (release-0.1) $ git checkout master
operator (master) $ git commit -a -m 'Start 0.2 development cycle'
operator (master) $ git push origin master
```

## cli

- Change the apimachinery dependency to `release-*` branch
- Revendor `glide slow`
- Make sure master branch compiles, commit any changes due to vendoring & push to origin master
- Create a matching `release-*` branch
- Apply a matching `X.Y.Z` tag
- Push to origin `release-*` branch
- Build and upload cli binaries to Github.
- Go back to master branch and change the dependency on `apimachinery` back to master branch.

```console
cli (master) $ glide slow
cli (master) $ ./hack/make.py
cli (master) $ git commit -a -m 'Vendor apimachinery release-0.1'
cli (master) $ git push origin master

cli (master) $ git checkout release-0.1
cli (release-0.1) $ git merge master
cli (release-0.1) $ git tag -fa 0.1.0
cli (release-0.1) $ git push origin release-0.1 --tags
cli (release-0.1) $ rm -rf dist
cli (release-0.1) $ env APPSCODE_ENV=prod ./hack/make.py build

cli (release-0.1) $ git checkout master
cli (master) $ git commit -a -m 'Start 0.2 development cycle'
cli (master) $ git push origin master
```
