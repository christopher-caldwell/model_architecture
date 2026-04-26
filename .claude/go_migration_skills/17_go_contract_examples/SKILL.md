# Skill: Go Contract Examples

## Purpose

Keep HTTP, GraphQL, README, and client examples aligned with the migrated Go app and Rust reference contract. This is a teaching/template project, so examples are part of the product.

## Maintain

Maintain HTTP request/response examples, Bruno collections or equivalent, GraphQL operations/variables/responses, route documentation, README snippets, curl examples, expected errors, and auth examples.

## Example rules

Examples must match actual DTOs. If `POST /books/{isbn}/copies` uses ISBN in the path, the JSON body must not include `isbn` unless specified. If `author_name` is derived from parent book, creation examples must not send it.

## Update procedure

When a use case changes, update domain/application, HTTP DTOs, GraphQL schema/DTOs, examples, docs, and conformance checks together.

## HTTP examples

For every route, include method, path, auth requirement, request body if any, success response, and common error response.

## GraphQL examples

For every operation, include operation document, variables, success response, and common error response.

## Consistency checks

Route names match code, body fields match DTOs, response fields match mappers, status values match domain/persistence mappings, identifiers match public contracts, and errors match transport mappings.

## Anti-patterns

Examples copied from an old Rust state, stale fields, GraphQL examples for missing operations, README routes that differ from code, or examples that imply nonexistent behavior.

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
