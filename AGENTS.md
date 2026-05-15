# AGENTS.md - KubeDB CLI

This file provides instructions for AI coding agents working in the KubeDB `kubectl` plugin Go repository.

## Project Overview

`kubectl-dba` is the official `kubectl` plugin for KubeDB (Kubernetes database operator platform by AppsCode). It provides DBA-style commands for KubeDB-managed databases: describe, connect/exec, restart, pause/resume, debug, monitor, show credentials, generate appbinding/secrets for remote replicas, and insert/verify test data. Distributed as a single static binary via GitHub Releases and as the `dba` plugin in the krew-index.

Supported databases include: Cassandra, ClickHouse, Druid, Elasticsearch, FerretDB, Hazelcast, Ignite, Kafka, MariaDB, Memcached, MongoDB, MSSQLServer, MySQL, Oracle, PerconaXtraDB, PgBouncer, Pgpool, PostgreSQL, ProxySQL, RabbitMQ, Redis, SingleStore, Solr, ZooKeeper.

## Build & Development Commands

Builds run inside the `ghcr.io/appscode/golang-dev:1.25` Docker image (set by `BUILD_IMAGE`). The output binary is `bin/kubectl-dba-$(GOOS)-$(GOARCH)`.

```bash
# Default: format + build for host OS/ARCH
make

# Build only (host platform)
make build

# Cross-build a specific platform
make build GOOS=linux GOARCH=arm64

# Build all release platforms (linux/amd64, linux/arm, linux/arm64, windows/amd64, darwin/amd64, darwin/arm64)
make all-build

# Compressed release artifacts (.tar.gz / .zip + checksums)
make all-build COMPRESS=yes

# Format Go sources (runs hack/fmt.sh inside docker)
make fmt

# Run unit tests
make test
make unit-tests

# Lint (golangci-lint inside docker)
make lint

# Full CI pipeline: verify check-license lint build unit-tests
make ci

# License header management
make add-license
make check-license

# Verify go.mod / vendor / generated files are current
make verify         # verify-gen + verify-modules

# Print version / build metadata
make version

# Generate krew plugin manifest (requires a git tag)
make gen-krew-manifest

# Remove build artifacts
make clean
```

Release flow: `make qa` (non-prod, branch builds) and `make release` (only with `APPSCODE_ENV=prod` and a git tag). The `.github/workflows/release.yml` job runs `make release COMPRESS=yes`, uploads tarballs to GitHub Releases, then opens a PR to `appscode/krew-index` via `hack/krew/plugin.yaml`.

The `hack/build.sh` script injects ldflags (`main.Version`, `main.GitTag`, `main.CommitHash`, etc.) read by `cmd/kubectl-dba/version.go`.

## Project Structure

```
cmd/kubectl-dba/
  main.go               # Entry point: calls cmds.NewKubeDBCommand
  version.go            # Build-time ldflag variables -> gomodules.xyz/x/version

pkg/
  cmds/                 # Cobra command wiring (one file per top-level command)
    root.go             # NewKubeDBCommand builds the command tree and groups
    describe.go connect.go exec.go data.go debug.go monitor.go
    pause.go resume.go restart.go remote_replica.go mssql.go
    show_credentials.go completion.go options.go
  common/               # Shared SQL helpers (mssql, mysql, postgres)
  connect/              # `kubectl dba connect` per-database implementations
  credentials/          # `kubectl dba show-credentials` per-database implementations
  data/                 # `kubectl dba data insert/drop/verify` test-data clients
    redisutil/          # Redis cluster helpers
  debug/                # `kubectl dba debug` per-database collectors + gitops.go
  describer/            # `kubectl dba describe` printers (kubectl describe-style)
  events/               # Event sorting helpers
  lib/                  # Shared client helpers
  monitor/              # `kubectl dba monitor` (Prometheus, alerts, dashboards)
    alerts/ connection/ dashboard/
  pauser/ resumer/      # Per-database pause/resume logic + archiver
  printer/              # Table/YAML printers
  remote_replica/       # AppBinding + Secret generation for MySQL/Postgres replicas
  restarter/            # Per-database restart logic

hack/
  build.sh test.sh fmt.sh coverage.sh e2e.sh   # Wrapped by Makefile
  krew/plugin.yaml      # krew manifest template (populated by gen-krew-manifest)
  gendocs/              # Tool to regenerate kubectl-dba reference docs
  license/              # License header template (used by ltag)
  scripts/

.github/workflows/
  ci.yml                # PR/push: runs `make ci`
  release.yml           # Tag push: cross-build, GitHub Release, krew PR
  release-tracker.yml update-docs.yml

vendor/                 # Vendored deps (built with GOFLAGS=-mod=vendor)
```

## Key Packages / APIs

- `pkg/cmds.NewKubeDBCommand(in, out, err) *cobra.Command` (`pkg/cmds/root.go`) - assembles the root Cobra command, kubeconfig flags via `genericclioptions.NewConfigFlags`, and grouped subcommands using `k8s.io/kubectl/pkg/util/templates`.
- `cmd/kubectl-dba/main.go` - thin `main` that constructs the root command, calls `logs.Init`, and executes.
- Per-command constructors follow the convention `NewCmd<Verb>(parent string, f cmdutil.Factory, streams genericiooptions.IOStreams) *cobra.Command` (see `describe.go`, `restart.go`, `pause.go`, `resume.go`, `show_credentials.go`). Connection-style commands use the shorter `NewCmd<Verb>(f cmdutil.Factory)` form (`connect.go`, `exec.go`, `data.go`, `debug.go`, `monitor.go`, `remote_replica.go`, `mssql.go`).
- Command groups defined in `root.go`: Troubleshooting/Debugging, Database Ops, Pause/Resume, Database Connection, Data Insert/Verify, Debug, Remote Replica appbinding, MSSQLServer, Metrics.
- Per-database engine packages (`pkg/connect`, `pkg/data`, `pkg/debug`, `pkg/pauser`, `pkg/resumer`, `pkg/restarter`, `pkg/credentials`, `pkg/describer`) each follow the same one-file-per-database layout (`mysql.go`, `postgres.go`, `mongodb.go`, `redis.go`, ...) plus a small `const.go` / `helpers.go` / `common.go` and a top-level dispatcher (e.g. `pauser.go`, `resumer.go`, `restarter.go`, `show_cred.go`).

## Testing

- Unit tests run via `make test` (alias for `make unit-tests`), which invokes `hack/test.sh cmd pkg` inside the build container. The script just runs `go test -installsuffix "static" ./cmd/... ./pkg/...` with `GOFLAGS=-mod=vendor`.
- `hack/coverage.sh` generates coverage output; `hack/e2e.sh` is reserved for e2e runs.
- CI (`.github/workflows/ci.yml`) runs `make ci` (`verify check-license lint build unit-tests`) on Go 1.25 / Ubuntu 24.04 for every PR and push to `master`.
- There are currently no `_test.go` files outside vendor; most validation is via build + lint + manual/cluster testing.

## Dependencies

Module path: `kubedb.dev/cli`, Go `1.25.5`.

Internal (KubeDB / AppsCode ecosystem):

- `kubedb.dev/apimachinery` - KubeDB CRD types (Database, Ops, Archiver, Schema, etc.)
- `kubedb.dev/db-client-go` - shared per-engine client helpers used by `connect`/`data`/`debug`
- `kmodules.xyz/client-go`, `kmodules.xyz/custom-resources`, `kmodules.xyz/cert-manager-util`, `kmodules.xyz/monitoring-agent-api`
- `kubeops.dev/petset` - PetSet (StatefulSet replacement) types
- `stash.appscode.dev/apimachinery` - Stash backup CRDs
- `gomodules.xyz/x` (version, logs), `gomodules.xyz/go-sh`, `gomodules.xyz/pointer`, `gomodules.xyz/runtime`

External:

- `github.com/spf13/cobra` - CLI framework
- `k8s.io/{api,apimachinery,cli-runtime,client-go,component-base,kubectl}` v0.34.x - Kubernetes client and `kubectl` plugin scaffolding
- `sigs.k8s.io/controller-runtime` v0.22.4 (replaced with the `kmodules/controller-runtime` fork)
- `github.com/cert-manager/cert-manager` v1.19.4
- Database drivers / clients: `github.com/go-sql-driver/mysql`, `github.com/redis/go-redis/v9`, `github.com/elastic/go-elasticsearch/v{5..9}`, `github.com/opensearch-project/opensearch-go/v{1..3}` (all indirect via `db-client-go`)
- `github.com/prometheus/{client_golang,common}` - metrics integration

`go.mod` `replace` directives (do not casually edit):

```
github.com/Masterminds/sprig/v3 => github.com/gomodules/sprig/v3 v3.2.3-0.20220405051441-0a8a99bac1b8
sigs.k8s.io/controller-runtime  => github.com/kmodules/controller-runtime v0.22.5-0.20251227114913-f011264689cd
github.com/imdario/mergo        => github.com/imdario/mergo v0.3.6
k8s.io/apiserver                => github.com/kmodules/apiserver v0.34.4-0.20251227112449-07fa35efc6fc
```

Builds use vendored dependencies (`GOFLAGS=-mod=vendor`). After bumping any module, run `make verify-modules` to refresh `go.mod`/`go.sum`/`vendor/`.

## Code Conventions

- License headers (`hack/license/`) are required on every Go file; enforce with `make add-license` / `make check-license` (uses `ltag`).
- `golangci.yml` enables the standard linter set plus `unparam`, with `gofmt` configured to rewrite `interface{}` to `any` and `goimports` enabled. Generated files (`generated.*\.go`) and `client/`, `vendor/` are excluded.
- Use `cmdutil.Factory` and `genericiooptions.IOStreams` from `k8s.io/cli-runtime` for all new subcommands - never call `kubectl` config helpers directly.
- Per-database parallel structure: when adding support for a new engine, add a sibling file in each relevant `pkg/<feature>/<engine>.go` and wire it into the dispatcher (`pkg/<feature>/<feature>.go` or the corresponding `pkg/cmds/<verb>.go`).
- The plugin binary is named `kubectl-dba` (krew name: `dba`); preserve this when editing `Makefile`, `hack/krew/plugin.yaml`, or release workflow asset names.
- Do not commit `bin/`, `.go/`, or `vendor/` changes that are not the result of a deliberate dependency update.
