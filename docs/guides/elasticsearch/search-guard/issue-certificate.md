---
title: Search Guard Certificate
menu:
  docs_0.9.0:
    identifier: es-issue-certificate-search-guard
    name: Issue Certificate
    parent: es-search-guard-elasticsearch
    weight: 25
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Issue TLS Certificates

Search Guard requires certificates to enable TLS. KubeDB creates necessary certificates automatically. However, if you want to use your own certificates, you can provide them through `spec.certificateSecret` field of Elasticsearch object.

This tutorial will show you how to generate certificates for Search Guard and use them with Elasticsearch database.

In KubeDB Elasticsearch, keystore and truststore files in JKS format are used instead of certificates and private keys in PEM format.

KubeDB applies same **truststore**  for both transport layer TLS and REST layer TLS.

But, KubeDB distinguishes between the following types of keystore for security purpose.

- **transport layer keystore** are used to identify and secure traffic between Elasticsearch nodes on the transport layer
- **http layer keystore** are used to identify Elasticsearch clients on the REST and transport layer.
- **sgadmin keystore** are used as admin client that have elevated rights to perform administrative tasks.

## Before You Begin

At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install KubeDB cli on your workstation and KubeDB operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial.

```console
$ kubectl create ns demo
namespace/demo created

$ kubectl get ns demo
NAME    STATUS  AGE
demo    Active  5s
```

You also need to have [*OpenSSL*](https://www.openssl.org/source/) and Java *keytool* for generating all required artifacts.

In order to find out if you have OpenSSL installed, open a terminal and type

```console
$ openssl version
OpenSSL 1.0.2g  1 Mar 2016
```

Make sure it’s version 1.0.1k or higher

And check *keytool* by calling

```console
keytool
```

If already installed, it will print a list of available commands.

To keep generated files separated, open a new terminal and create a directory `/tmp/kubedb/certs`

```console
mkdir -p /tmp/kubedb/certs
cd /tmp/kubedb/certs
```

> Note: YAML files used in this tutorial are stored in [docs/examples/elasticsearch](https://github.com/kubedb/cli/tree/master/docs/examples/elasticsearch) folder in GitHub repository [kubedb/cli](https://github.com/kubedb/cli).

## Generate truststore

First, we need root certificate to sign other server & client certificates. And also this certificate is imported as *truststore*.

You need to follow these steps

1. Get root certificate configuration file

    ```console
    $ wget https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/elasticsearch/search-guard/openssl-config/openssl-ca.ini
    ```

    ```ini
    [ ca ]
    default_ca = CA_default

    [ CA_default ]
    private_key     = root-key.pem
    default_days    = 1000        # how long to certify for
    default_md      = sha256      # use public key default MD
    copy_extensions = copy        # Required to copy SANs from CSR to cert

    [ req ]
    prompt             = no
    default_bits       = 4096
    distinguished_name = ca_distinguished_name

    [ ca_distinguished_name ]
    O  = Elasticsearch Operator
    CN = KubeDB Com. Root CA
    ```

2. Set a password of your keystore and truststore files

    ```console
    $ export KEY_PASS=secret
    ```

    > Note: You need to provide this KEY_PASS in your Secret as `key_pass`

3. Generate private key and certificate

    ```console
    $ openssl req -x509 -config openssl-ca.ini -newkey rsa:4096 -sha256 -nodes -out root.pem -keyout root-key.pem -batch -passin "pass:$KEY_PASS"
    ```

    Here,

    - `root-key.pem` holds Private Key
    - `root.key`holds CA Certificate

4. Finally, import certificate as keystore

    ```console
    $ keytool -import -file root.pem -keystore root.jks -storepass $KEY_PASS -srcstoretype pkcs12 -noprompt
    ```

    Here,

    - `root.jks` is truststore for Elasticsearch

## Generate keystore

Steps to generate certificate and keystore for Elasticsearch

1. Get certificate configuration file
2. Generate private key and certificate signing request (CSR)
3. Sign certificate using root certificate
4. Generate PKCS12 file using root certificate
5. Import PKCS12 as keystore

You need to follow these steps to generate three keystore.

To sign certificate, we need another configuration file.

```console
$ wget https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/elasticsearch/search-guard/openssl-config/openssl-sign.ini
```

```ini
[ ca ]
default_ca = CA_default

[ CA_default ]
base_dir      = .
certificate   = $base_dir/root.pem          # The CA certifcate
private_key   = $base_dir/root-key.pem      # The CA private key
new_certs_dir = $base_dir                   # Location for new certs after signing
database      = $base_dir/index.txt         # Database index file
serial        = $base_dir/serial.txt        # The current serial number
unique_subject = no                         # Set to 'no' to allow creation of several certificates with same subject.

default_days    = 1000        # how long to certify for
default_md      = sha256      # use public key default MD
email_in_dn     = no
copy_extensions = copy        # Required to copy SANs from CSR to cert

[ req ]
default_bits       = 4096
default_keyfile    = root-key.pem
distinguished_name = ca_distinguished_name

[ ca_distinguished_name ]
O  = Elasticsearch Operator
CN = KubeDB Com. Root CA

[ signing_req ]
keyUsage               = digitalSignature, keyEncipherment

[ signing_policy ]
organizationName       = optional
commonName             = supplied
```

Here,

- `certificate` denotes CA certificate path
- `private_key` denotes CA key path

Also, you need to create a `index.txt` file and `serial.txt` file with value `01`

```console
touch index.txt
echo '01' > serial.txt
```

### Node

Following configuration is used to generate CSR for node certificate.

```ini
[ req ]
prompt             = no
default_bits       = 4096
distinguished_name = node_distinguished_name
req_extensions     = node_req_extensions

[ node_distinguished_name ]
O  = Elasticsearch Operator
CN = sg-elasticsearch

[ node_req_extensions ]
keyUsage            = digitalSignature, keyEncipherment
extendedKeyUsage    = serverAuth, clientAuth
subjectAltName      = @alternate_names

[ alternate_names ]
DNS.1 = localhost
RID.1 = 1.2.3.4.5.5
```

Here,

- `RID.1=1.2.3.4.5.5` is used in node certificate. All certificates with registeredID `1.2.3.4.5.5` is considered as valid certificate for transport layer.

Now run following commands

```console
$ wget https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/elasticsearch/search-guard/openssl-config/openssl-node.ini
$ openssl req -config openssl-node.ini -newkey rsa:4096 -sha256 -nodes -out node-csr.pem -keyout node-key.pem
$ openssl ca -config openssl-sign.ini -batch -policy signing_policy -extensions signing_req -out node.pem -infiles node-csr.pem
$ openssl pkcs12 -export -certfile root.pem -inkey node-key.pem -in node.pem -password "pass:$KEY_PASS" -out node.pkcs12
$ keytool -importkeystore -srckeystore node.pkcs12  -storepass $KEY_PASS  -srcstoretype pkcs12 -srcstorepass $KEY_PASS  -destkeystore node.jks -deststoretype pkcs12
```

Generated `node.jks` will be used as keystore for transport layer TLS.

### Client

Following configuration is used to generate CSR for client certificate.

```ini
[ req ]
prompt             = no
default_bits       = 4096
distinguished_name = client_distinguished_name
req_extensions     = client_req_extensions

[ client_distinguished_name ]
O  = Elasticsearch Operator
CN = sg-elasticsearch

[ client_req_extensions ]
keyUsage            = digitalSignature, keyEncipherment
extendedKeyUsage    = serverAuth, clientAuth
subjectAltName      = @alternate_names

[ alternate_names ]
DNS.1 = localhost
DNS.2 = sg-elasticsearch.demo.svc
```

Here,

- `sg-elasticsearch` is used as a Common Name so that host `sg-elasticsearch` is verified as valid Client.

Now run following commands

```console
$ wget https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/elasticsearch/search-guard/openssl-config/openssl-client.ini
$ openssl req -config openssl-client.ini -newkey rsa:4096 -sha256 -nodes -out client-csr.pem -keyout client-key.pem
$ openssl ca -config openssl-sign.ini -batch -policy signing_policy -extensions signing_req -out client.pem -infiles client-csr.pem
$ openssl pkcs12 -export -certfile root.pem -inkey client-key.pem -in client.pem -password "pass:$KEY_PASS" -out client.pkcs12
$ keytool -importkeystore -srckeystore client.pkcs12  -storepass $KEY_PASS  -srcstoretype pkcs12 -srcstorepass $KEY_PASS  -destkeystore client.jks -deststoretype pkcs12
```

Generated `client.jks` will be used as keystore for http layer TLS.

### sgadmin

Following configuration is used to generate CSR for sgadmin certificate.

```ini
[ req ]
prompt             = no
default_bits       = 4096
distinguished_name = sgadmin_distinguished_name
req_extensions     = sgadmin_req_extensions

[ sgadmin_distinguished_name ]
O  = Elasticsearch Operator
CN = sgadmin

[ sgadmin_req_extensions ]
keyUsage            = digitalSignature, keyEncipherment
extendedKeyUsage    = serverAuth, clientAuth
subjectAltName      = @alternate_names

[ alternate_names ]
DNS.1 = localhost
```

Here,

- `sgadmin` is used as Common Name. Because in searchguard, certificate with `sgadmin` common name is considered as admin certificate.

Now run following commands

```console
$ wget https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/elasticsearch/search-guard/openssl-config/openssl-sgadmin.ini
$ openssl req -config openssl-sgadmin.ini -newkey rsa:4096 -sha256 -nodes -out sgadmin-csr.pem -keyout sgadmin-key.pem
$ openssl ca -config openssl-sign.ini -batch -policy signing_policy -extensions signing_req -out sgadmin.pem -infiles sgadmin-csr.pem
$ openssl pkcs12 -export -certfile root.pem -inkey sgadmin-key.pem -in sgadmin.pem -password "pass:$KEY_PASS" -out sgadmin.pkcs12
$ keytool -importkeystore -srckeystore sgadmin.pkcs12  -storepass $KEY_PASS  -srcstoretype pkcs12 -srcstorepass $KEY_PASS  -destkeystore sgadmin.jks -deststoretype pkcs12
```

Generated `sgadmin.pkcs12` will be used as keystore for admin usage.

## Create Secret

Now create a Secret with these certificates to use in your Elasticsearch object.

```console
$ kubectl create secret -n demo generic sg-elasticsearch-cert \
                --from-file=root.pem \
                --from-file=root.jks \
                --from-file=node.jks \
                --from-file=client.jks \
                --from-file=sgadmin.jks \
                --from-literal=key_pass=$KEY_PASS

secret/sg-elasticsearch-cert created
```

> Note: `root.pem` is added in Secret so that user can use these to connect Elasticsearch

Use this Secret `sg-elasticsearch-cert` in your Elasticsearch object.

## Create a Elasticsearch database

Below is the Elasticsearch object created in this tutorial.

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Elasticsearch
metadata:
  name: sg-elasticsearch
  namespace: demo
spec:
  version: "6.3-v1"
  authPlugin: "SearchGuard"
  enableSSL: true
  certificateSecret:
    secretName: sg-elasticsearch-cert
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
```

Here,

- `spec.certificateSecret` specifies Secret with certificates those will be used in Elasticsearch database.

Create example above with following command

```console
$ kubectl create -f https://raw.githubusercontent.com/kubedb/cli/0.9.0/docs/examples/elasticsearch/search-guard/sg-elasticsearch.yaml
elasticsearch.kubedb.com/sg-elasticsearch created
```

KubeDB operator sets the `status.phase` to `Running` once the database is successfully created.

```console
$ kubectl get es -n demo sg-elasticsearch -o wide
NAME               VERSION   STATUS    AGE
sg-elasticsearch   6.3-v1    Running   1m
```

## Cleaning up

To cleanup the Kubernetes resources created by this tutorial, run:

```console
$ kubectl patch -n demo es/sg-elasticsearch -p '{"spec":{"terminationPolicy":"WipeOut"}}' --type="merge"
$ kubectl delete -n demo es/sg-elasticsearch

$ kubectl delete ns demo
```

## Next Steps

- Learn how to use TLS certificates to connect Elasticsearch from [here](/docs/guides/elasticsearch/search-guard/use-tls.md).
- Learn how to generate [search-guard configuration](/docs/guides/elasticsearch/search-guard/configuration.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
