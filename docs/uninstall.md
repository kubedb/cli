---
title: KubeDB Uninstall
menu:
  docs_0.7.1:
    identifier: uninstall-kubedb
    name: Uninstall
    parent: getting-started
    weight: 30
menu_name: docs_0.7.1
section_menu_id: getting-started
url: /docs/0.7.1/getting-started/uninstall/
aliases:
  - /docs/0.7.1/uninstall/
---

> New to KubeDB? Please start [here](/docs/tutorials/README.md).

# Uninstall KubeDB
Please follow the steps below to uninstall KubeDB:

1. Delete the deployment and service used for KubeDB operator.
```console
$ kubectl delete deployment -l app=kubedb -n <operator-namespace>
$ kubectl delete service -l app=kubedb -n <operator-namespace>

# Delete RBAC objects, if --rbac flag was used.
$ kubectl delete serviceaccount -l app=kubedb -n <operator-namespace>
$ kubectl delete clusterrolebindings -l app=kubedb -n <operator-namespace>
$ kubectl delete clusterrole -l app=kubedb -n <operator-namespace>
```

2. Now, wait several seconds for KubeDB to stop running. To confirm that KubeDB operator pod(s) have stopped running, run:
```console
$ kubectl get pods --all-namespaces -l app=kubedb
```

3. To keep a copy of your existing KubeDB objects, run:
```console
kubectl get postgres.kubedb.com --all-namespaces -o yaml > postgres.yaml
kubectl get elasticsearch.kubedb.com --all-namespaces -o yaml > elastic.yaml
kubectl get snapshot.kubedb.com --all-namespaces -o yaml > snapshot.yaml
kubectl get dormant-database.kubedb.com --all-namespaces -o yaml > data.yaml
```

4. To delete existing KubeDB objects from all namespaces, run the following command in each namespace one by one.
```
kubectl delete postgres.kubedb.com --all --cascade=false
kubectl delete elasticsearch.kubedb.com --all --cascade=false
kubectl delete snapshot.kubedb.com --all --cascade=false
kubectl delete dormant-database.kubedb.com --all --cascade=false
```

5. Delete the old CRD-registration.
```console
kubectl delete thirdpartyresource -l app=kubedb
```
