/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import "kubedb.dev/apimachinery/apis/kafka"

const (
	Finalizer = kafka.GroupName + "/finalizer"
)

const (
	LabelRole = kafka.GroupName + "/role"
	RoleStats = "stats"

	ComponentKafka   = "kafka"
	DefaultStatsPath = "/metrics"
)

// ConnectCluster constants

const (
	ConnectClusterUser            = "connect"
	ConnectClusterContainerName   = "connect-cluster"
	ConnectClusterModeEnv         = "CONNECT_CLUSTER_MODE"
	ConnectClusterPrimaryPortName = "primary"
	ConnectClusterPortName        = "connect"
	ConnectClusterRESTPort        = 8083
	ConnectClusterUserEnv         = "CONNECT_CLUSTER_USER"
	ConnectClusterPasswordEnv     = "CONNECT_CLUSTER_PASSWORD"

	ConnectClusterOperatorVolumeConfig = "connect-operator-config"
	ConnectClusterCustomVolumeConfig   = "connect-custom-config"
	ConnectorPluginsVolumeName         = "connector-plugins"
	ConnectClusterAuthSecretVolumeName = "connect-cluster-auth"
	ConnectClusterOffsetFileDirName    = "connect-stand-offset"

	ConnectClusterGroupID                     = "group.id"
	ConnectClusterPluginPath                  = "plugin.path"
	ConnectClusterRestAdvertisedHostName      = "rest.advertised.host.name"
	ConnectClusterRestAdvertisedPort          = "rest.advertised.port"
	ConnectClusterOffsetStorage               = "offset.storage.file.filename"
	ConnectClusterKeyConverter                = "key.converter"
	ConnectClusterValueConverter              = "value.converter"
	ConnectClusterKeyConverterSchemasEnable   = "key.converter.schemas.enable"
	ConnectClusterValueConverterSchemasEnable = "value.converter.schemas.enable"
	ConnectClusterJsonConverterName           = "org.apache.kafka.connect.json.JsonConverter"
	ConnectClusterStatusStorageTopic          = "status.storage.topic"
	ConnectClusterConfigStorageTopic          = "config.storage.topic"
	ConnectClusterOffsetStorageTopic          = "offset.storage.topic"
	ConnectClusterBootstrapServers            = "bootstrap.servers"
	ConnectClusterStatusStorageTopicName      = "connect-status"
	ConnectClusterConfigStorageTopicName      = "connect-configs"
	ConnectClusterOffsetStorageTopicName      = "connect-offsets"

	ConnectClusterListeners               = "listeners"
	ConnectClusterServerCertsVolumeName   = "server-certs"
	KafkaClientCertVolumeName             = "kafka-client-ssl"
	KafkaClientKeystoreKey                = "client.keystore.jks"
	KafkaClientTruststoreKey              = "client.truststore.jks"
	ConnectClusterBasicAuthKey            = "rest.extension.classes"
	KafkaConnectRestAdvertisedListener    = "rest.advertised.listener"
	ConnectClusterKeyPassword             = "listeners.https.ssl.key.password"
	ConnectClusterKeystorePassword        = "listeners.https.ssl.keystore.password"
	ConnectClusterKeystoreLocation        = "listeners.https.ssl.keystore.location"
	ConnectClusterTruststoreLocation      = "listeners.https.ssl.truststore.location"
	ConnectClusterTruststorePassword      = "listeners.https.ssl.truststore.password"
	ConnectClusterClientAuthentication    = "listeners.https.ssl.client.authentication"
	ConnectClusterIdentificationAlgorithm = "listeners.https.ssl.endpoint.identification.algorithm"
	ConnectClusterBasicAuthValue          = "org.apache.kafka.connect.rest.basic.auth.extension.BasicAuthSecurityRestExtension"

	ConnectClusterOffsetFileDir        = "/var/log/connect"
	ConnectClusterServerCertVolumeDir  = "/var/private/ssl"
	ConnectClusterPluginPathDir        = "/opt/kafka/libs"
	ConnectClusterAuthSecretVolumePath = "/var/private/basic-auth"
	KafkaClientCertDir                 = "/var/private/kafka-client-ssl"
	ConnectClusterOffsetFileName       = "/var/log/connect/connect.offsets"
	ConnectorPluginsVolumeDir          = "/opt/kafka/libs/connector-plugins"
	ConnectClusterCustomConfigPath     = "/opt/kafka/config/connect-custom-config"
	ConnectClusterOperatorConfigPath   = "/opt/kafka/config/connect-operator-config"
	KafkaClientKeystoreLocation        = "/var/private/kafka-client-ssl/client.keystore.jks"
	KafkaClientTruststoreLocation      = "/var/private/kafka-client-ssl/client.truststore.jks"
)

// SchemaRegistry constants

const (
	SchemaRegistryPrimaryPortName = "primary"
	SchemaRegistryPortName        = "registry"
	ApicurioRegistryRESTPort      = 8080
	SchemaRegistryContainerName   = "schema-registry"
	SchemaRegistryConfigFileName  = "application.properties"

	SchemaRegistryStorageBackendTypeMemory = "mem"
	SchemaRegistryStorageBackendTypeKafka  = "kafkasql"
	SchemaRegistryStorageBackendTypeSQL    = "sql"

	SchemaRegistryOperatorVolumeConfig = "registry-operator-config"
	SchemaRegistryOperatorConfigPath   = "/deployments/config"
)

// RestProxy constants

const (
	RestProxyPrimaryPortName      = "primary"
	RestProxyPortName             = "restproxy"
	RestProxyRESTPort             = 8082
	RestProxyContainerName        = "rest-proxy"
	RestProxyOperatorVolumeConfig = "rest-proxy-operator-config"
	RestProxyOperatorConfigPath   = "/opt/karapace/config"

	RestProxyKarapaceLogLevel          = "log_level"
	RestProxyKarapaceLogLevelWarning   = "WARNING"
	RestProxyKarapaceLogLevelInfo      = "INFO"
	RestProxyKarapaceLogLevelDebug     = "DEBUG"
	RestProxyKarapaceHostName          = "host"
	RestProxyKarapacePortName          = "port"
	RestProxyKafkaBootstrapURI         = "bootstrap_uri"
	RestProxyKafkaSecurityProtocolName = "security_protocol"
	RestProxyKafkaSASLMechanismName    = "sasl_mechanism"
	RestProxyKafkaSASLUsername         = "sasl_plain_username"
	RestProxyKafkaSASLPassword         = "sasl_plain_password"
	RestProxyKafkaSSLCAFile            = "ssl_cafile"
	RestProxyKafkaSSLCertFile          = "ssl_certfile"
	RestProxyKafkaSSLKeyFile           = "ssl_keyfile"
	RestProxyKafkaSSLCAFilePath        = "/var/private/kafka-client-ssl/ca.crt"
	RestProxyKafkaSSLCertFilePath      = "/var/private/kafka-client-ssl/tls.crt"
	RestProxyKafkaSSLKeyFilePath       = "/var/private/kafka-client-ssl/tls.key"
	RestProxyConfigFileName            = "rest.config.json"
)
