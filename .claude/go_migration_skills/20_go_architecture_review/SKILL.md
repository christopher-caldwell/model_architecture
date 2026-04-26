# Skill: Go Architecture Review

## Purpose

Audit generated Go code against the Rust reference, language-neutral spec, and Go migration architecture.

## Inputs

Use the Rust inventory, language-neutral spec, current Go code, database schema/migrations, HTTP/GraphQL examples, tests, and recent diff.

## Phase 1: reference parity

Check same entities, statuses, commands, queries, rules, route contracts, GraphQL operations, persistence semantics, auth behavior, and error behavior. Separate missing behaviors from invented behaviors.

## Phase 2: architecture boundaries

Check domain has no persistence/transport/bootstrap/framework imports; application has no concrete persistence or transport imports; persistence implements ports and maps rows; HTTP/GraphQL are thin adapters; bootstrap wires concrete dependencies; auth remains separate.

## Phase 3: Go idiom

Check package names, no `index.go`, no dumping-ground packages, correct `context.Context` use, error wrapping/inspection, purposeful interfaces, safe transactions, and table-driven tests.

## Phase 4: persistence correctness

Check creates return real persisted rows, row structs do not leak, unknown statuses fail clearly, DB defaults/timestamps are not fabricated, UOW write repos share a transaction, read repos do not mutate, and SQL matches schema/spec.

## Phase 5: transport contract correctness

Check HTTP routes, request bodies, response DTOs, GraphQL schema, error mapping, and auth behavior.

## Phase 6: tests and examples

Check domain rule tests, application fake-UOW tests, persistence mapping tests, HTTP/GraphQL tests, examples, and conformance tests.

## Severity levels

Blocker: behavior or architecture no longer represents the reference app. High: code works but misleads future templates. Medium: quality/consistency issue. Low: style/readability polish.

## Report format

Produce verdict, blockers, high-priority issues, medium/low issues, what is solid, recommended fix order, and whether the Go implementation remains a faithful migration.

## Anti-patterns to flag

Literal Rust syntax translated into Go, generic tutorial architecture replacing this project, framework-first package layout, ORM active record domain, transport DTOs as domain models, SQL in application/domain, business logic in handlers/resolvers, missing UOW, and unreviewed route/schema drift.

## Non-negotiable guardrails

- Treat the Rust implementation as the behavioral reference, not syntax to copy.
- Preserve the same app, use cases, route contracts, GraphQL operations, persistence semantics, and architectural boundaries.
- Use idiomatic Go. Do not write Rust-in-Go clothing.
- Do not invent framework patterns, service locators, active-record models, generic `models` dumping grounds, or `index.go` files.
- Do not move business rules into HTTP handlers, GraphQL resolvers, SQL adapters, bootstrap wiring, generated code, or tests.
- Do not let transport DTOs, database rows, ORM models, SQL types, driver types, or generated GraphQL types cross inward into domain or application code.
- When a Rust construct does not translate directly, preserve the intent using normal Go mechanisms.
- Make uncertainty explicit. If the Rust reference does not contain a behavior, do not invent it.

## Required completion review

Before considering work from this skill complete, run two additional passes.

### Review pass 1: parity and boundary check

- Confirm the result still matches the Rust reference behavior.
- Confirm no dependency points inward incorrectly.
- Confirm no transport, persistence, framework, or generated-code concern leaked into domain/application.
- Confirm no generic Go tutorial pattern replaced the architecture of this project.
- Confirm every generated or changed artifact is justified by the Rust inventory or language-neutral spec.

### Review pass 2: Go idiom and drift check

- Confirm the Go code is idiomatic for package naming, errors, context, interfaces, transactions, and tests.
- Confirm the model did not create Rust-shaped Go abstractions where Go has a clearer native idiom.
- Confirm names are explicit and project-specific, not vague names like `models`, `common`, `manager`, `service2`, or `index`.
- Confirm all changed examples/docs/tests still agree with the generated code.
- Confirm unknowns were documented instead of guessed.
