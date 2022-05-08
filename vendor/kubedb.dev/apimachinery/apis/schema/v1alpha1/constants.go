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
	DatabaseSchemaConditionTypeDBServerReady  DatabaseSchemaConditionType = "DatabaseServerReady"
	DatabaseSchemaMessageDBServerNotCreated   DatabaseSchemaMessage       = "Database Server is not created yet"
	DatabaseSchemaMessageDBServerProvisioning DatabaseSchemaMessage       = "Database Server is provisioning"
	DatabaseSchemaMessageDBServerReady        DatabaseSchemaMessage       = "Database Server is Ready"

	DatabaseSchemaConditionTypeVaultReady  DatabaseSchemaConditionType = "VaultReady"
	DatabaseSchemaMessageVaultNotCreated   DatabaseSchemaMessage       = "VaultServer is not created yet"
	DatabaseSchemaMessageVaultProvisioning DatabaseSchemaMessage       = "VaultServer is provisioning"
	DatabaseSchemaMessageVaultReady        DatabaseSchemaMessage       = "VaultServer is Ready"

	DatabaseSchemaConditionTypeDoubleOptInNotPossible DatabaseSchemaConditionType = "DoubleOptInNotPossible"
	DatabaseSchemaMessageDoubleOptInNotPossible       DatabaseSchemaMessage       = "Double OptIn is not possible between the applied Schema & Database server"

	DatabaseSchemaConditionTypeSecretEngineReady DatabaseSchemaConditionType = "SecretEngineReady"
	DatabaseSchemaMessageSecretEngineNotCreated  DatabaseSchemaMessage       = "SecretEngine is not created yet"
	DatabaseSchemaMessageSecretEngineCreating    DatabaseSchemaMessage       = "SecretEngine is being creating"
	DatabaseSchemaMessageSecretEngineSuccess     DatabaseSchemaMessage       = "SecretEngine phase is success"

	DatabaseSchemaConditionTypeRoleReady        DatabaseSchemaConditionType = "RoleReady"
	DatabaseSchemaMessageDatabaseRoleNotCreated DatabaseSchemaMessage       = "Database Role is not created yet"
	DatabaseSchemaMessageDatabaseRoleCreating   DatabaseSchemaMessage       = "Database Role is being creating"
	DatabaseSchemaMessageDatabaseRoleSuccess    DatabaseSchemaMessage       = "Database Role is success"

	DatabaseSchemaConditionTypeSecretAccessRequestReady DatabaseSchemaConditionType = "SecretAccessRequestReady"
	DatabaseSchemaMessageSecretAccessRequestNotCreated  DatabaseSchemaMessage       = "SecretAccessRequest is not created yet"
	DatabaseSchemaMessageSecretAccessRequestWaiting     DatabaseSchemaMessage       = "SecretAccessRequest is waiting for approval"
	DatabaseSchemaMessageSecretAccessRequestApproved    DatabaseSchemaMessage       = "SecretAccessRequest has been approved"
	DatabaseSchemaMessageSecretAccessRequestExpired     DatabaseSchemaMessage       = "SecretAccessRequest has been expired"

	DatabaseSchemaConditionTypeDBCreationUnsuccessful DatabaseSchemaConditionType = "DatabaseCreationUnsuccessful"
	DatabaseSchemaMessageSchemaNameConflicted         DatabaseSchemaMessage       = "Schema name is conflicted"
	DatabaseSchemaMessageDBCreationUnsuccessful       DatabaseSchemaMessage       = "Internal error occurred when creating database"

	DatabaseSchemaConditionTypeInitScriptCompleted DatabaseSchemaConditionType = "InitScriptCompleted"
	DatabaseSchemaMessageInitScriptNotApplied      DatabaseSchemaMessage       = "InitScript is not applied yet"
	DatabaseSchemaMessageInitScriptRunning         DatabaseSchemaMessage       = "InitScript is running"
	DatabaseSchemaMessageInitScriptCompleted       DatabaseSchemaMessage       = "InitScript is completed"
	DatabaseSchemaMessageInitScriptSucceeded       DatabaseSchemaMessage       = "InitScript is succeeded"
	DatabaseSchemaMessageInitScriptFailed          DatabaseSchemaMessage       = "InitScript is failed"

	DatabaseSchemaConditionTypeRepositoryFound DatabaseSchemaConditionType = "RepositoryFound"
	DatabaseSchemaMessageRepositoryNotCreated  DatabaseSchemaMessage       = "Repository is not created yet"
	DatabaseSchemaMessageRepositoryFound       DatabaseSchemaMessage       = "Repository has been found"

	DatabaseSchemaConditionTypeAppBindingFound DatabaseSchemaConditionType = "AppBindingFound"
	DatabaseSchemaMessageAppBindingNotCreated  DatabaseSchemaMessage       = "AppBinding is not created yet"
	DatabaseSchemaMessageAppBindingFound       DatabaseSchemaMessage       = "AppBinding is Found"

	DatabaseSchemaConditionTypeRestoreCompleted   DatabaseSchemaConditionType = "RestoreSessionCompleted"
	DatabaseSchemaMessageRestoreSessionNotCreated DatabaseSchemaMessage       = "RestoreSession is not created yet"
	DatabaseSchemaMessageRestoreSessionRunning    DatabaseSchemaMessage       = "RestoreSession is running"
	DatabaseSchemaMessageRestoreSessionSucceed    DatabaseSchemaMessage       = "RestoreSession is succeeded"
	DatabaseSchemaMessageRestoreSessionFailed     DatabaseSchemaMessage       = "RestoreSession is failed"
)

const (
	MySQLEncryptionEnabled  string = "'Y'"
	MySQLEncryptionDisabled string = "'N'"
)
