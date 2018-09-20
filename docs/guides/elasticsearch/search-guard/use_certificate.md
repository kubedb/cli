---
title: Search Guard Use Certificate
menu:
  docs_0.8.0:
    identifier: es-use-certificate-search-guard
    name: Use Certificate
    parent: es-search-guard-elasticsearch
    weight: 20
menu_name: docs_0.8.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Use TLS certificate

If `enableSSL` is set to be true in Elasticsearch object, only HTTPS calls are allowed to database server.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster.
If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace "demo" created

$ kubectl get ns demo
NAME    STATUS  AGE
demo    Active  5s
```

> Note: Yaml files used in this tutorial are stored in [docs/examples/elasticsearch](https://github.com/kubedb/cli/tree/master/docs/examples/elasticsearch) folder in github repository [kubedb/cli](https://github.com/kubedb/cli).

We need an Elasticsearch object in `Running` phase where `enableSSL` is set to be `true`.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: ssl-elasticsearch
  namespace: demo
spec:
  version: "5.6"
  replicas: 2
  enableSSL: true
```

If Elasticsearch object `ssl-elasticsearch` doesn't exists, create it first.

```console
$ kubedb create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0-beta.0/docs/examples/elasticsearch/search-guard/ssl-elasticsearch.yaml
elasticsearch "ssl-elasticsearch" created
```

```console
$ kubedb get es -n demo ssl-elasticsearch
NAME                STATUS    AGE
ssl-elasticsearch   Running   17m
```

## HTTPS request to Elasticsearch

If `enableSSL` is set to be `true`, only HTTPS calls are allowed to Elasticsearch server. If certificates are not provided when Elasticsearch is created,
KubeDB operator will create necessary certificates and use those in Search Guard.

Lets check the certificate, KubeDB created for Elasticsearch `ssl-elasticsearch`.

```console
$ kubectl get secret -n demo ssl-elasticsearch-cert -o yaml
```

```yaml
apiVersion: v1
data:
  client.jks: /u3+7QAAAAIAAAABAAAA...mVv0I52GubpXTAahXDo=
  node.jks: /u3+7QAAAAIAAAABAAAA...pn6opk0qoxabtPTP30c=
  root.jks: /u3+7QAAAAIAAAABAAAA...rjIEWtBA1IMnDcB2JJm5
  root.pem: LS0tLS1CRUdJTiBDRVJU...VElGSUNBVEUtLS0tLQo=
  sgadmin.jks: /u3+7QAAAAIAAAABAAAA...12OXut1U7gYnEyJsBg==
  key_pass: NnRhN3h2
kind: Secret
metadata:
  creationTimestamp: 2018-02-19T09:51:45Z
  labels:
    kubedb.com/kind: Elasticsearch
    kubedb.com/name: ssl-elasticsearch
  name: ssl-elasticsearch-cert
  namespace: demo
  resourceVersion: "754"
  selfLink: /api/v1/namespaces/demo/secrets/ssl-elasticsearch-cert
  uid: 7efdaf31-155a-11e8-a001-42010a8000d5
type: Opaque
```

#### Connect Elasticsearch

In this tutorial, we will expose ClusterIP Service `ssl-elasticsearch` to connect database from local.

```console
$ kubectl expose svc -n demo ssl-elasticsearch --name=ssl-es-exposed --port=9200 --protocol=TCP --type=NodePort
service "ssl-es-exposed" exposed
```

```console
$ kubectl get svc -n demo ssl-es-exposed
NAME             TYPE       CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
ssl-es-exposed   NodePort   10.110.138.210   <none>        9200:30582/TCP   2m
```

Elasticsearch `ssl-elasticsearch` is exposed with following URL

```console
$ minikube service ssl-es-exposed -n demo --https --url
https://192.168.99.100:30582
```

To connect Elasticsearch server securely, now you need to use DNS endpoints of client certificate which are:

- localhost
- *ssl-elasticsearch*.demo.svc

Lets use `ssl-elasticsearch.svc.demo` as host name

```console
curl https://ssl-elasticsearch.demo.svc:30582
```

> Note: You need to set `ssl-elasticsearch.svc.demo` as DNS entry of IP `192.168.99.100` (minikube IP)

As TLS on HTTP layer is enabled, we need to provide root/ca certificate.

To get the root certificate data from Secret, run following command

```console
$ kubectl get secrets -n demo ssl-elasticsearch-cert -o jsonpath='{.data.\root\.pem}' | base64 --decode > root.pem
```

Now try to connect, it will give `Unauthorized` reply. That means, provided certificate works

```console
$ curl https://ssl-elasticsearch.demo.svc:30582 --cacert root.pem
Unauthorized⏎
```

Run following command to get `admin` user password

```console
$ kubectl get secrets -n demo ssl-elasticsearch-auth -o jsonpath='{.data.\ADMIN_PASSWORD}' | base64 -d
uv2io5au⏎
```

Now run following commands to connect to Elasticsearch server in secure mode with basic auth information.

```console
export es_service=https://ssl-elasticsearch.demo.svc:30582
export es_admin_pass=$(kubectl get secrets -n demo ssl-elasticsearch-auth -o jsonpath='{.data.\ADMIN_PASSWORD}' | base64 -d)
curl --user "admin:$es_admin_pass" "$es_service/_cluster/health?pretty" --cacert root.pem
```

```json
{
  "cluster_name" : "ssl-elasticsearch",
  "status" : "green",
  "timed_out" : false,
  "number_of_nodes" : 2,
  "number_of_data_nodes" : 2,
  "active_primary_shards" : 1,
  "active_shards" : 2,
  "relocating_shards" : 0,
  "initializing_shards" : 0,
  "unassigned_shards" : 0,
  "delayed_unassigned_shards" : 0,
  "number_of_pending_tasks" : 0,
  "number_of_in_flight_fetch" : 0,
  "task_max_waiting_in_queue_millis" : 0,
  "active_shards_percent_as_number" : 100.0
}
```

In summary,

- If `enableSSL` is not set, you do not need certificate to validate client, but still you need basic auth.
- If `enableSSL` is set, you need root certificate to validate client.

If certificate Secret is not provided when creating Elasticsearch, one will be created for user.

> Note: Do not need to provide client certificate. Client is verified by valid host name.

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo es/ssl-elasticsearch -p '{"spec":{"doNotPause":false}}' --type="merge"
$ kubectl delete -n demo es/ssl-elasticsearch

$ kubectl patch -n demo drmn/ssl-elasticsearch -p '{"spec":{"wipeOut":true}}' --type="merge"
$ kubectl delete -n demo drmn/ssl-elasticsearch

$ kubectl delete ns demo
```

## Next Steps

- Learn how to [create TLS certificates](/docs/guides/elasticsearch/search-guard/certificate.md).
- Learn how to generate [search-guard configuration](/docs/guides/elasticsearch/search-guard/configuration.md).
- Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
