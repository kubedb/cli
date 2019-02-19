---
title: Using Kibana with KubeDB Elasticsearch
menu:
  docs_0.9.0:
    identifier: es-kibana-connect
    name: Use Kibana
    parent: es-kibana
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Using Kibana with KubeDB Elasticsearch

This tutorial will show you how to connect Kibana with an Elasticsearch cluster deployed with KubeDB.

If you don't know how to use Kibana, please visit [here](https://www.elastic.co/guide/en/kibana/current/introduction.html).

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

As KubeDB uses [Search Guard](https://search-guard.com/) plugin for authentication and authorization, you have to know how to configure Search Guard for both Elasticsearch cluster and Kibana. If you don't know, please visit [here](https://docs.search-guard.com/latest/main-concepts).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace/demo created

$ kubectl get ns demo
NAME    STATUS  AGE
demo    Active  5s
```

> Note: YAML files used in this tutorial are stored in [docs/examples/elasticsearch](https://github.com/kubedb/cli/tree/master/docs/examples/elasticsearch) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Overview

At first, we will create some necessary Search Guard configuration and roles to give a user access to an Elasticsearch cluster from Kibana. We will create a secret with this configuration files. Then we will provide this secret in `spec.databaseSecret` field of Elasticsearch crd so that our Elasticsearch cluster start with this configuration. We will also configure Elasticsearch cluster with a [custom configuration](/docs/guides/elasticsearch/custom-config/overview.md) file.

Then, we will deploy Kibana with Search Guard plugin installed. We will configure Kibana to connect with our Elasticsearch cluster.

Finally, we will perform some operation from Kibana UI to ensure that Kibana is working well with our Elasticsearch cluster.

For this tutorial, we will use Elasticsearch 6.3.0 with Search Guard plugin 23.1 and Kibana 6.3.0 with Search Guard plugin 14 installed.

## Deploy Elasticsearch Cluster

Let's create necessary Search Guard configuration files. Here, we will create two users `admin` and `kibanauser`. User `admin` will have all permissions on the cluster and user `kibanauser` will have some limited permission. Here, are the contents of Search Guard configuration files,

**sg_action_groups.yml:**

```yaml
UNLIMITED:
  readonly: true
  permissions:
    - "*"

###### INDEX LEVEL ######

INDICES_ALL:
  readonly: true
  permissions:
    - "indices:*"

###### CLUSTER LEVEL ######
CLUSTER_MONITOR:
  readonly: true
  permissions:
    - "cluster:monitor/*"

CLUSTER_COMPOSITE_OPS_RO:
  readonly: true
  permissions:
    - "indices:data/read/mget"
    - "indices:data/read/msearch"
    - "indices:data/read/mtv"
    - "indices:data/read/coordinate-msearch*"
    - "indices:admin/aliases/exists*"
    - "indices:admin/aliases/get*"
    - "indices:data/read/scroll"

CLUSTER_COMPOSITE_OPS:
  readonly: true
  permissions:
    - "indices:data/write/bulk"
    - "indices:admin/aliases*"
    - "indices:data/write/reindex"
    - CLUSTER_COMPOSITE_OPS_RO
```

**sg_roles.yaml:**

```yaml
sg_all_access:
  readonly: true
  cluster:
    - UNLIMITED
  indices:
    '*':
      '*':
        - UNLIMITED
  tenants:
    admin_tenant: RW

# For the kibana user
sg_kibana_user:
  readonly: true
  cluster:
      - CLUSTER_MONITOR
      - CLUSTER_COMPOSITE_OPS
      - cluster:admin/xpack/monitoring*
      - indices:admin/template*
  indices:
    '*':
      '*':
        - INDICES_ALL
```

**sg_internal_users.yml:**

```yaml
#password is: admin@secret
admin:
  readonly: true
  hash: $2y$12$skma87wuFFtxtGWegeAiIeTtUH1nnOfIRZzwwhBlzXjg0DdM4gLeG
  roles:
    - admin

#password is: kibana@secret
kibanauser:
  readonly: true
  hash: $2y$12$dk2UrPTjhgCRbFOm/gThX.aJ47yH0zyQcYEuWiNiyw6NlVmeOjM7a
  roles:
    - kibanauser
```

Here, we have used `admin@secret` password for `admin` user and  `kibana@secret` password for `kibanauser` user. You can use `htpasswd` to generate the bcrypt encrypted password hashes.

```console
$htpasswd -bnBC 12 "" <password_here>| tr -d ':\n'
```

**sg_roles_mapping.yml:**

```yaml
sg_all_access:
  readonly: true
  backendroles:
    - admin

sg_kibana_user:
  readonly: true
  backendroles:
    - kibanauser
```

**sg_config.yml:**

```yaml
searchguard:
  dynamic:
    authc:
      kibana_auth_domain:
        enabled: true
        order: 0
        http_authenticator:
          type: basic
          challenge: false
        authentication_backend:
          type: internal
      basic_internal_auth_domain: 
        http_enabled: true
        transport_enabled: true
        order: 1
        http_authenticator:
          type: basic
          challenge: true
        authentication_backend:
          type: internal
```

Now, create a secret with these Search Guard configuration files.

```console
 $ kubectl create secret generic -n demo es-auth \
	--from-literal=ADMIN_USERNAME=admin \
	--from-literal=ADMIN_PASSWORD=admin@secret \
	--from-file=./sg_action_groups.yml \
	--from-file=./sg_config.yml \
	--from-file=./sg_internal_users.yml \
	--from-file=./sg_roles_mapping.yml \
	--from-file=./sg_roles.yml
secret/es-auth created
```

Verify that the secret has desired configuration files,

```yaml
$ kubectl get secret -n demo es-auth -o yaml
apiVersion: v1
data:
  sg_action_groups.yml: <base64 encoded content>
  sg_config.yml: <base64 encoded content>
  sg_internal_users.yml: <base64 encoded content>
  sg_roles.yml: <base64 encoded content>
  sg_roles_mapping.yml: <base64 encoded content>
kind: Secret
metadata:
  ...
  name: es-auth
  namespace: demo
  ...
type: Opaque
```

As we are using Search Guard plugin for authentication, we need to ensure that `x-pack` security is not enabled. We will ensure that by providing `xpack.security.enabled: false` in `common-config.yml` file and we will use this file to configure our Elasticsearch cluster. We will also configure `searchguard.restapi` to ensure that `kibanauser` can use REST API on the cluster.

 Content of `common-config.yml`,

```yaml
xpack.security.enabled: false
searchguard.restapi.roles_enabled: ["sg_all_access","sg_kibana_user"]
```

Create a ConfigMap using this file,

```console
$ kubectl create configmap -n demo es-custom-config \
                        --from-file=./common-config.yml
configmap/es-custom-config created
```

Verify that the ConfigMap has desired configuration,

```yaml
$ kubectl get configmap -n demo es-custom-config -o yaml
apiVersion: v1
data:
  common-config.yaml: |-
    xpack.security.enabled: false
    searchguard.restapi.roles_enabled: ["sg_all_access","sg_kibana_user"]
kind: ConfigMap
metadata:
  creationTimestamp: 2018-08-18T06:53:04Z
  name: es-custom-config
  namespace: demo
  resourceVersion: "12171"
  selfLink: /api/v1/namespaces/demo/configmaps/es-custom-config
  uid: 5b2adaeb-a2b3-11e8-ba38-080027975c84
```

Now, create Elasticsearch crd specifying  `spec.databaseSecret` and `spec.configSource` field.

```console
$ kubectl apply -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/elasticsearch/kibana/es-kibana-demo.yaml
elasticsearch.kubedb.com/es-kibana-demo created
```

Below is the YAML for the Elasticsearch crd we just created.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: es-kibana-demo
  namespace: demo
spec:
  version: "6.3.0-v1"
  replicas: 1
  databaseSecret:
    secretName: es-auth
  configSource:
    configMap:
      name: es-custom-config
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
```

Now, wait for few minutes. KubeDB will create necessary secrets, services, and statefulsets.

Check resources created in demo namespace by KubeDB,

```console
$ kubectl get all -n demo -l=kubedb.com/name=es-kibana-demo
NAME                   READY     STATUS    RESTARTS   AGE
pod/es-kibana-demo-0   1/1       Running   0          39s

NAME                            TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/es-kibana-demo          ClusterIP   10.104.1.206    <none>        9200/TCP   44s
service/es-kibana-demo-master   ClusterIP   10.111.58.230   <none>        9300/TCP   44s

NAME                              DESIRED   CURRENT   AGE
statefulset.apps/es-kibana-demo   1         1         42s
```

Once everything is created, Elasticsearch will go to Running state. Check that Elasticsearch is in running state.

```console
$ kubectl get es -n demo es-kibana-demo
NAME             VERSION    STATUS    AGE
es-kibana-demo   6.3.0-v1   Running   1m
```

Now, check elasticsearch log to see if the cluster is ready to accept requests,

```console
$ kubectl logs -n demo es-kibana-demo-0 -f
...
Starting runit...
...
Elasticsearch Version: 6.3.0
Search Guard Version: 6.3.0-23.1
Connected as CN=sgadmin,O=Elasticsearch Operator
Contacting elasticsearch cluster 'elasticsearch' and wait for YELLOW clusterstate ...
Clustername: es-kibana-demo
Clusterstate: GREEN
Number of nodes: 1
Number of data nodes: 1
...
Done with success
...
```

Once you see `Done with success` success line in the log, the cluster is ready to accept requests. Now, it is time to connect with Kibana.

## Deploy Kibana

In order to connect the Elasticsearch cluster that we have just deployed, we need to configure `kibana.yml` with appropriate configuration.

KubeDB has created a service with name`es-kibana-demo` in `demo` namespace for the Elasticsearch cluster. We will use this service in `elasticsearch.url` field. Kibana will use this service to connect with the Elasticsearch cluster.

Let's, configure `kibana.yml` as below,

```yaml
xpack.security.enabled: false
server.host: 0.0.0.0

elasticsearch.url: "http://es-kibana-demo.demo.svc:9200"
elasticsearch.username: "kibanauser"
elasticsearch.password: "kibana@secret"

searchguard.auth.type: "basicauth"
searchguard.cookie.secure: false

```

Notice the `elasticsearch.username` and `elasticsearch.password` field. Kibana will connect to Elasticsearch cluster with this credentials. They must match with the credentials we have provided in `sg_internal_users.yml` file while creating the cluster.

Now, create a ConfigMap with `kibana.yml` file. We will mount this ConfigMap in Kibana deployment so that Kibana starts with this configuration.

```console
$ kubectl create configmap -n demo kibana-config \
                        --from-file=./kibana.yml
configmap/kibana-config created
```

Finally, deploy Kibana deployment,
```console
$ kubectl apply -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/elasticsearch/kibana/kibana-deployment.yaml
deployment.apps/kibana created
```

Below is the YAML for the Kibana deployment we just created.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kibana
  namespace: demo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kibana
  template:
    metadata:
      labels:
        app: kibana
    spec:
      containers:
      - name: kibana
        image: kubedb/kibana:6.3.0
        volumeMounts:
        - name:  kibana-config
          mountPath: /usr/share/kibana/config
      volumes:
      - name:  kibana-config
        configMap:
          name: kibana-config
```

Now, wait for few minutes. Let the Kibana pod  go in`Running` state. Check pod is in `Running` using this command,

```console
 $ kubectl get pods -n demo -l app=kibana
NAME                      READY     STATUS    RESTARTS   AGE
kibana-84b8cbcf7c-mg699   1/1       Running   0          3m
```

Now, watch the Kibana pod's log to see if Kibana is ready to access,

```console
$ kubectl logs -n demo kibana-84b8cbcf7c-mg699 -f
...
{"type":"log","@timestamp":"2018-08-18T07:22:16Z","tags":["listening","info"],"pid":1,"message":"Server running at http://0.0.0.0:5601"}
```

Once you see `"message":"Server running at http://0.0.0.0:5601"` in the log, Kibana is ready. Now it is time to access Kibana UI.

Kibana is running on port `5601` in of `kibana-84b8cbcf7c-mg699` pod. In order to access Kibana UI from outside of the cluster, we will use [port forwarding](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster).

First, open a new terminal and run,

```console
$ kubectl port-forward -n demo kibana-84b8cbcf7c-mg699 5601
Forwarding from 127.0.0.1:5601 -> 5601
Forwarding from [::1]:5601 -> 5601
```

Now, open `localhost:5601` in your browser. When you open the address, you will be greeted with Search Guard login UI.

![Search Guard Login UI](/docs/images/elasticsearch/kibana/search-guard-login-ui.png)

Login with following credentials: `username: kibanauser` and `password: kibana@secret`.

After login, you will be redirected to Kibana Home UI.

![Kibana Home](/docs/images/elasticsearch/kibana/kibana-home.png)

Now, it is time to perform some operations on our cluster from the Kibana UI.

## Use Kibana

We can use Dev Tool's console of Kibana UI to create Index and insert data in the index. Let's create an Index,

```json
PUT /shakespeare
{
 "mappings": {
  "doc": {
   "properties": {
    "speaker": {"type": "keyword"},
    "play_name": {"type": "keyword"},
    "line_id": {"type": "integer"},
    "speech_number": {"type": "integer"}
   }
  }
 }
}
```

Now, insert some demo data in the Index,

```json
//  demo data-1
POST /shakespeare/doc
{
    "index": {
        "_index": "shakespeare",
        "_id": 1
    },
    "type": "scene",
    "line_id": 2,
    "play_name": "Henry IV",
    "speech_number": "",
    "line_number": "",
    "speaker": "",
    "text_entry": "SCENE I. London. The palace."
}

// demo data-2

POST /shakespeare/doc
{
    "index": {
        "_index": "shakespeare",
        "_id": 2
    },
    "type": "line",
    "line_id": 3,
    "play_name": "Henry IV",
    "speech_number": "",
    "line_number": "",
    "speaker": "",
    "text_entry": "Enter KING HENRY, LORD JOHN OF LANCASTER, the EARL of WESTMORELAND, SIR WALTER BLUNT, and others"
}
```

Now, let's create index_pattern.

![Create Index_Pattern](/docs/images/elasticsearch/kibana/kibana-create-index.png)

Once we have created an index_pattern, we can use the Discovery UI.

![](/docs/images/elasticsearch/kibana/kibana-discovery-ui.png)

## Cleanup

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo es/es-kibana-demo -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"

$ kubectl delete -n demo es/es-kibana-demo

$ kubectl delete  -n demo configmap/es-custom-config

$ kubectl delete -n demo configmap/kibana-config

$ kubectl delete -n demo deployment/kibana

$ kubectl delete ns demo
```

To uninstall KubeDB follow this [guide](/docs/setup/uninstall.md).
