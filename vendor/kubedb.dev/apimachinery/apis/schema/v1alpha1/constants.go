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

import kmapi "kmodules.xyz/client-go/api/v1"

const (
	DatabaseNameAdmin  = "admin"
	DatabaseNameConfig = "config"
	DatabaseNameLocal  = "local"
	DatabaseForEntry   = "kubedb_system"
	SYSDatabase        = "sys"
)

const (
	// DeletionPolicyDelete allows the created objects to be deleted
	DeletionPolicyDelete DeletionPolicy = "Delete"
	// DeletionPolicyDoNotDelete Rejects attempt to delete using ValidationWebhook.
	DeletionPolicyDoNotDelete DeletionPolicy = "DoNotDelete"
)

const (
	DatabaseSchemaPhasePending     DatabaseSchemaPhase = "Pending"
	DatabaseSchemaPhaseInProgress  DatabaseSchemaPhase = "InProgress"
	DatabaseSchemaPhaseTerminating DatabaseSchemaPhase = "Terminating"
	DatabaseSchemaPhaseCurrent     DatabaseSchemaPhase = "Current"
	DatabaseSchemaPhaseFailed      DatabaseSchemaPhase = "Failed"
	DatabaseSchemaPhaseExpired     DatabaseSchemaPhase = "Expired"
)

const (
	DatabaseSchemaConditionTypeDBServerReady  kmapi.ConditionType = "DatabaseServerReady"
	DatabaseSchemaMessageDBServerNotCreated   string              = "Database Server is not created yet"
	DatabaseSchemaMessageDBServerProvisioning string              = "Database Server is provisioning"
	DatabaseSchemaMessageDBServerReady        string              = "Database Server is Ready"

	DatabaseSchemaConditionTypeVaultReady  kmapi.ConditionType = "VaultReady"
	DatabaseSchemaMessageVaultNotCreated   string              = "VaultServer is not created yet"
	DatabaseSchemaMessageVaultProvisioning string              = "VaultServer is provisioning"
	DatabaseSchemaMessageVaultReady        string              = "VaultServer is Ready"

	DatabaseSchemaConditionTypeDoubleOptInNotPossible kmapi.ConditionType = "DoubleOptInNotPossible"
	DatabaseSchemaMessageDoubleOptInNotPossible       string              = "Double OptIn is not possible between the applied Schema & Database server"

	DatabaseSchemaConditionTypeSecretEngineReady kmapi.ConditionType = "SecretEngineReady"
	DatabaseSchemaMessageSecretEngineNotCreated  string              = "SecretEngine is not created yet"
	DatabaseSchemaMessageSecretEngineCreating    string              = "SecretEngine is being creating"
	DatabaseSchemaMessageSecretEngineSuccess     string              = "SecretEngine phase is success"

	DatabaseSchemaConditionTypeRoleReady        kmapi.ConditionType = "RoleReady"
	DatabaseSchemaMessageDatabaseRoleNotCreated string              = "Database Role is not created yet"
	DatabaseSchemaMessageDatabaseRoleCreating   string              = "Database Role is being creating"
	DatabaseSchemaMessageDatabaseRoleSuccess    string              = "Database Role is success"

	DatabaseSchemaConditionTypeSecretAccessRequestReady kmapi.ConditionType = "SecretAccessRequestReady"
	DatabaseSchemaMessageSecretAccessRequestNotCreated  string              = "SecretAccessRequest is not created yet"
	DatabaseSchemaMessageSecretAccessRequestWaiting     string              = "SecretAccessRequest is waiting for approval"
	DatabaseSchemaMessageSecretAccessRequestApproved    string              = "SecretAccessRequest has been approved"
	DatabaseSchemaMessageSecretAccessRequestExpired     string              = "SecretAccessRequest has been expired"

	DatabaseSchemaConditionTypeDBCreationUnsuccessful kmapi.ConditionType = "DatabaseCreationUnsuccessful"
	DatabaseSchemaMessageSchemaNameConflicted         string              = "Schema name is conflicted"
	DatabaseSchemaMessageDBCreationUnsuccessful       string              = "Internal error occurred when creating database"

	DatabaseSchemaConditionTypeInitScriptCompleted kmapi.ConditionType = "InitScriptCompleted"
	DatabaseSchemaMessageInitScriptNotApplied      string              = "InitScript is not applied yet"
	DatabaseSchemaMessageInitScriptRunning         string              = "InitScript is running"
	DatabaseSchemaMessageInitScriptCompleted       string              = "InitScript is completed"
	DatabaseSchemaMessageInitScriptSucceeded       string              = "InitScript is succeeded"
	DatabaseSchemaMessageInitScriptFailed          string              = "InitScript is failed"

	DatabaseSchemaConditionTypeRepositoryFound kmapi.ConditionType = "RepositoryFound"
	DatabaseSchemaMessageRepositoryNotCreated  string              = "Repository is not created yet"
	DatabaseSchemaMessageRepositoryFound       string              = "Repository has been found"

	DatabaseSchemaConditionTypeAppBindingFound kmapi.ConditionType = "AppBindingFound"
	DatabaseSchemaMessageAppBindingNotCreated  string              = "AppBinding is not created yet"
	DatabaseSchemaMessageAppBindingFound       string              = "AppBinding is Found"

	DatabaseSchemaConditionTypeRestoreCompleted   kmapi.ConditionType = "RestoreSessionCompleted"
	DatabaseSchemaMessageRestoreSessionNotCreated string              = "RestoreSession is not created yet"
	DatabaseSchemaMessageRestoreSessionRunning    string              = "RestoreSession is running"
	DatabaseSchemaMessageRestoreSessionSucceed    string              = "RestoreSession is succeeded"
	DatabaseSchemaMessageRestoreSessionFailed     string              = "RestoreSession is failed"
)

const (
	MySQLEncryptionEnabled  string = "'Y'"
	MySQLEncryptionDisabled string = "'N'"
)
