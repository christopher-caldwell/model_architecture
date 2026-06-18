# Independent Audit: reviewer-03

## Scope Reviewed
I reviewed the architecture decision of whether `server/crates/application` and `server/crates/domain` should be merged into one crate. I used the repository instructions, the Understand Anything knowledge graph, the required project skills, `audits/application-domain-merge-audit.config.yaml`, and representative files under the allowed crate paths. I did not review implementation details for an actual refactor and did not read any other reviewer reports.

## Executive Summary
I recommend keeping `application` and `domain` as separate crates for this project. The current split is not accidental ceremony: it encodes the clean/onion dependency rule in Cargo, lets persistence implement domain ports without depending on use-case orchestration, and lets transports map domain values without gaining direct access to application internals. Merging is possible in Rust if strict module conventions are retained, but it would trade a hard compile-time package boundary for softer discipline in a codebase whose stated purpose is to demonstrate clean/onion architecture and CQRS.

## Findings

### Finding 1: The Separate Crates Encode The Primary Inward Dependency Rule
**Severity:** high
**Confidence:** 9
**Category:** design
**Location:** `server/Cargo.toml:1-10`, `server/crates/application/Cargo.toml:6-13`, `server/crates/domain/Cargo.toml:6-10`, `server/crates/persistence/Cargo.toml:6-12`
**Status:** recommendation

**Issue:** Merging `application` and `domain` would remove the strongest mechanical boundary that currently enforces the project's intended dependency direction.

**Evidence:** The workspace declares `application` and `domain` as distinct packages in `server/Cargo.toml:1-10`. `application` depends on `domain` at `server/crates/application/Cargo.toml:6-13`; `domain` has no dependency back to `application` in `server/crates/domain/Cargo.toml:6-10`; and `persistence` depends on `domain`, not `application`, in `server/crates/persistence/Cargo.toml:6-12`. The Understand Anything graph also classifies the layers separately: Domain is "Business entities, invariants, typed domain errors, state transitions, and repository ports"; Application CQRS is "Commands and queries that orchestrate use cases, transactions, ports, and domain decisions."

**Impact:** The current package graph makes the onion rule visible and compiler-checkable at the crate level. If the crates are merged, Rust can still preserve modules, but it no longer prevents use-case code, domain code, and ports from growing cross-references inside one package. That is a meaningful loss for a reference implementation whose architecture is the product.

**Recommended Fix:** Do not merge by default. If the team proceeds anyway, keep `domain` and `application` as explicit top-level modules in the merged crate, add architecture tests or dependency-lint checks, and keep public exports narrow enough that domain modules cannot import application modules.

### Finding 2: Persistence Currently Depends On Domain Ports Without Depending On Use Cases
**Severity:** high
**Confidence:** 9
**Category:** design
**Location:** `server/crates/domain/src/uow.rs:8-50`, `server/crates/domain/src/member/port.rs:5-19`, `server/crates/persistence/src/uow.rs:3-14`, `server/crates/persistence/src/uow.rs:35-148`
**Status:** risk

**Issue:** A merge would blur the important distinction between domain repository contracts and application command/query orchestration.

**Evidence:** Domain owns repository and unit-of-work abstractions: `UnitOfWorkPort` and `WriteUnitOfWorkFactory` are defined in `server/crates/domain/src/uow.rs:8-50`, and member read/write ports are defined in `server/crates/domain/src/member/port.rs:5-19`. Persistence implements those domain ports directly in `server/crates/persistence/src/uow.rs:35-148`. The persistence crate imports only `domain` abstractions and SQLx infrastructure in `server/crates/persistence/src/uow.rs:3-14`.

**Impact:** This is a clean adapter relationship: persistence satisfies the core's ports without depending on application use cases. If the crates are merged, persistence either depends on a broader "core" crate that also contains command/query orchestration, or the refactor must introduce extra module boundaries to recreate the same separation. The former weakens adapter replaceability; the latter largely gives back the simplicity benefit of merging.

**Recommended Fix:** Keep repository ports and unit-of-work traits in the domain crate. If reducing crate count is mandatory, document and enforce that persistence may import only the merged crate's domain-facing modules, not command/query modules.

### Finding 3: Domain Contains Real Business Rules, Not Just Shared DTOs
**Severity:** medium
**Confidence:** 9
**Category:** maintainability
**Location:** `server/crates/domain/src/member/logic.rs:5-58`, `server/crates/domain/src/book_copy/logic.rs:5-63`, `server/crates/domain/src/loan/logic.rs:4-27`, `server/crates/domain/src/member/errors.rs:1-13`
**Status:** observation

**Issue:** Treating the merge as a harmless packaging simplification understates the role of the domain crate.

**Evidence:** Member transitions and guards live in `server/crates/domain/src/member/logic.rs:5-58`, including `suspend`, `reactivate`, `ensure_can_borrow`, `ensure_within_loan_limit`, and creation preparation. Book copy borrowability and state transitions live in `server/crates/domain/src/book_copy/logic.rs:5-63`. Loan returnability and preparation live in `server/crates/domain/src/loan/logic.rs:4-27`. Typed business errors are separate domain types, for example `MemberError` in `server/crates/domain/src/member/errors.rs:1-13`.

**Impact:** The domain crate provides a small but meaningful business language that application commands consume. Keeping it separate reinforces that reusable decisions belong in domain methods while commands orchestrate workflows. Merging increases the chance that new contributors place business conditions directly in commands because the code is physically closer and no crate boundary signals a conceptual transition.

**Recommended Fix:** Preserve the domain crate. If merged, keep domain modules independently testable and document the rule that commands must call domain guards/transitions rather than inspect state directly.

### Finding 4: Application Commands Correctly Orchestrate Domain Decisions And Transactions
**Severity:** informational
**Confidence:** 8
**Category:** design
**Location:** `server/crates/application/src/commands/lending.rs:59-130`, `server/crates/application/src/commands/catalog.rs:88-182`, `server/crates/application/src/commands/membership.rs:62-112`
**Status:** observation

**Issue:** The current boundary aligns with the actual code shape: application is a workflow layer, not a duplicate domain layer.

**Evidence:** `check_out_book_copy` builds a unit of work, loads the member and copy, calls domain guards, counts active loans, creates a prepared loan, and commits once in `server/crates/application/src/commands/lending.rs:59-103`. Return and lost-copy flows similarly coordinate domain guards, writes, commits, and code-side timestamps in `server/crates/application/src/commands/lending.rs:106-130` and `server/crates/application/src/commands/lending.rs:132-163`. Catalog and membership commands call domain transitions before update writes in `server/crates/application/src/commands/catalog.rs:88-182` and `server/crates/application/src/commands/membership.rs:62-112`.

**Impact:** This separation is currently doing useful explanatory work. The code demonstrates "domain decides, application orchestrates" cleanly. Merging would not improve correctness here; it would mostly reduce package count while making the distinction less obvious.

**Recommended Fix:** Keep the current layering. If merged, preserve the command/query module structure and add review checklist items that reject business-rule branching in command code unless it delegates to a domain guard or transition.

### Finding 5: A Merge Would Broaden The Dependency Surface Seen By Transports
**Severity:** medium
**Confidence:** 8
**Category:** maintainability
**Location:** `server/crates/http_server/Cargo.toml:6-23`, `server/crates/graphql_server/Cargo.toml:6-20`, `server/crates/http_server/src/router/books/schemas.rs:1-46`, `server/crates/graphql_server/src/router/graphql/catalog/mod.rs:1-92`
**Status:** risk

**Issue:** HTTP and GraphQL currently depend on `domain` for DTO conversion and error mapping while using `server_bootstrap` to reach application use cases. Merging application and domain would make that direct domain dependency a direct dependency on a broader crate containing commands and queries.

**Evidence:** HTTP depends on both `domain` and `server_bootstrap` in `server/crates/http_server/Cargo.toml:6-23`; GraphQL does the same in `server/crates/graphql_server/Cargo.toml:6-20`. HTTP schema mapping imports domain `Book` and `BookCreationPayload` in `server/crates/http_server/src/router/books/schemas.rs:1-46`. GraphQL catalog mapping imports domain `Book`, `BookCreationPayload`, `BookCopy`, and `BookCopyStatus` in `server/crates/graphql_server/src/router/graphql/catalog/mod.rs:1-92`. Handlers/resolvers call commands through bootstrap dependencies, for example `server/crates/http_server/src/router/books/post_handlers.rs:28-79` and `server/crates/graphql_server/src/router/graphql/catalog/mutations.rs:14-48`.

**Impact:** The current setup lets transports see domain values for edge mapping without bypassing the composition root for use-case execution. After a merge, transports would have easier access to application command/query internals unless exports are carefully constrained. That increases the risk of future transport-level workflow composition, which the project rules explicitly reject.

**Recommended Fix:** Keep the split, or introduce transport-facing DTO/read-model exports so transports no longer need a direct dependency on domain. If merged, avoid glob re-exports from the merged crate and keep command constructors inaccessible except through bootstrap.

### Finding 6: The Main Argument For Merging Is Ergonomic, Not Architectural
**Severity:** low
**Confidence:** 8
**Category:** maintainability
**Location:** `server/crates/application/src/lib.rs:1-3`, `server/crates/domain/src/lib.rs:1-5`, `server/crates/application/src/ports/gen_ident.rs:1-3`, `server/crates/application/src/ports/services.rs:1`
**Status:** observation

**Issue:** There is some ceremony in the split, but I did not find evidence that it is causing architectural harm.

**Evidence:** The top-level lib files are tiny: `application` exports commands, ports, and queries in `server/crates/application/src/lib.rs:1-3`; `domain` exports business modules and `uow` in `server/crates/domain/src/lib.rs:1-5`. Application-owned ports are also small, with `IdentGeneratorPort` in `server/crates/application/src/ports/gen_ident.rs:1-3` and an empty `services.rs` at `server/crates/application/src/ports/services.rs:1`.

**Impact:** A merge could reduce manifests and imports, and for a small application that might be a reasonable production tradeoff. In this repository, however, the extra crate is cheap and supports the architecture demonstration. The ergonomic benefit does not outweigh the loss of a clear package-level boundary.

**Recommended Fix:** Leave the crates separate. Clean up minor ceremony instead: remove or fill the empty application services module, keep exports intentional, and document why transport crates depend on domain directly.

## Non-Issues / Things Checked
- I checked for obvious outer-layer dependencies in `domain` and did not find imports of SQLx, Axum, async-graphql, persistence, bootstrap, or transports.
- I checked for obvious framework leakage in `application` and did not find SQLx, Axum, async-graphql, persistence, bootstrap, or transport imports.
- Direct transport imports of domain types looked acceptable for DTO conversion and error mapping, not business workflow ownership.
- Persistence write methods follow the narrow write-result pattern in representative SQL: inserts return generated IDs or generated identifiers only, and updates do not return hydrated rows.
- Application commands use domain guards/transitions for the main business rules and commit once at the end of representative write workflows.

## Assumptions
- The repository is intended to remain a clean/onion architecture + CQRS reference implementation, so explicit architecture boundaries are more valuable than minimizing crate count.
- The audit target is the decision itself, not a specific proposed merged-crate design with compensating lint rules.
- Build performance and publish/package management are not the primary drivers for this decision.
- The current direct transport dependency on `domain` is intentional edge mapping, not a sign that domain should be folded into application.

## Open Questions
- Is the project optimizing for teaching/reference clarity, or for a production service where fewer crates may be preferred?
- Would the team add automated module-boundary checks if the crates were merged?
- Are transport DTOs expected to continue converting directly from domain entities, or is there a future plan for application-owned read models?
- Should repository ports remain in domain long term, or does the team intend to move some ports to application as use-case-specific contracts?
