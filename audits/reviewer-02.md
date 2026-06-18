# Independent Audit: reviewer-02

## Scope Reviewed

I reviewed the architecture decision of whether to merge `server/crates/application` and `server/crates/domain` into one crate. Materials reviewed included the repository instructions in the prompt, `.understand-anything/knowledge-graph.json`, the repo-local architecture skills, `audits/application-domain-merge-audit.config.yaml`, workspace and crate `Cargo.toml` files, and representative files in `server/crates/application`, `server/crates/domain`, `server/crates/persistence`, `server/crates/server_bootstrap`, `server/crates/http_server`, and `server/crates/graphql_server`.

I did not review any other reviewer reports or consensus output.

## Executive Summary

I recommend not merging `application` and `domain` for this repository. The current split is doing useful architectural work: Rust crate dependencies enforce the intended inward direction, keep the domain as a small stable core, and force persistence and transports to depend only on the domain types and ports they actually need. A merge would be mechanically possible, but it would replace compile-time boundaries with convention in a project whose main purpose is to demonstrate clean/onion architecture and CQRS boundaries.

The strongest counterargument is that the combined code size is small and the application/domain split has some conceptual overlap around repository and unit-of-work ports. That is a reason to clarify port ownership, not to collapse the crate boundary.

## Findings

### Finding 1: Keep `domain` as a separate crate because it is the primary enforceable inward boundary

**Severity:** high  
**Confidence:** 9  
**Category:** design  
**Location:** `server/crates/application/Cargo.toml:6-13`, `server/crates/domain/Cargo.toml:6-10`, `server/crates/application/src/commands/lending.rs:1-8`, `server/crates/domain/src/member/logic.rs:5-58`, `server/crates/domain/src/book_copy/logic.rs:5-62`  
**Status:** recommendation

**Issue:** Merging `application` and `domain` would remove the current compile-time one-way dependency from application use cases into domain business rules. Inside one crate, domain modules could import application commands, application ports, or orchestration helpers without Cargo detecting a layer violation.

**Evidence:** `application` explicitly depends on `domain` in `server/crates/application/Cargo.toml:6-13`, while `domain` has no dependency back to `application` in `server/crates/domain/Cargo.toml:6-10`. Application commands use domain types and ports, for example `LendingCommands` imports `domain::{book_copy, loan, member, uow}` in `server/crates/application/src/commands/lending.rs:1-8`. The domain currently owns reusable business decisions such as `Member::ensure_can_borrow`, `Member::ensure_within_loan_limit`, and `MemberCreationPayload::prepare` in `server/crates/domain/src/member/logic.rs:5-58`, plus book-copy transitions and defaults in `server/crates/domain/src/book_copy/logic.rs:5-62`.

**Impact:** The repository is a clean/onion architecture reference implementation. Its central teaching point is that business decisions remain in domain code and use-case choreography remains in application code. The current crate split makes that direction visible and enforceable. A merge would make the same separation depend on code review discipline and module naming. That is a meaningful regression for a reference implementation even if the runtime behavior would not change.

**Recommended Fix:** Do not merge the crates. Keep `domain` as the leaf crate and `application` as the CQRS use-case crate. If import ergonomics are the problem, prefer a small facade or re-exports from `server_bootstrap` for transport-facing application inputs and command/query types, while preserving `application -> domain`.

### Finding 2: A merge would broaden adapter dependencies and weaken the port boundary

**Severity:** medium  
**Confidence:** 8  
**Category:** maintainability  
**Location:** `server/crates/persistence/Cargo.toml:6-12`, `server/crates/persistence/src/uow.rs:3-14`, `server/crates/persistence/src/uow.rs:35-148`, `server/crates/server_bootstrap/Cargo.toml:6-12`, `server/crates/server_bootstrap/src/deps.rs:68-111`  
**Status:** risk

**Issue:** Persistence currently depends only on the domain crate for entities, repository traits, and the unit-of-work interface. If application and domain are merged, persistence must depend on a larger crate that also contains commands, queries, and application-owned service ports.

**Evidence:** `persistence` depends on `domain`, not `application`, in `server/crates/persistence/Cargo.toml:6-12`. Its SQL unit of work imports and implements domain ports in `server/crates/persistence/src/uow.rs:3-14` and `server/crates/persistence/src/uow.rs:35-148`. Bootstrap is the composition root that depends on both `application` and `persistence` in `server/crates/server_bootstrap/Cargo.toml:6-12`, then wires SQL repositories into application commands and queries in `server/crates/server_bootstrap/src/deps.rs:68-111`.

**Impact:** The current shape keeps adapters pointed at narrow core contracts. After a merge, the adapter crate would be able to reference application commands/queries directly because those symbols would live in the same dependency as the domain ports. That does not create an immediate Rust cycle, but it does make accidental adapter-to-use-case coupling easier and less detectable. It also makes the core dependency less replaceable because adapter code no longer consumes the smallest relevant crate.

**Recommended Fix:** Preserve `persistence -> domain` and `server_bootstrap -> application + persistence`. If a future merge is still pursued, add explicit replacement guardrails before merging: top-level `domain` and `application` modules, restricted public re-exports, dependency/import checks in CI, and documented rules forbidding persistence from importing application modules. Those guardrails would still be weaker than the current crate boundary.

### Finding 3: The current transport symmetry depends on a shared application layer; merging does not improve that and may make the example less clear

**Severity:** medium  
**Confidence:** 8  
**Category:** design  
**Location:** `server/crates/http_server/src/router/loan/post_handlers.rs:25-40`, `server/crates/graphql_server/src/router/graphql/lending/mutations.rs:13-31`, `server/crates/application/src/commands/lending.rs:59-103`  
**Status:** observation

**Issue:** One goal of the repository is to show HTTP and GraphQL as replaceable adapters over the same use case. The current separate `application` crate makes that use-case layer explicit. A combined application/domain crate would still work, but it would make the boundary less visible to readers.

**Evidence:** The HTTP checkout endpoint parses a request, builds `CheckOutBookCopyInput`, calls `deps.lending.commands.check_out_book_copy`, and maps the response in `server/crates/http_server/src/router/loan/post_handlers.rs:25-40`. The GraphQL mutation follows the same shape in `server/crates/graphql_server/src/router/graphql/lending/mutations.rs:13-31`. The actual workflow lives in `server/crates/application/src/commands/lending.rs:59-103`, where the command loads locked state, asks domain guards, checks loan count, writes, commits, and returns the result.

**Impact:** The separate application crate is a useful signpost: transports call one command/query, while domain code remains reusable business logic. Combining the crates would make the architecture easier to misread as a generic "core" crate rather than a CQRS application layer over a domain model.

**Recommended Fix:** Keep the separate `application` crate. If discoverability is a problem, improve documentation or module docs around command/query entry points instead of merging the physical crates.

### Finding 4: Port ownership has some ambiguity, but that does not justify a full application/domain merge

**Severity:** low  
**Confidence:** 7  
**Category:** maintainability  
**Location:** `server/crates/domain/src/uow.rs:8-51`, `server/crates/domain/src/member/port.rs:6-18`, `server/crates/domain/src/loan/port.rs:11-26`, `server/crates/application/src/ports/gen_ident.rs:1-3`  
**Status:** observation

**Issue:** The domain crate owns repository and unit-of-work ports, while the application crate owns at least one external service port, `IdentGeneratorPort`. A reviewer could reasonably ask whether transaction-oriented ports are application concerns rather than pure domain concerns.

**Evidence:** `UnitOfWorkPort` and `WriteUnitOfWorkFactory` live in `server/crates/domain/src/uow.rs:8-51`. Repository ports live in domain modules such as `server/crates/domain/src/member/port.rs:6-18` and `server/crates/domain/src/loan/port.rs:11-26`. By contrast, `IdentGeneratorPort` lives in `server/crates/application/src/ports/gen_ident.rs:1-3`, and bootstrap implements it in `server/crates/server_bootstrap/src/deps.rs:48-54`.

**Impact:** This is the main legitimate pressure toward a less granular core boundary. However, the repository's documented rule is that domain includes repository ports and unit-of-work traits. The current design is internally consistent with that rule and keeps persistence dependent only on domain. Merging the crates would resolve the naming discomfort by removing the enforcement boundary, which is a disproportionate fix.

**Recommended Fix:** Keep the crates separate. If the team wants stricter terminology, document why repository and unit-of-work ports are domain-owned in this reference implementation, or consider a smaller targeted change such as a clearly named `domain::ports`/`domain::uow` section. Do not use this ambiguity as the reason to merge application and domain wholesale.

## Non-Issues / Things Checked

- The small combined size of `application` and `domain` is not, by itself, an architectural reason to merge. The split is about dependency direction and review clarity, not line count.
- I did not find evidence that application currently imports Axum, async-graphql, SQLx, or transport DTOs.
- I checked representative HTTP and GraphQL paths and found they call the same application command rather than duplicating the checkout workflow.
- Domain logic has local tests for key guards, transitions, enum parsing, and preparation defaults; merging crates would not obviously improve this coverage.
- Direct transport dependencies on `domain` for DTO mapping and error mapping are acceptable in the current architecture because transports still call application use cases for workflow.

## Assumptions

- The primary goal remains a clean/onion architecture + CQRS reference implementation, not minimizing crate count.
- There is no measured compile-time, publishing, or dependency-management problem severe enough to outweigh architectural clarity.
- A proposed merge would be a physical crate merge while intending to preserve the same conceptual application/domain separation.
- The repository's documented rule that domain owns repository ports and unit-of-work traits is intentional.

## Open Questions

- What concrete pain is motivating the merge: compile time, dependency ergonomics, teaching simplicity, public API clutter, or something else?
- Would the team accept CI import-boundary checks if the crates were merged, and what tool would enforce them reliably inside a single crate?
- Is the long-term goal to demonstrate replaceable persistence adapters, multiple transports, and reusable domain rules to readers who benefit from physical crate boundaries?
- Should transport crates continue to depend directly on `domain`, or should `server_bootstrap` re-export selected domain/application-facing types to reduce transport manifest dependencies?
