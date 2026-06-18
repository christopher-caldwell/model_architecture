# Independent Audit: reviewer-04

## Scope Reviewed

I reviewed the architecture decision of whether to merge `server/crates/application` and `server/crates/domain` into one Rust crate. I used the repository instructions, `.understand-anything/knowledge-graph.json`, the required local architecture skill documents, `audits/application-domain-merge-audit.config.yaml`, and representative files from `server/crates/application`, `server/crates/domain`, `server/crates/persistence`, `server/crates/server_bootstrap`, `server/crates/http_server`, `server/crates/graphql_server`, and the workspace Cargo manifests. I did not review other audit reports.

## Executive Summary

I recommend not merging `application` and `domain` for this repository. The current split is not just packaging ceremony; it enforces and documents the central lesson of the reference architecture: domain owns reusable business truth, while application owns CQRS use-case choreography. A merge would reduce one workspace member but would also weaken compile-time dependency signals, broaden adapter dependencies, and make future drift between business rules and use-case orchestration easier.

The strongest argument for a merge would be local ergonomics in a small demo. I did not find enough evidence that this benefit outweighs the architectural cost, especially because transports and persistence already consume the domain crate directly and the composition root already hides most application wiring from transports.

## Findings

### Finding 1: Keep the domain and application crate boundary as an architectural enforcement point

**Severity:** high
**Confidence:** 9
**Category:** design
**Location:** `.understand-anything/knowledge-graph.json` layers; `server/Cargo.toml` lines 1-11; `server/crates/application/Cargo.toml` lines 6-13; `server/crates/domain/Cargo.toml` lines 6-10
**Status:** recommendation

**Issue:** Merging `application` and `domain` would remove a compiler-visible boundary that currently reinforces the repository's core clean/onion architecture rule: use-case orchestration depends on business rules, not the other way around.

**Evidence:** The knowledge graph identifies separate `Domain` and `Application CQRS` layers. The workspace currently lists `crates/application` and `crates/domain` as distinct members in `server/Cargo.toml` lines 2-9. `application` depends on `domain` at `server/crates/application/Cargo.toml` line 7, while `domain` has no dependency on `application` in `server/crates/domain/Cargo.toml` lines 6-10. Representative command code in `server/crates/application/src/commands/lending.rs` lines 63-103 loads state, calls domain guards, persists, and commits. Representative domain code in `server/crates/domain/src/member/logic.rs` lines 5-58 and `server/crates/domain/src/book_copy/logic.rs` lines 5-63 owns transitions, guards, and preparation defaults.

**Impact:** In a reference implementation, the boundary is part of the teaching and review surface. If the crates are merged, the project loses a simple mechanical check that `domain` remains independent from CQRS orchestration. Reviewers would have to rely more heavily on convention, module naming, and discipline to prevent domain rules, application workflows, and ports from drifting together.

**Recommended Fix:** Keep `server/crates/domain` and `server/crates/application` as separate crates. If import ergonomics are the concern, prefer narrow re-exports from `server_bootstrap` or small module cleanup inside the existing crates rather than collapsing the boundary.

### Finding 2: A merge would move adapter coupling problems rather than eliminate them

**Severity:** high
**Confidence:** 8
**Category:** maintainability
**Location:** `server/crates/domain/src/uow.rs` lines 8-50; `server/crates/persistence/Cargo.toml` lines 6-12; `server/crates/persistence/src/uow.rs` lines 3-14 and 35-148; `server/crates/server_bootstrap/src/deps.rs` lines 64-111
**Status:** risk

**Issue:** The current domain crate owns repository and unit-of-work ports that persistence implements without depending on application. If `domain` is merged into `application`, persistence must either depend on a broader combined crate that also contains command/query use cases, or the ports must be extracted again into another boundary.

**Evidence:** `domain::uow` defines `UnitOfWorkPort`, `WriteUnitOfWork`, and `WriteUnitOfWorkFactory` in `server/crates/domain/src/uow.rs` lines 8-50. `persistence` depends on `domain`, not `application`, in `server/crates/persistence/Cargo.toml` line 7. `SqlUnitOfWork` implements domain ports in `server/crates/persistence/src/uow.rs` lines 35-148. The composition root wires persistence adapters into application commands and queries in `server/crates/server_bootstrap/src/deps.rs` lines 68-92.

**Impact:** A merged crate would make persistence depend on a larger "core" artifact that includes use-case orchestration it does not need. That broadens compile-time coupling and makes it less obvious that persistence is an adapter implementing core ports rather than a collaborator with application workflows. Extracting the ports to avoid this would recreate the same conceptual boundary under another name.

**Recommended Fix:** Preserve the current split. If the team still merges, create strict top-level modules such as `core::domain` and `core::application`, keep persistence imports limited to `core::domain::{..., port, uow}`, and add an architectural check that persistence cannot import command/query modules.

### Finding 3: The merge has low payoff because outer layers already consume stable domain and bootstrap surfaces

**Severity:** medium
**Confidence:** 8
**Category:** maintainability
**Location:** `server/crates/http_server/Cargo.toml` lines 6-23; `server/crates/graphql_server/Cargo.toml` lines 6-21; `server/crates/http_server/src/router/loan/post_handlers.rs` lines 25-40; `server/crates/graphql_server/src/router/graphql/lending/mutations.rs` lines 13-30; `server/crates/server_bootstrap/src/deps.rs` lines 14-15
**Status:** observation

**Issue:** Collapsing the two inner crates would not materially simplify the transport call shape. HTTP and GraphQL already call through `ServerDeps` and application commands, while direct domain imports are mostly for DTO conversion and typed error mapping.

**Evidence:** HTTP depends on `domain` and `server_bootstrap` in `server/crates/http_server/Cargo.toml` lines 7-9. GraphQL does the same in `server/crates/graphql_server/Cargo.toml` lines 7-9. The HTTP checkout handler parses a request, builds `CheckOutBookCopyInput`, calls one command, and maps the result in `server/crates/http_server/src/router/loan/post_handlers.rs` lines 25-40. The GraphQL mutation follows the same pattern in `server/crates/graphql_server/src/router/graphql/lending/mutations.rs` lines 13-30. `server_bootstrap` re-exports application commands and queries in `server/crates/server_bootstrap/src/deps.rs` lines 14-15.

**Impact:** A merge mostly renames imports or broadens dependencies; it does not remove the need for domain DTO mapping, typed domain errors, command/query inputs, or composition-root wiring. The expected simplification is therefore small compared with the lost boundary signal.

**Recommended Fix:** Do not merge for import-count reduction alone. If the desired improvement is API ergonomics, expose a narrower facade from `server_bootstrap` or application command modules while keeping the domain crate independently consumable.

### Finding 4: If a merge is chosen anyway, it needs explicit architectural guardrails

**Severity:** medium
**Confidence:** 7
**Category:** design
**Location:** `server/crates/application/src/commands/error.rs` lines 1-19; `server/crates/http_server/src/router/errors.rs` lines 51-96; `server/crates/application/src/commands/lending.rs` lines 75-96 and 117-159
**Status:** recommendation

**Issue:** A merged crate can still preserve clean/onion layering, but only if it adds replacement safeguards for the lost crate boundary. Without those safeguards, command code can more easily absorb domain decisions, and domain code can more easily grow use-case concerns.

**Evidence:** `CommandError` currently wraps typed domain errors in `server/crates/application/src/commands/error.rs` lines 1-19, and HTTP maps those errors at the transport edge in `server/crates/http_server/src/router/errors.rs` lines 51-96. The lending command calls domain guards and transitions rather than reimplementing policy in `server/crates/application/src/commands/lending.rs` lines 75-96 and 117-159. These are good existing shapes, but after a merge they become conventions within one crate rather than cross-crate dependency constraints.

**Impact:** The likely failure mode is gradual erosion rather than immediate breakage: commands may start branching directly on fields, domain modules may start importing application-only services, and ports may become less narrow because all symbols live under the same package.

**Recommended Fix:** If merging proceeds, require all of the following: keep separate `domain` and `application` top-level modules; deny imports from `application` inside `domain`; keep typed domain errors separate from `CommandError`; keep repository and unit-of-work ports in the domain module; add architecture tests or a static import check; and document that the merge is packaging-only, not a responsibility merge.

## Non-Issues / Things Checked

- The current direct dependency from transports to `domain` looks acceptable for DTO mapping and protocol-specific error mapping; handlers and resolvers reviewed still call one command/query rather than composing business workflows.
- The current `application -> domain` dependency is one-way in manifests and representative imports.
- Persistence implements domain ports and does not need to depend on application commands or queries today.
- The command examples reviewed avoid post-write hydration and reconstruct returned values from loaded state, generated values, and code-side timestamps.
- Keeping repository ports in `domain` is consistent with this repository's documented style, even though some onion variants place ports in application.

## Assumptions

- The decision under review is a packaging and architecture-boundary decision, not a request to implement a refactor.
- The repository is intended to remain a clean/onion architecture + CQRS reference implementation, so clarity and enforceable boundaries are part of the product value.
- Build-time crate count is not currently causing a measurable performance or developer-experience problem.
- The reviewed representative files are sufficient to assess the architectural dependency direction and intended responsibility split.

## Open Questions

- What concrete pain is motivating the merge: compile time, import ergonomics, conceptual simplicity for readers, or preparation for a different deployment/package shape?
- Would the team accept a smaller facade or re-export cleanup instead of a physical crate merge?
- If merged, will the project add automated architecture checks to replace the current crate-level dependency enforcement?
- Is the long-term goal to keep persistence depending only on domain contracts, or is depending on a combined application/domain crate considered acceptable for this reference implementation?
