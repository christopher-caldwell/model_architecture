# Independent Audit: reviewer-01

## Scope Reviewed

I reviewed the architecture decision of whether to merge `server/crates/application` and `server/crates/domain` into one crate. Materials reviewed were the repository instructions in the prompt, `.understand-anything/knowledge-graph.json`, the required `.agents/skills/*` architecture guides, `audits/application-domain-merge-audit.config.yaml`, the relevant `Cargo.toml` files under `server/`, and representative files from `server/crates/application`, `server/crates/domain`, `server/crates/persistence`, `server/crates/server_bootstrap`, `server/crates/http_server`, and `server/crates/graphql_server`.

I did not read any other reviewer report or final synthesis file.

## Executive Summary

My recommendation is not to merge `application` and `domain` for this repository in its current form. The current split is not accidental packaging overhead; it is the main compile-time enforcement mechanism for the reference architecture's dependency direction. `domain` owns entities, invariants, errors, repository ports, and unit-of-work traits. `application` owns CQRS use-case orchestration and depends inward on `domain`. `persistence` currently implements only `domain` ports without depending on application use cases.

A merge could still preserve conceptual module boundaries if implemented carefully, but it would replace Rust crate-level enforcement with convention. For a clean/onion + CQRS reference implementation whose purpose is to demonstrate replaceable adapters and stable core boundaries, that is a meaningful regression. The only clear benefit I found is modest workspace simplification, which does not outweigh the increased coupling and review burden.

## Findings

### Finding 1: Merging removes the strongest enforcement of the domain/application boundary

**Severity:** high  
**Confidence:** 9  
**Category:** design  
**Location:** `server/Cargo.toml`; `server/crates/domain/Cargo.toml`; `server/crates/application/Cargo.toml`; `.understand-anything/knowledge-graph.json` layer map; `.agents/skills/onion-cqrs-project-map/SKILL.md`  
**Status:** risk

**Issue:** The current architecture uses separate Rust crates to enforce that `domain` cannot depend on `application`, transports, persistence, SQLx, Axum, async-graphql, bootstrap wiring, or application-specific services. Merging the crates would make that boundary primarily a module convention inside one compilation unit.

**Evidence:** The workspace declares `crates/application` and `crates/domain` as separate members in `server/Cargo.toml`. `server/crates/application/Cargo.toml` depends on `domain = { path = "../domain" }`, while `server/crates/domain/Cargo.toml` has no dependency on `application`. The knowledge graph also classifies them as separate layers: `Domain` for business entities, invariants, typed domain errors, transitions, and repository ports; `Application CQRS` for commands, queries, transactions, ports, and domain decisions. The project guide states the intended direction as `transport -> server_bootstrap -> application -> domain`, with `persistence -> domain`.

**Impact:** In a reference architecture, the dependency boundary is part of the lesson. With separate crates, an accidental domain import of application use cases is impossible at compile time. After a merge, code in a `domain` module could call or import command/query code unless the project adds replacement guardrails. That increases the chance that reusable business decisions drift into use-case orchestration, weakening clean/onion separation.

**Recommended Fix:** Keep `server/crates/application` and `server/crates/domain` separate. If the team still merges them, require explicit internal modules such as `core::domain` and `core::application`, add architecture linting or dependency-boundary checks, and document that `domain` modules must not import application modules.

### Finding 2: A merge broadens persistence's dependency from domain ports to the full application surface

**Severity:** medium  
**Confidence:** 9  
**Category:** maintainability  
**Location:** `server/crates/persistence/Cargo.toml`; `server/crates/persistence/src/uow.rs:3`; `server/crates/domain/src/uow.rs:1`; `server/crates/domain/src/member/port.rs:5`; `server/crates/server_bootstrap/src/deps.rs:4`  
**Status:** risk

**Issue:** Persistence currently depends only on the domain crate to implement repository ports and unit-of-work traits. If application and domain become one crate, persistence would need to depend on a crate that also contains CQRS commands, queries, application inputs, and application-owned service ports.

**Evidence:** `server/crates/persistence/Cargo.toml` depends on `domain`, not `application`. `server/crates/persistence/src/uow.rs` implements `domain::uow::UnitOfWorkPort` and repository traits such as `BookWriteRepoPort`, `BookCopyWriteRepoPort`, `LoanWriteRepoPort`, and `MemberWriteRepoPort`. `server/crates/domain/src/uow.rs` defines the unit-of-work abstraction in terms of domain repository ports. `server/crates/domain/src/member/port.rs` defines the member read/write repository ports. Separately, `server/crates/server_bootstrap/src/deps.rs` wires persistence adapters into application commands and queries.

**Impact:** The current shape lets persistence be an adapter for domain contracts. A merged crate would make persistence depend on the command/query layer even if it only uses domain-facing ports. That widens the adapter's compile-time dependency and makes it easier for persistence to start relying on application use-case types or application service ports. It also pulls application crate dependencies such as `nanoid` and `tracing` into any consumer that only needs domain contracts.

**Recommended Fix:** Keep repository ports and unit-of-work abstractions in a crate that persistence can depend on without importing application use cases. In practice, that means preserving the existing `domain` crate. If merging is mandatory, consider a three-way split instead: a minimal `domain` or `core_ports` crate for entities and ports, plus application use cases separately.

### Finding 3: The current split cleanly supports shared HTTP and GraphQL use cases

**Severity:** medium  
**Confidence:** 8  
**Category:** design  
**Location:** `server/crates/http_server/src/router/loan/post_handlers.rs:25`; `server/crates/graphql_server/src/router/graphql/lending/mutations.rs:13`; `server/crates/application/src/commands/lending.rs:59`; `server/crates/http_server/src/router/errors.rs:52`; `server/crates/graphql_server/src/router/graphql/mod.rs:67`  
**Status:** observation

**Issue:** The code already demonstrates the desired clean/onion outcome: transports are thin adapters over shared application use cases, while domain errors remain typed until protocol-specific mapping. A crate merge does not improve this behavior and may obscure why the split exists.

**Evidence:** The HTTP checkout handler builds `CheckOutBookCopyInput`, calls `deps.lending.commands.check_out_book_copy(input)`, and maps the result to `LoanResponseBody`. The GraphQL mutation performs the same use case through the same `LendingCommands::check_out_book_copy`. The command itself loads locked state, calls domain decisions such as `member.ensure_can_borrow()` and `book_copy.ensure_can_be_borrowed()`, persists through the unit of work, commits once, and returns a domain value. HTTP and GraphQL error mappers match on the same `CommandError` variants and domain error types at their protocol edges.

**Impact:** This is strong evidence that the current crate boundary is carrying architectural value, not just ceremony. Merging application and domain would not reduce duplicated transport workflows because those workflows are already centralized. The likely change is weaker visibility into which layer owns each decision.

**Recommended Fix:** Retain the split and treat it as part of the reference design. Any simplification should target smaller issues, such as naming or re-export ergonomics, without collapsing the domain/application crate boundary.

### Finding 4: Domain tests and domain-only dependencies benefit from independent compilation

**Severity:** low  
**Confidence:** 8  
**Category:** testing  
**Location:** `server/crates/domain/src/member/logic.rs:61`; `server/crates/domain/src/book_copy/logic.rs:65`; `server/crates/domain/src/loan/logic.rs:29`; `server/crates/domain/Cargo.toml`; `server/crates/application/Cargo.toml`  
**Status:** risk

**Issue:** The domain crate currently has focused tests for guards, transitions, preparation defaults, and enum behavior, with a small dependency surface. Merging application into the same crate would make these tests live in a larger crate with unrelated application dependencies and public surface.

**Evidence:** `member::logic`, `book_copy::logic`, and `loan::logic` contain unit tests for state transitions, guard failures, and preparation defaults. The domain manifest depends only on `anyhow`, `async-trait`, `chrono`, and `thiserror`. The application manifest adds application-use-case dependencies such as `domain`, `nanoid`, and `tracing`.

**Impact:** This is not a correctness break by itself, but it weakens the fast, isolated feedback loop around business rules. As the application layer grows, domain-only test runs and dependency review become noisier.

**Recommended Fix:** Keep domain logic and its tests in the dedicated crate. If a merge proceeds, preserve a strict internal domain test module structure and ensure domain tests do not require application fixtures or application service setup.

### Finding 5: The main benefit of merging is ergonomic, not architectural

**Severity:** informational  
**Confidence:** 8  
**Category:** maintainability  
**Location:** `server/crates/application/src/commands/lending.rs:3`; `server/crates/application/src/commands/error.rs:1`; `server/crates/server_bootstrap/src/lib.rs:4`; `server/crates/http_server/Cargo.toml`; `server/crates/graphql_server/Cargo.toml`  
**Status:** observation

**Issue:** The codebase has some import and re-export overhead from multiple crates, but I did not find evidence that the separate crates are causing duplicated workflows or awkward cross-layer implementation. The current friction appears to be mostly ergonomic.

**Evidence:** Application command files import domain types and errors directly. `server_bootstrap` re-exports application input/error types for transports. HTTP and GraphQL crates both depend on `domain` for DTO conversion and error mapping, and on `server_bootstrap` for wired command/query access. These are visible seams, but they align with the documented adapter pattern.

**Impact:** Ergonomic cleanup may be worthwhile, but a full merge is a broad response to a narrow issue. Collapsing crates to reduce imports would trade away compile-time architectural enforcement for convenience.

**Recommended Fix:** Prefer smaller ergonomic improvements: targeted re-exports, clearer module names, or helper constructors. Do not merge crates solely to reduce import paths.

## Non-Issues / Things Checked

- I checked for obvious reverse dependencies from `domain` into `application`, `persistence`, `server_bootstrap`, `http_server`, `graphql_server`, `axum`, `async_graphql`, or `sqlx`; I found no such references in `server/crates/domain/src`.
- The presence of domain types in HTTP and GraphQL DTO/error mapping is acceptable for this project because transports map returned domain values and typed domain errors at the edge.
- `application` using domain repository ports and unit-of-work traits is acceptable; commands are the correct place for orchestration and transaction timing.
- `application::ports::gen_ident::IdentGeneratorPort` being application-owned is acceptable because member identifier generation is a use-case dependency wired in `server_bootstrap`, not a reusable domain invariant.
- `server_bootstrap` re-exporting application command/query/input types is acceptable as composition-root ergonomics and does not itself argue for a domain/application merge.

## Assumptions

- The proposed merge means a single Rust crate containing both current `server/crates/application` and `server/crates/domain` code, not merely a package rename or facade crate.
- The repository is intended to remain a clean/onion architecture + CQRS reference implementation where dependency direction is a first-class teaching and maintenance concern.
- Future features will increase application use-case complexity, making compile-time boundaries more valuable over time.
- I treated architectural risk from weakened boundaries as in scope even if the current code could still compile after a careful merge.

## Open Questions

- Is there a concrete build-time, publishing, or developer-experience problem caused by the two-crate split that cannot be solved with re-exports or tooling?
- Would the team be willing to add automated architecture checks if the crates are merged?
- Is there a future plan for alternate application layers over the same domain model, or is the current CQRS application layer expected to be the only use-case layer permanently?
- Should repository ports remain in `domain`, or does the team want to revisit whether some ports belong in `application` before making any crate-boundary decision?
