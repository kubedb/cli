# `kubectl dba dc-dr`: cross data center DR commands

Day-2 operations for a KubeDB database distributed across data centers: trigger and
monitor planned switchovers, release a failover held by the RPO budget, move the
failover authority, pin a data center, and diagnose a failover that is not happening.

Every command below was live-verified against a three data center cluster (two Member
DCs plus an Arbiter) before being documented here.

## Which cluster each command talks to

DC-DR spans three different API surfaces, and every command documents which one it
needs. Getting this wrong is the most common mistake, so it is stated per command in
`--help` as well.

| surface | what lives there | how the CLI reaches it |
|---|---|---|
| **hub** | the KubeDB operator, the Postgres CRs, the PlacementPolicies | the ordinary kubeconfig flags: `--kubeconfig`, `$KUBECONFIG`, `--context`, `-n` |
| **coordination control plane** | the `primary-dc*` Leases (the failover authority) | the `--coord-*` flags below |
| **spoke** (one data center) | that DC's agent, its marker ConfigMap, and the human-owned pin ConfigMaps | the ordinary kubeconfig flags, pointed AT that spoke |

### Reaching the coordination control plane

It is a separate apiserver, so it needs its own kubeconfig. Commands that read or
write Leases (`active-dc`, `handoff`, `debug failover`, `debug fence`) accept, in
priority order:

```
--coord-kubeconfig /path/to/coord.yaml              # a file
--coord-kubeconfig-configmap [namespace/]name       # key "kubeconfig" in a ConfigMap on the CURRENT cluster
--coord-kubeconfig-secret [namespace/]name          # key "kubeconfig" in a Secret on the CURRENT cluster
                                                    # default: dc-failover/coord-kubeconfig
--coord-namespace dc-failover                       # namespace holding the Leases
```

The default works out of the box on a standard install: the dr-controlplane chart
mints `dc-failover/coord-kubeconfig` on the hub, so pointing `$KUBECONFIG` at the hub
is usually all that is needed.

## Command summary

| command | what it does | talks to |
|---|---|---|
| [`switchover`](#switchover) | planned, zero-RPO move of the primary role | hub |
| [`status`](#status) | one-shot state and switchover progress | hub |
| [`abort`](#abort) | abort an in-flight switchover | hub |
| [`accept-data-loss`](#accept-data-loss) | release a failover held by the RPO budget | hub |
| [`handoff`](#handoff) | move the failover authority (the failover lever) | hub + coordination |
| [`pin-primary`](#pin-primary) | break-glass: this DC stays primary, no failover | spoke |
| [`pin-standby`](#pin-standby) | this DC never promotes | spoke |
| [`active-dc`](#active-dc) | who holds the primary role right now | hub + coordination |
| [`debug`](#debug) | diagnose failover / switchover / fence symptoms | hub + coordination |

---

## switchover

Planned, zero-RPO move of the primary role to another data center.

```sh
kubectl dba dc-dr switchover DB_NAME -n NS --to DC
```

Sets `dr.kubedb.com/switchover-to`. The operator quiesces the active primary
(write-locks it so its LSN freezes), waits for the target to replay to that exact
LSN, hands off the Lease, and clears the annotation. No committed row is lost.

**Requires the active primary to be up and reachable.** All three of the operator's
gates measure by dialing it and fail closed, so a dead primary can never be switched
away from; use [`handoff`](#handoff) for that case.

Refuses up front: a target that is not a Member DC (an Arbiter or Witness can never
hold the primary role), and a target that is already active (prints a no-op notice).

```
$ kubectl dba dc-dr switchover cli-test -n demo --to dc-a
Switchover of demo/cli-test to "dc-a" requested (scope primary-dc-clitest, from PlacementPolicy cli-test-pp failoverPolicy trigger (Group "clitest")).
The operator will quiesce, wait for catch-up, and hand off; zero committed rows are lost.
Monitor:  kubectl dba dc-dr status cli-test -n demo
Abort:    kubectl dba dc-dr abort cli-test -n demo

$ kubectl dba dc-dr switchover cli-test -n demo --to dc-c
Error: "dc-c" is not a Member data center of this database (members: [dc-a dc-b]); an Arbiter or Witness DC can never become primary
```

If the database's scope is Global, the command says so: every database in that scope
switches with it.

## status

One-shot DC-DR picture, plus step-by-step switchover progress when one is in flight.

```sh
kubectl dba dc-dr status DB_NAME -n NS
```

**It does not follow.** Run it again to see the next state; that keeps the output
pasteable into a ticket. The step states are derived from the same status fields the
operator's own gates read, so the display cannot drift from the real decision.

```
$ kubectl dba dc-dr status pg-dcdr -n demo
Database:      demo/pg-dcdr (Ready)
Failover scope: primary-dc
                (PlacementPolicy dcdr-postgres failoverPolicy trigger (Global))
Member DCs:    dc-a, dc-b
Active DC:     dc-b   DR phase: Steady
Protected:     true  (data center "dc-a" is streaming 0 bytes behind, within the 16777216 byte budget)

Data centers:
  NAME       ROLE      WRITABLE  HEALTHY   LAG(BYTES)   STREAMER
  dc-a       Member    -         true      0            pg-dcdr-dc-a-1
  dc-b       Member    true      true      -            pg-dcdr-dc-b-1
  dc-c       Arbiter   -         true      -            <none>

No switchover in flight.
Trigger one:  kubectl dba dc-dr status pg-dcdr -n demo --to <dc>
```

During a switchover, the progress block replaces that last section:

```
Planned switchover to "dc-a", started 11s ago:
  [done]    1. target "dc-a" validated: healthy and lag known, within the switchover budget
  [done]    2. quiesce requested on the active DC (dc-b)
  [NOW]     3. quiesce IN EFFECT: active primary write-locked, its LSN frozen
  [pending] 4. target caught up to the frozen LSN (now 0 bytes, needs <= 8192)
  [pending] 5. primary-DC Lease handed off to "dc-a"
  [pending] 6. old DC demoted to standby, annotations cleared, DR phase back to Steady

  NEXT: waiting for the write-lock to take hold on dc-b. This needs the active primary to be UP
        and reachable; a dead primary can never satisfy it (use the failover path).

  Abort:  kubectl dba dc-dr abort cli-test -n demo
  Re-run this command to see the next step; it does not follow.
```

## abort

```sh
kubectl dba dc-dr abort DB_NAME -n NS
```

Sets `dr.kubedb.com/switchover-abort`; its PRESENCE is the signal. The hub clears the
quiesce so the original active DC resumes writes, then removes every switchover
annotation.

**Do not abort by deleting `switchover-to`.** In a scope shared by several databases
the hub re-propagates that annotation to every sibling on each pass, so a bare removal
silently reappears. A switchover that cannot progress also auto-aborts on its own after
`dr.kubedb.com/switchover-timeout` (default 10m).

## accept-data-loss

Releases a cross-DC failover that the RPO budget is holding.

```sh
kubectl dba dc-dr accept-data-loss DB_NAME -n NS --yes
```

When the surviving DC lags more than `spec.replication.bestEffortCrossDCLagBytesForFailover`
(or its lag cannot be measured), promotion is **held**: the lost DC's un-replicated WAL
is unrecoverable, so the choice between an outage and a larger-than-budgeted loss
belongs to a human. This records that decision.

The command prints the current protection verdict first, and `--yes` is mandatory. Both
promotion paths (the hub gate and the coordinator's data-plane gate) honor it within
seconds, and the operator removes the annotation automatically once the failover it
authorized lands, so it cannot linger and approve a later, unrelated loss.

```
$ kubectl dba dc-dr accept-data-loss cli-test -n demo --yes
Current protection verdict: data center "dc-b" is streaming 0 bytes behind, within the 16777216 byte budget
NOTE: the database currently reads protected=true; nothing seems held. Setting the annotation is still
safe: it is only honored by a promotion that the budget is actively refusing, and it is auto-removed after use.
Data-loss acceptance recorded on demo/cli-test.
```

## handoff

Moves the failover authority by annotating the Lease itself.

```sh
kubectl dba dc-dr handoff (DB_NAME | --lease NAME) --to DC --yes
```

This is the **scope-local failover lever**, and the right tool when the active DC's
database is down but its data center is alive (nothing fails over on its own in that
state: the agent renews the Lease on its own health, not the database's). No quiesce
and no catch-up wait happen, so loss is bounded by the RPO budget rather than zero.
For a healthy primary prefer [`switchover`](#switchover).

It moves **every** database in the scope. Never stop a DC's agent to force a failover
instead: one agent serves every scope that DC holds, so stopping it expires all of them.

`--lease` works with no database present, which is how orphan scopes are moved.
Refuses a non-Member target, and refuses outright while a break-glass pin holds the
scope. Re-issuing the same target is safe: the command clears and re-sets the
annotation so the agents always see a real transition.

```
$ kubectl dba dc-dr handoff cli-test -n demo --to dc-b --yes
Handoff of primary-dc-clitest requested: dc-a -> dc-b.
The holder releases within seconds and dc-b acquires on its next retry tick; the annotation clears itself.
Verify:  kubectl dba dc-dr active-dc --lease primary-dc-clitest
```

Measured: the target held the role about 6 seconds later.

## pin-primary

Break-glass. This data center stays primary and stays writable, even through a
coordination-plane outage.

```sh
# run against the SPOKE of the DC being pinned
kubectl dba dc-dr pin-primary (--scope LEASE | --db DB) --yes
kubectl dba dc-dr pin-primary --scope LEASE --remove --yes
```

Creates the human-owned `<scope>-override` ConfigMap on that DC's spoke, which is
precisely why it still works when the hub is unreachable. Two effects: the DC's agent
mirrors the pin onto the Lease so every other Member defers permanently, and the DC's
coordinator forces its leader active regardless of marker state, so a sustained
control-plane outage no longer fences it read-only.

Without `--yes` it prints what it would do and exits non-zero.

```
$ kubectl dba dc-dr pin-primary --scope primary-dc-clitest --yes     # against dc-b
ConfigMap dc-failover/primary-dc-clitest-override created on the current cluster.
This data center is now PINNED primary for scope primary-dc-clitest:
  - its agent mirrors the pin onto the Lease, so no other Member contends;
  - its coordinator keeps the local leader writable even if the marker goes stale.
  IMPORTANT: only honored on the scope's LAST KNOWN HOLDER. On any other DC the agent refuses it.
  IMPORTANT: while pinned, that DC dying means NO failover happens.
```

The pin is then visible from anywhere, and blocks a handoff:

```
$ kubectl dba dc-dr active-dc --lease primary-dc-clitest
Active DC:  dc-b
PINNED:     break-glass override holds this scope on dc-b; it cannot fail over until the
            override ConfigMap is removed from that DC's spoke.

$ kubectl dba dc-dr handoff --lease primary-dc-clitest --to dc-a --yes
Error: scope primary-dc-clitest is PINNED to "dc-b" by a break-glass override; remove that DC's
override ConfigMap first (kubectl dba dc-dr pin-primary --remove), or the handoff cannot complete
```

Caveats, all live-verified: it is honored **only on the scope's last known holder** (on
any other DC the agent refuses it and logs why, so it never promotes a standby); while
it stands there is no split-brain protection for the scope and nothing takes over if
that DC dies; and on control-plane recovery the recreated Lease is claimed first-come,
so keep the pin until the authority is back with the right holder, then remove it.

## pin-standby

This data center never promotes.

```sh
# run against the SPOKE of the DC being held
kubectl dba dc-dr pin-standby (--scope LEASE | --db DB) --yes
kubectl dba dc-dr pin-standby --scope LEASE --remove --yes
```

Creates `<scope>-standby-hold` on that DC's spoke. While it exists the DC never
contends for the Lease, never promotes (it refuses **even an explicit handoff naming
it**, verified live), and its coordinator refuses destructive cross-DC rewinds of the
data it holds. It fails **closed**: an unreadable ConfigMap is treated as held, so a
flaky apiserver never silently drops the protection.

It is deliberately ignored on the DC that is currently active, because demoting the
active DC without a quiesce is unsafe. Move the primary away with a switchover first.

Removing the hold takes effect immediately: the agent re-evaluates that scope's
contention the moment the ConfigMap disappears. Verified live, a handoff that had been
vetoed for 30s completed 10 seconds after the hold was released, with no further
action.

## active-dc

Who holds the primary role, from the authority itself.

```sh
kubectl dba dc-dr active-dc DB_NAME -n NS
kubectl dba dc-dr active-dc --lease primary-dc-orders
kubectl dba dc-dr active-dc DB_NAME -n NS -q          # just the DC name, for scripts
```

Given a database, its scope is resolved exactly as the operator resolves it (the
PlacementPolicy's `failoverPolicy.trigger`, falling back to the legacy
`dr.kubedb.com/failover-group` annotation, else Global) and the matching Lease is read.

```
$ kubectl dba dc-dr active-dc pg-dcdr -n demo
Active DC:  dc-b
Lease:      dc-failover/primary-dc  (scope from PlacementPolicy dcdr-postgres failoverPolicy trigger (Global))
Renewed:    2s ago (lease duration 45s)
Transitions:41
Members:    dc-a,dc-b
```

It warns when the Lease is expired, shows an in-flight handoff, shows a break-glass
pin, and notes when the database's own `status.activeDC` disagrees (the status trails
the Lease by a reconcile; the Lease is the authority).

## debug

```sh
kubectl dba dc-dr debug failover DB_NAME -n NS      # nothing failed over
kubectl dba dc-dr debug switchover DB_NAME -n NS    # switchover not completing
kubectl dba dc-dr debug fence DB_NAME -n NS         # database up but read-only
```

Each walks the real decision chain and reports every condition in order, marking
`OK` / `WARN` / `FAIL`, with a concrete remedy on anything that is not OK, and a count
of blocking conditions at the end.

`debug failover` checks: the database is distributed and armed; the scope is
registered (an unregistered scope has no Lease, so nothing can ever move); the Lease
exists, who holds it, and whether it is still being renewed; break-glass pins and
in-flight handoffs; each candidate DC's health Lease; the RPO budget verdict; stale or
failed ForceFailOver ops; and the DC-DR conditions including the retry cap.

Real output from a database whose failovers had been failing overnight:

```
$ kubectl dba dc-dr debug failover pg-dcdr -n demo
Diagnosing failover for demo/pg-dcdr

  [OK  ] database is DC-DR distributed and armed
  [OK  ] failover scope resolves to primary-dc
         PlacementPolicy dcdr-postgres failoverPolicy trigger (Global)
  [WARN] holder "dc-b" is renewing normally, so the authority will NOT move on its own
         the Lease was renewed 5s ago, inside its 45s duration: that data center's agent is alive
         -> this is by design: no database-level condition (client errors, QPS, lag, a crashed
            postgres) ever moves the Lease. If the DATABASE is down but the DC is alive, either let
            the DC's own raft promote a local peer, or move the scope deliberately:
            kubectl dba dc-dr handoff pg-dcdr -n demo --to <other-dc> --yes
  [OK  ] no break-glass pin on the Lease
  [OK  ] candidate DC "dc-a" is alive (health renewed 0s ago)
  [OK  ] protection is confirmed (RPO budget satisfied)
  [WARN] no failed ForceFailOver ops are blocking
         failed: pg-dcdr-dcdr-failover-92vjc, ... (5 total)
         -> read their status for the real cause, then delete them
  [FAIL] the ForceFailOver retry cap is not tripped
         5 consecutive ForceFailOver attempts to promote data center "dc-b" failed; the hub has
         stopped creating new attempts and needs manual intervention
         -> Fix the underlying failure, delete the failed ops, and it resumes

1 blocking condition(s) found.
```

`debug switchover` reports which gate is holding, and calls out the case people hit
most: every gate measures by dialing the active primary, so an unreachable primary
stalls the switchover permanently (abort and use `handoff`).

`debug fence` explains a database that is up but refusing writes: whether the
authority is healthy, whether the Lease is being renewed (a stale Lease means every
DC's marker goes stale and all of them fence, by design), and what to do when the
coordination plane itself is the broken thing.

---

## Worked examples

**Planned move, monitored.**

```sh
kubectl dba dc-dr switchover pg-dcdr -n demo --to dc-a
kubectl dba dc-dr status pg-dcdr -n demo     # re-run until step 6 is done
```

**The active DC's database is dead but the DC is alive.** Nothing fails over on its
own; that is by design. Confirm, then move the scope:

```sh
kubectl dba dc-dr debug failover pg-dcdr -n demo
kubectl dba dc-dr handoff pg-dcdr -n demo --to dc-b --yes
kubectl dba dc-dr active-dc pg-dcdr -n demo
```

**A failover is held by the RPO budget.**

```sh
kubectl dba dc-dr status pg-dcdr -n demo                    # read protectionMessage
kubectl dba dc-dr accept-data-loss pg-dcdr -n demo --yes
```

**The coordination plane is down and the primary must keep writing.** Run against the
active DC's spoke, before its marker goes stale (roughly 90s from the last renewal):

```sh
kubectl dba dc-dr pin-primary --scope primary-dc --yes --kubeconfig ~/.kube/dc-b.yaml
# ... after recovery, once the Lease is back with the right holder:
kubectl dba dc-dr pin-primary --scope primary-dc --remove --yes --kubeconfig ~/.kube/dc-b.yaml
```

**Keep a data center out of the primary role permanently.**

```sh
kubectl dba dc-dr pin-standby --scope primary-dc --yes --kubeconfig ~/.kube/dc-a.yaml
```
