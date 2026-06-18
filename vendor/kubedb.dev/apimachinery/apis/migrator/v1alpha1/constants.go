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
	KindClusterRole           = "ClusterRole"
	KindRole                  = "Role"
	KindPersistentVolumeClaim = "PersistentVolumeClaim"
	KindJob                   = "Job"
)

const (
	MigratorJobClusterRole = "migrator-job"
	MigratorJobPrefix      = "migrator"
	SidecarContainerName   = "status-reporter"
	MigratorGRPCPort       = 50051
	MigratorPVCSuffix      = "pvc"
	PVCVolumeName          = "migrator-data"
	PVCVolumeMountPath     = "/data"
	ConfigVolName          = "migrator-config"
	ConfigVolMountPath     = "/etc/migrator"
	ConfigFileName         = "config.yaml"
	MigratorConfigSuffix   = "config"
	ConfigPath             = ConfigVolMountPath + "/" + ConfigFileName
	TLSVolumePrefix        = "tls"
	TLSMountPathPrefix     = "/etc/tls"
	SourceDatabaseRole     = "source"
	TargetDatabaseRole     = "target"
)

// Conditions Related Constants
const (
	MigratorJobTriggered = "MigratorJobTriggered"

	DestroySignalSend = "DestroySignalSend"

	// MigrationRunning Migration status conditions
	MigrationRunning       = "MigrationRunning"
	ReasonMigrationRunning = "MigrationInProgress"

	MigrationSucceeded       = "MigrationSucceeded"
	ReasonMigrationSucceeded = "MigrationCompleted"

	MigrationFailed       = "MigrationFailed"
	ReasonMigrationFailed = "MigrationError"
)

// ============ CLI Constants ==================
const (
	MySQLDump   = "mysqldump"
	MySQLCli    = "mysql"
	MariaDBCli  = "mariadb"
	MariaDBDump = "mariadb-dump"
)

// ============ Snapshot Pipeline Constant ==================
const (
	SnapshotWorker     = 3
	SnapshotSinker     = 3
	SnapshotBuffer     = 10
	SnapshotReadBatch  = 5000
	SnapshotWriteBatch = 500
)

// ============ SQLite Related Constants ==================
const (
	// SQLite Constants
	sqliteMaxOpenConns = 35
	sqliteMaxIdleConns = 10
	sqliteBusyTimeout  = 5000
	SqliteFile         = "migration.db"

	// Status Constants
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
	StatusSkipped    = "skipped"

	// Phase Constants
	PhaseSchema    = "schema"
	PhaseSnapshot  = "snapshot"
	PhaseStreaming = "streaming"

	// Table Names
	TableMigrationInfo = "migration_info"
	TablePhases        = "phases"
	TableSnapshot      = "snapshot"
	TableStreaming     = "streaming"

	// Position methods
	PositionMethodGTID    = "gtid"
	PositionMethodFilePos = "file_pos"

	// migration_info columns
	ColMigrationID      = "migration_id"
	ColCreatedAt        = "created_at"
	ColRestartedAt      = "restarted_at"
	ColRestartCount     = "restart_count"
	ColLastErrorMessage = "last_error_message"

	// phases columns
	ColPhaseID          = "id"
	ColPhaseName        = "name"
	ColPhaseEnabled     = "enabled"
	ColPhaseStatus      = "status"
	ColPhaseStartedAt   = "started_at"
	ColPhaseCompletedAt = "completed_at"
	ColPhaseErrorMsg    = "error_message"

	// snapshot columns
	ColTableID           = "table_id"
	ColDatabaseName      = "database_name"
	ColTableName         = "table_name"
	ColStatus            = "status"
	ColRowsCopied        = "rows_copied"
	ColLastInsertedBatch = "last_inserted_batch"
	ColLastInsertedKey   = "last_inserted_key"
	ColStartedAt         = "started_at"
	ColCompletedAt       = "completed_at"
	ColErrorMessage      = "error_message"

	// streaming columns
	ColPositionMethod = "position_method"
	ColGTIDSet        = "gtid_set"
	ColBinlogFile     = "binlog_file"
	ColBinlogPos      = "binlog_pos"
	ColUpdatedAt      = "updated_at"
)
