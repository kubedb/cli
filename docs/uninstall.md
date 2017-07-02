> New to KubeDB? Please start [here](/docs/tutorial.md).

# Uninstall KubeDB
Please follow the steps below to uninstall KubeDB:

1. Delete the deployment and service used for KubeDB operator.
```sh
$ kubectl delete deployment -l app=kubedb -n <operator-namespace>
$ kubectl delete service -l app=kubedb -n <operator-namespace>
```

2. Now, wait several seconds for KubeDB to stop running. To confirm that KubeDB operator pod(s) have stopped running, run:
```sh
$ kubectl get pods --all-namespaces -l app=kubedb
```

3. To keep a copy of your existing KubeDB objects, run:
```sh
kubectl get postgres.kubedb.com --all-namespaces -o yaml > postgres.yaml
kubectl get elastic.kubedb.com --all-namespaces -o yaml > elastic.yaml
kubectl get snapshot.kubedb.com --all-namespaces -o yaml > snapshot.yaml
kubectl get dormant-database.kubedb.com --all-namespaces -o yaml > data.yaml
```

4. To delete existing KubeDB objects from all namespaces, run the following command in each namespace one by one.
```
kubectl delete postgres.kubedb.com --all --cascade=false
kubectl delete elastic.kubedb.com --all --cascade=false
kubectl delete snapshot.kubedb.com --all --cascade=false
kubectl delete dormant-database.kubedb.com --all --cascade=false
```

5. Delete the old TPR-registration.
```sh
kubectl delete thirdpartyresource -l app=kubedb
```
