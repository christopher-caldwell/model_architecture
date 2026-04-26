# Skill: Go HTTP Transport

## Purpose

Implement the HTTP adapter while preserving the Rust route contract and keeping handlers thin.

## Responsibilities

HTTP owns routes, DTOs, path/query/body extraction, transport validation, application calls, response mapping, HTTP error mapping, and auth middleware attachment. It does not own business rules, SQL, transactions, repos, or domain transitions.

## Route parity

Use the language-neutral spec and Rust reference. Preserve routes such as books, book copies, members, loans, member loans, and suspension routes exactly unless the spec is intentionally changed.

## DTOs

Request and response structs are transport DTOs. Do not reuse request DTOs as domain payloads without mapping. Do not put parent path fields redundantly in bodies unless specified. Do not include derived fields such as author name in create-copy bodies if Rust/spec derives them from the parent book.

## Handler flow

Handlers extract request data, map to application input, call command/query, map domain result to response DTO, and write response or mapped error.

## Error mapping

Map errors at the HTTP edge using `errors.Is`/`errors.As`: not found to 404, shape errors to 400, conflicts to 409 or spec-defined status, unauthenticated to 401, forbidden to 403, unexpected to 500.

## Router choice

Use standard `net/http` with Go 1.22+ patterns, `chi`, or another explicit project-approved router. Do not choose a framework that reshapes the architecture without approval.

## Anti-patterns

Handlers importing persistence, starting transactions, checking business rules directly, leaking request DTOs inward, adding JSON tags to domain entities just for handlers, or drifting route names.

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
