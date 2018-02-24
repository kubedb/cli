---
title: KubeDB Uninstall
menu:
  docs_0.8.0-beta.0:
    identifier: uninstall-kubedb
    name: Uninstall
    parent: setup
    weight: 20
menu_name: docs_0.8.0-beta.0
section_menu_id: setup
---

> New to KubeDB? Please start [here](/docs/guides/README.md).

# Uninstall KubeDB
Please follow the steps below to uninstall KubeDB:

- Delete the deployment and service used for KubeDB operator.

```console
$ curl -fsSL https://raw.githubusercontent.com/kubedb/cli/0.8.0-alpha.0/hack/deploy/kubedb.sh \
    | bash -s -- --uninstall [--namespace=NAMESPACE]

+ kubectl delete deployment -l app=kubedb -n kube-system
deployment "kubedb-operator" deleted
+ kubectl delete service -l app=kubedb -n kube-system
service "kubedb-operator" deleted
+ kubectl delete serviceaccount -l app=kubedb -n kube-system
No resources found
+ kubectl delete clusterrolebindings -l app=kubedb -n kube-system
No resources found
+ kubectl delete clusterrole -l app=kubedb -n kube-system
No resources found
```

- Now, wait several seconds for KubeDB to stop running. To confirm that KubeDB operator pod(s) have stopped running, run:

```console
$ kubectl get pods --all-namespaces -l app=kubedb
```

- To keep a copy of your existing KubeDB objects, run:

```console
kubectl get postgres.kubedb.com --all-namespaces -o yaml > postgres.yaml
kubectl get elasticsearch.kubedb.com --all-namespaces -o yaml > elastic.yaml
kubectl get snapshot.kubedb.com --all-namespaces -o yaml > snapshot.yaml
kubectl get dormant-database.kubedb.com --all-namespaces -o yaml > data.yaml
```

- To delete existing KubeDB objects from all namespaces, run the following command in each namespace one by one.

```console
kubectl delete postgres.kubedb.com --all --cascade=false
kubectl delete elasticsearch.kubedb.com --all --cascade=false
kubectl delete snapshot.kubedb.com --all --cascade=false
kubectl delete dormant-database.kubedb.com --all --cascade=false
```

- Delete the old CRD-registration.

```console
kubectl delete crd -l app=kubedb
```
