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

// Backup and restore names, flags, and defaults shared with the Neo4j operator.
const (
	Neo4jDataSubPath               = "data"
	SnapshotIDAlphabet             = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	DataBackupRestorer             = "data-backup-restorer"
	DifferentialSnapshotSuffix     = "differential-snapshot"
	DifferentialBackupCommand      = "differential-backup"
	NamespaceFlag                  = "namespace"
	DifferentialSnapshotFlag       = "diff-snapshot-name"
	DifferentialBackupConfigFlag   = "diff-backupconfig-name"
	SuccessfulLogHistoryLimitFlag  = "successful-log-history-limit"
	FailedLogHistoryLimitFlag      = "failed-log-history-limit"
	DifferentialBackupIntervalFlag = "differential-backup-interval"
	RestoreUntilFlag               = "restore-until"
	SeedServerNameParam            = "seedServerName"
	Neo4jAdminArgsParam            = "neo4jAdminArgs"
	RestoreTimeLayout              = "2006-01-02 15:04:05"
	DefaultLogHistoryLimit         = 5

	BackupDirNeo4jMetadata     = "metadata"
	BackupDirNeo4jData         = "data"
	BackupDirNeo4jDifferential = "differential"

	Neo4jDriver = "Neo4jAdmin"
)

const (
	RepositorySuffixFull     = "full"
	RepositorySuffixManifest = "manifest"
	SessionNameFull          = "full-backup"
	SessionNameManifest      = "manifest-backup"

	BackupConfigNameSuffix = "archiver"
	SnapshotNameSuffix     = "incremental-snapshot"
	SidekickNameSuffix     = "sidekick"
	WalgContainerName      = "wal-g"
	RestoreSessionName     = "manifest-restorer"

	RestoreJobNameBinlog  = "binlog-restorer"
	RestoreJobNameLog     = "log-restorer"
	RestoreJobNameOplog   = "oplog-restorer"
	RestoreJobNameWal     = "wal-restorer"
	RestoreCmdBinlogFetch = "binlog-fetch"
	RestoreCmdOplogReplay = "oplog-replay"
	RestoreCmdWalFetch    = "wal-fetch"

	BackupDirBinlog     = "binlog-backup"
	BackupDirOplog      = "oplog-backup"
	BackupDirWal        = "wal-backup"
	BackupCmdBinlogPush = "binlog-push"
	BackupCmdOplogPush  = "oplog-push"
	BackupCmdWalPush    = "wal-push"
)

// azure
const (
	WALG_AZ_PREFIX           = "WALG_AZ_PREFIX"
	AZURE_STORAGE_ACCOUNT    = "AZURE_STORAGE_ACCOUNT"
	AZURE_STORAGE_ACCESS_KEY = "AZURE_STORAGE_ACCESS_KEY"
	AZURE_STORAGE_KEY        = "AZURE_STORAGE_KEY"
	AZURE_ACCOUNT_KEY        = "AZURE_ACCOUNT_KEY"
)

// s3
const (
	WALG_S3_PREFIX          = "WALG_S3_PREFIX"
	AWS_ENDPOINT            = "AWS_ENDPOINT"
	AWS_REGION              = "AWS_REGION"
	AWS_S3_FORCE_PATH_STYLE = "AWS_S3_FORCE_PATH_STYLE"
	WALG_S3_CA_CERT_FILE    = "WALG_S3_CA_CERT_FILE"
	CA_CERT_DATA            = "CA_CERT_DATA"

	S3CredVolumeName = "s3-cred"
	S3CAMountPath    = "/s3-cred/public.crt"
)

// gcs
const (
	WALG_GS_PREFIX                 = "WALG_GS_PREFIX"
	GOOGLE_APPLICATION_CREDENTIALS = "GOOGLE_APPLICATION_CREDENTIALS"

	GoogleCredVolumeName = "google-cred"
	GoogleCredMountPath  = "/google-cred"
	GoogleCredFileName   = "GOOGLE_SERVICE_ACCOUNT_JSON_KEY"
)

// others
const (
	WALG_FILE_PREFIX                   = "WALG_FILE_PREFIX"
	OPLOG_PUSH_WAIT_FOR_BECOME_PRIMARY = "OPLOG_PUSH_WAIT_FOR_BECOME_PRIMARY"
)

// Sidekick's ExtraArgs keys
const (
	ExtraArgsKeyStorage      string = "Storage"
	ExtraArgsKeyLogBackupOpt string = "LogBackupOpt"
)

const (
	// TypeInitialBackupSucceeded indicates whether the initial backup session has succeeded or not.
	TypeInitialBackupSucceeded                      = "InitialBackupSucceeded"
	ReasonSuccessfullyCompletedInitialBackupSession = "SuccessfullyCompletedInitialBackup"
	ReasonFailedToCompleteInitialBackupSession      = "FailedToCompleteInitialBackup"
)

const (
	BackupDirClickHouseMetadata    = "metadata"
	BackupDirClickHouseData        = "data"
	BackupDirClickHouseIncremental = "incremental"

	ComponentNameMetadata = "metadata"
	ComponentNameData     = "data"

	ClickHouseDriver = "ClickHouseBackup"
)
