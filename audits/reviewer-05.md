# Independent Audit: reviewer-05

## Scope Reviewed
I reviewed the architecture decision of whether to merge `server/crates/application` and `server/crates/domain` into one crate. Materials reviewed were the repository AGENTS instructions, `.understand-anything/knowledge-graph.json`, the required `.agents/skills/*` architecture instructions, `audits/application-domain-merge-audit.config.yaml`, the workspace and relevant crate `Cargo.toml` files, and representative files under `server/crates/application`, `server/crates/domain`, `server/crates/persistence`, `server/crates/server_bootstrap`, `server/crates/http_server`, and `server/crates/graphql_server`.

I did not read any `audits/reviewer-*.md` file or `audits/final-consensus-audit.md`.

## Executive Summary
Recommendation: do not merge `application` and `domain` as the default architecture decision for this reference implementation.

The current two-crate split is doing useful architectural work. It physically enforces the inward dependency direction: `application` depends on `domain`; persistence implements domain ports without seeing application commands; bootstrap wires both; transports call use cases and map domain values/errors at the edge. A merge would not obviously improve correctness, performance, or ergonomics enough to justify removing that compile-time boundary. If a merge is still desired for packaging simplicity, it should be treated as a single `core` crate with explicit `domain` and `application` modules, facade exports, dependency linting, and tests that preserve the current layering rules.

## Findings

### Finding 1: Merging removes a valuable compile-time boundary between domain policy and use-case orchestration
**Severity:** high
**Confidence:** 9
**Category:** design
**Location:** `README.md` lines 3-14 and 20-34; `server/Cargo.toml` lines 1-9; `server/crates/application/Cargo.toml` lines 1-14; `server/crates/domain/Cargo.toml` lines 1-11; `.understand-anything/knowledge-graph.json` layer summaries
**Status:** recommendation

**Issue:** The current crate split encodes the repository's clean/onion architecture in Cargo itself. `domain` is a separate inner crate, and `application` is a separate CQRS orchestration crate that depends inward on it. Merging them would turn this from a compiler-enforced dependency rule into a convention.

**Evidence:** The README states that the project exists to keep core logic separated from transports and infrastructure and to allow multiple interfaces over the same core (`README.md` lines 3-14). It separately defines Domain as business models/rules and Application as commands/queries and required interfaces (`README.md` lines 20-27). The workspace currently lists `crates/application` and `crates/domain` as separate members (`server/Cargo.toml` lines 1-9). `application` depends on `domain` (`server/crates/application/Cargo.toml` line 7), while `domain` does not depend on `application`, persistence, bootstrap, HTTP, or GraphQL (`server/crates/domain/Cargo.toml` lines 1-11). The Understand Anything graph also classifies them as distinct layers: Domain for entities/invariants/errors/ports and Application CQRS for use-case orchestration.

**Impact:** The merge would make it easier for future changes to put workflow logic beside business invariants without any crate-level signal that a boundary has been crossed. In a reference implementation, the physical separation is part of the lesson: reusable business decisions live in `domain`, while multi-step use cases live in `application`. Removing that separation increases the maintenance burden and makes architecture regressions more likely.

**Recommended Fix:** Keep `domain` and `application` as separate crates. If packaging pressure requires a single crate, create an explicit `core` crate with `core::domain` and `core::application` modules, narrow public facades, dependency linting or architecture tests, and documentation that the merge is a packaging choice rather than a layering change.

### Finding 2: A merged crate would broaden persistence's dependency surface from domain ports to application use cases
**Severity:** medium
**Confidence:** 9
**Category:** maintainability
**Location:** `server/crates/persistence/Cargo.toml` lines 1-14; `server/crates/persistence/src/uow.rs` lines 3-14 and 35-148; `server/crates/domain/src/uow.rs` lines 8-50; `server/crates/server_bootstrap/src/deps.rs` lines 4-15 and 68-111
**Status:** risk

**Issue:** Persistence currently depends only on `domain` so it can implement repository and unit-of-work ports. If `domain` is folded into `application`, persistence would need to depend on the merged crate to see those ports and entities. That gives the adapter visibility to application commands, queries, command inputs, and application-owned ports unless the merged crate is carefully hidden behind module visibility.

**Evidence:** `persistence` depends on `domain` but not on `application` (`server/crates/persistence/Cargo.toml` line 7). Its SQL unit of work imports domain write repository ports and entities (`server/crates/persistence/src/uow.rs` lines 3-14) and implements those domain ports (`server/crates/persistence/src/uow.rs` lines 35-148). The unit-of-work abstraction lives in `domain` (`server/crates/domain/src/uow.rs` lines 8-50). Bootstrap is where `application` and `persistence` are composed together (`server/crates/server_bootstrap/src/deps.rs` lines 4-15 and 68-111).

**Impact:** The current structure prevents persistence from depending on or calling application workflows. A merged crate weakens that adapter boundary and makes accidental coupling easier, especially as the application command/query surface grows. It can also cause persistence to rebuild due to application-only changes that should not affect domain port contracts.

**Recommended Fix:** Preserve the current dependency shape: `persistence -> domain`, `application -> domain`, and `server_bootstrap -> application + persistence`. If merging proceeds, keep ports/entities in a clearly isolated public `domain` module and make application command/query modules inaccessible to persistence by visibility, facade exports, and review/lint rules.

### Finding 3: Transport edge mapping already relies on domain as a stable public contract; merging would blur that contract
**Severity:** medium
**Confidence:** 8
**Category:** design
**Location:** `server/crates/http_server/Cargo.toml` lines 1-25; `server/crates/graphql_server/Cargo.toml` lines 1-22; `server/crates/http_server/src/router/members/schemas.rs` lines 1-48; `server/crates/http_server/src/router/errors.rs` lines 1-96; `server/crates/graphql_server/src/router/graphql/membership/mod.rs` lines 1-62; `server/crates/graphql_server/src/router/graphql/mod.rs` lines 67-131
**Status:** risk

**Issue:** HTTP and GraphQL currently depend on `domain` for DTO conversion and protocol-specific error mapping while using `server_bootstrap` to reach application commands/queries. If `domain` and `application` become one crate, transports either depend directly on a much larger core crate or bootstrap must re-export more domain surface. Both options make the edge contract less crisp.

**Evidence:** Both transport manifests depend on `domain` and `server_bootstrap` (`server/crates/http_server/Cargo.toml` lines 7-8; `server/crates/graphql_server/Cargo.toml` lines 7-8). HTTP schema mapping converts `domain::member::Member` and `MemberCreationPayload` at the edge (`server/crates/http_server/src/router/members/schemas.rs` lines 1-48). HTTP error mapping matches domain errors from `CommandError` into status codes (`server/crates/http_server/src/router/errors.rs` lines 52-96). GraphQL similarly converts domain `Member` and `MemberStatus` into GraphQL objects/enums (`server/crates/graphql_server/src/router/graphql/membership/mod.rs` lines 1-62) and maps domain errors into GraphQL extension codes (`server/crates/graphql_server/src/router/graphql/mod.rs` lines 67-131).

**Impact:** The existing direct transport-to-domain dependency is acceptable because it is used for parse/map/error-boundary work. After a merge, that same dependency would also expose application internals to transports unless constrained. This increases the chance that future handlers or resolvers compose workflows directly instead of calling one command/query.

**Recommended Fix:** Keep `domain` as the stable edge mapping contract. If merging is chosen, require transports to depend on a narrow facade and prohibit direct use of application internals outside `server_bootstrap`-provided command/query handles.

### Finding 4: The observed benefits of merging appear small relative to the architectural cost
**Severity:** low
**Confidence:** 8
**Category:** maintainability
**Location:** `server/crates/domain/src/lib.rs` lines 1-5; `server/crates/application/src/lib.rs` lines 1-3; `server/crates/application/src/commands/mod.rs` lines 1-23; `server/crates/application/src/queries/catalog.rs` lines 1-46; `server/crates/application/src/commands/lending.rs` lines 59-163
**Status:** observation

**Issue:** The current split does not show enough friction to justify collapsing it. The crates are small, their responsibilities are explicit, and the application layer is not duplicating domain logic at a scale that would make the boundary costly.

**Evidence:** `domain` exports only the core modules (`server/crates/domain/src/lib.rs` lines 1-5). `application` exports commands, ports, and queries (`server/crates/application/src/lib.rs` lines 1-3) with small command input structs (`server/crates/application/src/commands/mod.rs` lines 6-18). Queries are thin repository orchestration over domain read ports (`server/crates/application/src/queries/catalog.rs` lines 1-46). Commands contain the expected workflow: build a unit of work, load rows for update, call domain guards/transitions, write, commit, and shape the return value (`server/crates/application/src/commands/lending.rs` lines 59-163).

**Impact:** A merge would mainly reduce a Cargo package and import line count. In exchange, it would remove a clear teaching and maintenance boundary in a repository whose stated purpose is demonstrating clean, replaceable architecture.

**Recommended Fix:** Do not merge for convenience alone. Address isolated friction through re-exports, clearer naming, module documentation, or architecture tests rather than collapsing the crates.

### Finding 5: Domain-only tests are a useful signal that would be diluted by a merged crate unless preserved deliberately
**Severity:** low
**Confidence:** 7
**Category:** testing
**Location:** `server/crates/domain/src/member/logic.rs` lines 61-153; `server/crates/domain/src/book_copy/logic.rs` lines 65-208; `server/crates/domain/src/loan/logic.rs` lines 29-82; repository test scan for `server/crates/application`
**Status:** observation

**Issue:** Domain tests currently live next to domain business decisions and can be run as tests for the domain crate. A merged crate would not necessarily break those tests, but it would reduce the visible separation between cheap invariant tests and heavier application orchestration tests.

**Evidence:** Member logic tests cover transitions, guards, loan limits, and preparation defaults (`server/crates/domain/src/member/logic.rs` lines 61-153). Book copy logic tests cover borrowability, maintenance/lost transitions, and prepare defaults (`server/crates/domain/src/book_copy/logic.rs` lines 65-208). Loan logic tests cover returnability and preparation (`server/crates/domain/src/loan/logic.rs` lines 29-82). A scan for test attributes found tests in domain modules but no corresponding test attributes under `server/crates/application`.

**Impact:** The current testing layout reinforces that domain rules are independently testable. A merge could make it easier for tests to drift toward use-case-level coverage while missing isolated domain invariant coverage.

**Recommended Fix:** Keep the domain crate and its focused tests. If merging proceeds, preserve a `domain` test module structure and add architecture tests that fail when domain modules import application, persistence, transport, or framework concerns.

## Non-Issues / Things Checked
- Direct `domain` imports in HTTP and GraphQL looked suspicious at first, but in the reviewed files they are used for DTO conversion and typed error mapping, which is consistent with the transport-adapter instructions.
- The `domain` crate contains `async-trait` and `anyhow` because it owns repository and unit-of-work port traits. That is not a reason by itself to merge; the project instructions explicitly allow domain repository ports and unit-of-work traits.
- Application commands read during write workflows, but the sampled reads are through write repositories and `FOR UPDATE`-style methods before writes, not post-write hydration.
- Persistence create methods construct returned domain entities from generated IDs, prepared input, and code-side `Utc::now()`, which supports the current boundary rather than arguing for a merge.
- The presence of small command input structs in `application` does not create meaningful duplication that would justify merging.

## Assumptions
- The audit target is the architecture decision only, not a request to implement or review a completed refactor.
- The repository is intended to remain a clean/onion + CQRS reference implementation where architectural clarity is an explicit goal.
- Future contributors will benefit from compiler-enforced crate boundaries more than from a slightly smaller workspace.
- If a merge is proposed, it would combine the public surfaces of `application` and `domain` unless additional module visibility and facade work is explicitly included.

## Open Questions
- What concrete pain is motivating the proposed merge: compile time, dependency management, naming, publishability, IDE ergonomics, or conceptual simplicity?
- Would a single `core` crate with strict `domain` and `application` modules satisfy the motivation while retaining enforceable boundaries?
- Should transports continue importing domain types directly, or should `server_bootstrap` or a dedicated API facade own the edge-facing contract?
- Should the project add architecture tests or lint rules to enforce dependency direction regardless of the merge decision?
