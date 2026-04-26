# Skill: Go Testing and Conformance

## Purpose

Test the Go migration in a way that proves it matches the Rust reference and language-neutral spec.

## Domain tests

Use table-driven tests for rules and transitions: suspended member cannot borrow, active member within limits can borrow, member over limit cannot borrow, copy status transitions, lost/maintenance/circulation behavior, and loan return behavior.

## Application tests

Use fake repos/UOW to test command orchestration: checkout happy path commits, suspended member fails without commit, loan limit failure does not write, create book copy uses path ISBN, report loss updates correct state, and return loan persists correct result.

## Persistence tests

Use integration tests when a test DB exists. Cover row mapping, status mapping, create returning actual persisted row, not-found behavior, constraint translation, and locking query behavior where applicable.

## HTTP tests

Test routes, request body shape, path parameter mapping, response mapping, error-to-status mapping, and absence of stale fields.

## GraphQL tests

Test operation existence, input mapping, output mapping, and error extension mapping.

## Conformance tests

Conformance tests prove behavior against the language-neutral spec. They should verify behavior, not Rust internals.

## Architecture tests

Consider static checks that domain does not import persistence/transport/bootstrap, application does not import concrete adapters, transport does not import persistence, no `index.go` files exist, and no generic `models` dumping ground exists.

## Fake UOW rules

Fakes should verify expected methods, commit on success, no commit on domain failure, rollback behavior if modeled, and saved values.

## Anti-patterns

Only testing handlers, testing domain rules only through HTTP, requiring a real DB for every command test, fakes reimplementing business logic, string-matching wrapped errors, or conformance tests blessing spec drift.

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
