---
title: KubeDB Uninstall
menu:
  docs_0.9.0-rc.0:
    identifier: uninstall-kubedb
    name: Uninstall
    parent: setup
    weight: 20
menu_name: docs_0.9.0-rc.0
section_menu_id: setup
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Uninstall KubeDB

To uninstall KubeDB operator, run the following command:

```console
$ curl -fsSL https://raw.githubusercontent.com/kubedb/cli/0.9.0-rc.0/hack/deploy/kubedb.sh \
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

The above command will leave the KubeDB crd objects as-is. If you wish to **nuke** all KubeDB crd objects, also pass the `--purge` flag. This will keep a copy of KubeDB crd objects in your current directory.


- Now, wait several seconds for KubeDB to stop running. To confirm that KubeDB operator pod(s) have stopped running, run:

    ```console
    $ kubectl get pods --all-namespaces -l app=kubedb
    ```

- To keep a copy of your existing KubeDB objects, run:

    ```console
    kubectl get postgres.kubedb.com --all-namespaces -o yaml > postgres.yaml
    kubectl get elasticsearch.kubedb.com --all-namespaces -o yaml > elasticsearch.yaml
    kubectl get memcached.kubedb.com --all-namespaces -o yaml > memcached.yaml
    kubectl get mongodb.kubedb.com --all-namespaces -o yaml > mongodb.yaml
    kubectl get mysql.kubedb.com --all-namespaces -o yaml > mysql.yaml
    kubectl get redis.kubedb.com --all-namespaces -o yaml > redis.yaml
    kubectl get snapshot.kubedb.com --all-namespaces -o yaml > snapshot.yaml
    kubectl get dormant-database.kubedb.com --all-namespaces -o yaml > data.yaml
    ```

- To delete existing KubeDB objects from all namespaces, run the following command in each namespace one by one.

    ```console
    kubectl delete postgres.kubedb.com --all --cascade=false
    kubectl delete elasticsearch.kubedb.com --all --cascade=false
    kubectl delete memcached.kubedb.com --all --cascade=false
    kubectl delete mongodb.kubedb.com --all --cascade=false
    kubectl delete mysql.kubedb.com --all --cascade=false
    kubectl delete redis.kubedb.com --all --cascade=false
    kubectl delete snapshot.kubedb.com --all --cascade=false
    kubectl delete dormant-database.kubedb.com --all --cascade=false
    ```

- Delete the old CRD-registration.

    ```console
    kubectl delete crd -l app=kubedb
    ```
