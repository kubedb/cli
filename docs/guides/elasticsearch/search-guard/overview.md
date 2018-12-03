---
title: Search Guard
menu:
  docs_0.9.0-rc.1:
    identifier: es-search-guard-search-guard
    name: Overview
    parent: es-search-guard-elasticsearch
    weight: 10
menu_name: docs_0.9.0-rc.1
section_menu_id: guides
---

> New to KubeDB? Please start [here](/docs/concepts/README.md).

# Search Guard

[Search Guard(®)](https://github.com/floragunncom/search-guard) is an Elasticsearch plugin that offers encryption, authentication, and authorization.
It supports fine grained role-based access control to clusters, indices, documents and fields.

- Search Guard authenticates the credentials against the configured authentication backend(s).
- Search Guard authorizes the user by retrieving a list of the user’s roles from the configured authorization backend

## TLS certificates

Search Guard relies heavily on the use of TLS, both for the REST and the transport layer of Elasticsearch. TLS is configured in the `elasticsearch.yml` file of Elasticsearch installation.

Following keys are used to configure location of keystore and truststore files.

Transport layer TLS

| Name                                              | Description                                                                   |
|---------------------------------------------------|:------------------------------------------------------------------------------|
| searchguard.ssl.transport.keystore_filepath       | Path to the keystore file, relative to the config/ directory (mandatory)      |
| searchguard.ssl.transport.keystore_password       | Keystore password                                                             |
| searchguard.ssl.transport.truststore_filepath     | Path to the truststore file, relative to the config/ directory (mandatory)    |
| searchguard.ssl.transport.truststore_password     | Truststore password                                                           |

REST layer TLS

| Name                                          | Description                                                                           |
|-----------------------------------------------|:--------------------------------------------------------------------------------------|
| searchguard.ssl.http.enabled                  | Whether to enable TLS on the REST layer or not                                        |
| searchguard.ssl.http.keystore_filepath        | Path to the keystore file, relative to the config/ directory (mandatory)              |
| searchguard.ssl.http.keystore_password        | Keystore password                                                                     |
| searchguard.ssl.http.truststore_filepath      | Path to the truststore file, relative to the config/ directory (mandatory)            |
| searchguard.ssl.http.truststore_password      | Truststore password                                                                   |


> Note: KubeDB Elasticsearch is configured with keystore and truststore files in JKS format

#### Configuring Admin certificates

Admin certificates are regular client certificates that have elevated rights to perform administrative tasks. You need an admin certificate to
change the Search Guard configuration via the *sgadmin* command line tool. Admin certificates are configured in `elasticsearch.yml` by simply stating their DN(s).

```yaml
searchguard.authcz.admin_dn:
  - CN=sgadmin, O=Elasticsearch Operator
```

#### Client authentication

With TLS client authentication enabled, REST clients can send a TLS certificate with the HTTP request to provide identity information to Search Guard.

- You can provide an admin certificate when using the REST API.
- You can provide Basic Auth with client certificates.

> Note: Search Guard accepts TLS client certificates if they are sent, but does not enforce them.

## Next Steps

- Learn how to [create TLS certificates](/docs/guides/elasticsearch/search-guard/issue-certificate.md).
- Learn how to [use TLS certificates](/docs/guides/elasticsearch/search-guard/use-tls.md) to connect Elasticsearch.
- Learn how to generate [search-guard configuration](/docs/guides/elasticsearch/search-guard/configuration.md).
- Want to hack on KubeDB? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
