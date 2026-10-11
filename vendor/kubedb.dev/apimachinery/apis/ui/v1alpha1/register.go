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

// Package v1alpha1 contains API Schema definitions for the kubedb v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=ui.kubedb.com
package v1alpha1

import (
	"kubedb.dev/apimachinery/apis/ui"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var SchemeGroupVersion = schema.GroupVersion{Group: ui.GroupName, Version: "v1alpha1"}

var (
	// TODO: move SchemeBuilder with zz_generated.deepcopy.go to k8s.io/api.
	// localSchemeBuilder and AddToScheme will stay in k8s.io/kubernetes.
	SchemeBuilder      runtime.SchemeBuilder
	localSchemeBuilder = &SchemeBuilder
	AddToScheme        = localSchemeBuilder.AddToScheme
)

func init() {
	// We only register manually written functions here. The registration of the
	// generated functions takes place in the generated files. The separation
	// makes the code compile even when the generated files are missing.
	localSchemeBuilder.Register(addKnownTypes)
}

// Kind takes an unqualified kind and returns a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// Adds the list of known types to api.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(
		SchemeGroupVersion,
		&DatabaseConnection{},
		&DatabaseConnectionList{},
		&DatabaseConfiguration{},
		&DatabaseConfigurationList{},
		&DatabaseConfigurationOptions{},
		&DatabaseSummary{},
		&CassandraInsight{},
		&CassandraInsightList{},
		&CassandraSchemaOverview{},
		&CassandraSchemaOverviewList{},
		&DB2Insight{},
		&DB2InsightList{},
		&DB2Queries{},
		&DB2QueriesList{},
		&DB2SchemaOverview{},
		&DB2SchemaOverviewList{},
		&DruidInsight{},
		&DruidInsightList{},
		&DruidTasks{},
		&DruidTasksList{},
		&DocumentDBInsight{},
		&DocumentDBInsightList{},
		&DocumentDBSchemaOverview{},
		&DocumentDBSchemaOverviewList{},
		&DruidSchemaOverview{},
		&DruidSchemaOverviewList{},
		&ElasticsearchInsight{},
		&ElasticsearchInsightList{},
		&ElasticsearchNodesStats{},
		&ElasticsearchNodesStatsList{},
		&ElasticsearchSchemaOverview{},
		&ElasticsearchSchemaOverviewList{},
		&MariaDBInsight{},
		&MariaDBInsightList{},
		&MariaDBQueries{},
		&MariaDBQueriesList{},
		&MariaDBSchemaOverview{},
		&MariaDBSchemaOverviewList{},
		&MilvusInsight{},
		&MilvusInsightList{},
		&MilvusSchemaOverview{},
		&MilvusSchemaOverviewList{},
		&MongoDBInsight{},
		&MongoDBInsightList{},
		&MongoDBQueries{},
		&MongoDBQueriesList{},
		&MongoDBSchemaOverview{},
		&MongoDBSchemaOverviewList{},
		&Neo4jInsight{},
		&Neo4jInsightList{},
		&Neo4jActivity{},
		&Neo4jActivityList{},
		&Neo4jSchemaOverview{},
		&Neo4jSchemaOverviewList{},
		&MySQLInsight{},
		&MySQLInsightList{},
		&MySQLQueries{},
		&MySQLQueriesList{},
		&MySQLSchemaOverview{},
		&MySQLSchemaOverviewList{},
		&PerconaXtraDBInsight{},
		&PerconaXtraDBInsightList{},
		&PerconaXtraDBQueries{},
		&PerconaXtraDBQueriesList{},
		&PerconaXtraDBSchemaOverview{},
		&PerconaXtraDBSchemaOverviewList{},
		&PgBouncerInsight{},
		&PgBouncerInsightList{},
		&PgBouncerPoolOverview{},
		&PgBouncerPoolOverviewList{},
		&PgBouncerServerOverview{},
		&PgBouncerServerOverviewList{},
		&PgBouncerSettings{},
		&PgBouncerSettingsList{},
		&PostgresInsight{},
		&PostgresInsightList{},
		&PostgresQueries{},
		&PostgresQueriesList{},
		&PostgresSchemaOverview{},
		&PostgresSchemaOverviewList{},
		&PostgresSettings{},
		&PostgresSettingsList{},
		&ProxySQLInsight{},
		&ProxySQLInsightList{},
		&ProxySQLQueries{},
		&ProxySQLQueriesList{},
		&ProxySQLSettings{},
		&ProxySQLSettingsList{},
		&RedisInsight{},
		&RedisInsightList{},
		&RedisQueries{},
		&RedisQueriesList{},
		&RedisSchemaOverview{},
		&RedisSchemaOverviewList{},
		&SinglestoreInsight{},
		&SinglestoreInsightList{},
		&SinglestoreQueries{},
		&SinglestoreQueriesList{},
		&SinglestoreSchemaOverview{},
		&SinglestoreSchemaOverviewList{},
	)

	scheme.AddKnownTypes(
		SchemeGroupVersion,
		&metav1.Status{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
