# Skill: Go Bootstrap / Composition Root

## Purpose

Wire the Go app without breaking boundaries. Bootstrap is where concrete dependencies meet abstractions.

## Package location

Use `internal/bootstrap`. `cmd/http_server/main.go` and `cmd/graphql_server/main.go` should be small.

## Bootstrap responsibilities

Bootstrap may load config, create DB pools, run migrations if intentionally chosen, construct repos, construct the UOW factory, construct commands/queries, construct the auth verifier, create dependency structs, and return cleanup functions.

Bootstrap must not own business rules, route behavior, resolver behavior, row mapping, domain transitions, or HTTP/GraphQL error mapping.

## Dependency struct

Expose a dependency struct containing command groups, query groups, and auth verifier. Preserve the Rust grouping intent, such as catalog/lending/membership if present.

## Constructor pattern

Use explicit construction. Return deps, cleanup, and error. Do not use global singletons. Do not hide construction failures.

## Binaries

Each binary loads config, calls bootstrap, builds its transport server, starts the server, and handles shutdown.

## Config

Domain/application should not read environment variables. Persistence and auth receive concrete config values from bootstrap.

## Anti-patterns

Global DB pools, handlers constructing repos, resolvers constructing UOW factories, application constructors opening DB connections, bootstrap containing business conditionals, or service locator maps of `string` to `any`.

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
