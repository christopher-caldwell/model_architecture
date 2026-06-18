# Final Consensus Audit

## Summary

Audit target: whether to merge `server/crates/application` and `server/crates/domain` into one crate in this Rust clean/onion architecture + CQRS reference implementation.

Five independent identical reviewers completed the architecture audit. The quorum threshold was 0.60, so findings supported by at least 3 of 5 reviewers meet quorum.

Overall result: unanimous recommendation not to merge `application` and `domain` as the default architecture decision. The reviewers agreed that the current split is not empty ceremony; it is the primary Cargo-enforced expression of the repo's inward dependency rule. A merge remains mechanically possible, but only as a packaging-only move with explicit replacement guardrails.

## Configuration Used

```yaml
audit_target: "Architecture decision: whether to merge server/crates/application and server/crates/domain into one crate in this Rust clean/onion architecture + CQRS reference implementation."
audit_type: "architecture_audit"
reviewer_count_requested: 5
reviewer_count_completed: 5
reviewer_mode: "identical_reviewers"
quorum_threshold: 0.60
include_minority_findings: true
minority_severity_floor: "high"
output_directory: "./audits"
individual_report_pattern: "reviewer-{number}.md"
final_report_filename: "final-consensus-audit.md"
```

## Consensus Findings

### Finding 1: Do Not Merge By Default; The Crate Boundary Enforces The Core Dependency Rule

**Severity:** high  
**Confidence:** 10/10  
**Support:** 5/5 reviewers  
**Quorum Status:** unanimous  
**Category:** design  
**Root Verification:** verified

**Issue:**  
Merging `application` and `domain` would remove the current compile-time boundary that enforces the intended direction: application use cases depend inward on domain business rules, while domain cannot depend on use-case orchestration.

**Evidence:**  
All reviewers cited the manifest graph. Root verification confirms `server/Cargo.toml` declares separate `crates/application` and `crates/domain` workspace members, `server/crates/application/Cargo.toml` depends on `domain`, and `server/crates/domain/Cargo.toml` does not depend on `application`. The README and project skills describe the same conceptual split: domain owns business rules; application owns commands/queries and use cases.

Representative code matches the split. `server/crates/application/src/commands/lending.rs` builds a unit of work, loads state, calls domain guards, writes, and commits. `server/crates/domain/src/member/logic.rs` and related domain modules own guards, transitions, preparation defaults, and typed business errors.

**Why It Matters:**  
The repository is a reference implementation. The package boundary is part of the architecture, not just a file organization preference. After a merge, the same rule would rely on convention, review discipline, and optional tooling instead of Cargo.

**Recommended Action:**  
Keep `application` and `domain` separate unless there is a concrete measured pain that outweighs the loss of enforcement. Prefer smaller ergonomic fixes first: clearer docs, narrower re-exports, cleanup of empty modules, or facade exports from `server_bootstrap`.

**Reviewers Supporting:**  
`reviewer-01`, `reviewer-02`, `reviewer-03`, `reviewer-04`, `reviewer-05`

### Finding 2: A Merge Would Broaden Persistence From Domain Contracts To A Larger Core Surface

**Severity:** high  
**Confidence:** 9/10  
**Support:** 5/5 reviewers  
**Quorum Status:** unanimous  
**Category:** maintainability  
**Root Verification:** verified

**Issue:**  
Persistence currently depends only on `domain` to implement repository ports and the write unit-of-work abstraction. If application and domain are merged, persistence must depend on a broader crate that also contains command/query use cases, unless another boundary is recreated.

**Evidence:**  
Root verification confirms `server/crates/persistence/Cargo.toml` depends on `domain`, not `application`. `server/crates/domain/src/uow.rs` defines `UnitOfWorkPort`, `WriteUnitOfWork`, and `WriteUnitOfWorkFactory`. `server/crates/persistence/src/uow.rs` imports domain entities and ports, implements `UnitOfWorkPort`, and implements each domain write repository port. `server_bootstrap` is the current composition root that wires persistence implementations into application commands and queries.

**Why It Matters:**  
The current shape makes persistence an adapter for domain contracts. A merged crate would make application use-case internals more visible to persistence and increase the chance of accidental coupling. If the project solves that by extracting ports or adding import checks, it has effectively recreated the same boundary under another name.

**Recommended Action:**  
Preserve `persistence -> domain` and `server_bootstrap -> application + persistence`. If a merge is ever pursued, require persistence to import only a `domain` submodule or facade and add an automated boundary check.

**Reviewers Supporting:**  
`reviewer-01`, `reviewer-02`, `reviewer-03`, `reviewer-04`, `reviewer-05`

### Finding 3: The Current Split Makes The Shared HTTP/GraphQL Use-Case Story Clear

**Severity:** medium  
**Confidence:** 8/10  
**Support:** 5/5 reviewers  
**Quorum Status:** unanimous  
**Category:** design  
**Root Verification:** verified

**Issue:**  
HTTP and GraphQL already share the same application use cases while mapping domain values and typed errors at the edge. A merge would not improve that behavior and could make the architectural example less clear.

**Evidence:**  
Reviewers cited the transport pattern: handlers and resolvers call `ServerDeps` command/query handles rather than composing workflows themselves. Root verification confirms both transport manifests depend on `domain` plus `server_bootstrap`; `server_bootstrap` re-exports application command/query types and wires repositories into `ServerDeps`.

**Why It Matters:**  
The separate `application` crate is a visible signpost for "transports call one use case." The separate `domain` crate remains a stable edge-mapping contract for DTO conversion and error mapping. Combining them would broaden what transports can see unless facade exports are carefully constrained.

**Recommended Action:**  
Keep the split. If the pain is transport import clutter, consider a narrower transport-facing facade instead of merging the inner crates.

**Reviewers Supporting:**  
`reviewer-01`, `reviewer-02`, `reviewer-03`, `reviewer-04`, `reviewer-05`

### Finding 4: The Benefit Of Merging Appears Mostly Ergonomic And Low Payoff

**Severity:** low  
**Confidence:** 8/10  
**Support:** 5/5 reviewers  
**Quorum Status:** unanimous  
**Category:** maintainability  
**Root Verification:** verified

**Issue:**  
The main observed benefit is fewer crates/imports. Reviewers did not find evidence that the split is causing duplicated workflows, incorrect layering, or significant implementation friction.

**Evidence:**  
The top-level `application` and `domain` `lib.rs` files are small. Application commands and queries are thin but meaningful orchestration layers. Domain modules contain real guards, transitions, typed errors, and preparation methods. The current use-case flows are already centralized in application.

**Why It Matters:**  
Reducing package count is a valid production preference in some projects, but here it trades away a clear architecture boundary in a repo whose purpose is to demonstrate that boundary.

**Recommended Action:**  
Do not merge for convenience alone. First clarify port ownership, remove/complete the empty application services module if desired, and improve re-exports or docs.

**Reviewers Supporting:**  
`reviewer-01`, `reviewer-02`, `reviewer-03`, `reviewer-04`, `reviewer-05`

### Finding 5: If A Merge Happens Anyway, It Must Be Packaging-Only With Replacement Guardrails

**Severity:** medium  
**Confidence:** 8/10  
**Support:** 5/5 reviewers  
**Quorum Status:** unanimous  
**Category:** design  
**Root Verification:** partially verified

**Issue:**  
A single crate can preserve clean/onion layering only if the project deliberately replaces the lost Cargo boundary with module boundaries, visibility rules, architecture tests, and documentation.

**Evidence:**  
Every reviewer proposed some version of the same fallback: a `core` crate with explicit `domain` and `application` modules, narrow public exports, and checks preventing `domain` from importing `application`, transports, persistence, or framework concerns. Root verification confirms the current code shape could be represented as modules, but no replacement boundary tooling currently exists in the repo.

**Why It Matters:**  
Without guardrails, the expected failure mode is gradual erosion: commands start inspecting fields instead of asking domain guards, domain modules start importing application services, or transports begin composing workflows directly.

**Recommended Action:**  
If there is a strong external reason to merge, do it only as:

```text
core
  domain/
  application/
```

Then add import-boundary checks, keep typed domain errors separate from `CommandError`, keep repository/UoW contracts under the domain-facing API, and make transports reach use cases only through bootstrap/facade handles.

**Reviewers Supporting:**  
`reviewer-01`, `reviewer-02`, `reviewer-03`, `reviewer-04`, `reviewer-05`

## Strong Minority Findings

No high-severity minority finding contradicted the consensus recommendation.

Two lower-severity observations appeared in only part of the reviewer set:

- Domain-only tests would become less visibly isolated in a merged crate. This is plausible but low severity and does not change the recommendation.
- Repository and unit-of-work port ownership is somewhat ambiguous because this project places those abstractions in `domain`. Reviewers treated this as a clarification/documentation issue, not a reason to merge.

## Split Or Disputed Findings

There was no meaningful split on the final recommendation. All reviewers recommended keeping the crates separate by default.

The only nuance is fallback posture:

- Some reviewers framed a single `core` crate as acceptable if packaging pressure is real.
- Others treated the current split as the better reference-architecture choice unless a concrete pain is proven.

These are compatible positions: the consensus is "do not merge by default; if forced, merge only with explicit guardrails."

## Likely False Positives

### Domain Is Not Pure, Therefore It Should Be Merged

**Support:** 0/5 as a recommendation; 2/5 noted related ambiguity  
**Why Not Accepted:**  
The domain crate includes async repository ports and UoW traits, so it is not a purely behavioral domain model. But the repository's own rules explicitly put repository ports and UoW traits in `domain`, and persistence currently benefits from depending on that smaller contract. This is a design choice to document or revisit, not sufficient evidence for merging application and domain.

### Small Crate Count Alone Justifies A Merge

**Support:** 0/5  
**Why Not Accepted:**  
All reviewers agreed the current code is small, but they treated the split as an architectural enforcement point. No reviewer identified enough concrete friction to justify removing that enforcement.

## Final Recommendation

Keep `server/crates/application` and `server/crates/domain` as separate crates.

The combined recommendation is stronger than the initial intuition that the split may be no-benefit. The reviewers found a real benefit: Cargo enforces the inward dependency direction, keeps persistence pointed at domain contracts instead of application use cases, and keeps the reference architecture legible.

The most pragmatic next step is not a merge. It is a small architecture cleanup pass:

1. Document why repository and unit-of-work ports live in `domain` for this reference implementation.
2. Clean up minor ergonomics, such as empty application port modules or noisy re-exports.
3. Consider a narrow facade if transport imports feel too broad.
4. Only revisit a packaging-only `core` crate if there is measured friction, and add boundary checks before merging.

## Caveats

- All five reviewers completed successfully.
- This audit evaluated the architecture decision, not a concrete refactor branch.
- Root verification checked representative files, manifests, and project rules, but did not run tests because no code behavior changed.
- The recommendation assumes this repo remains a clean/onion + CQRS reference implementation where explicit, enforceable boundaries are part of the value.

## Audit Files Reviewed

- `./audits/reviewer-01.md`
- `./audits/reviewer-02.md`
- `./audits/reviewer-03.md`
- `./audits/reviewer-04.md`
- `./audits/reviewer-05.md`
