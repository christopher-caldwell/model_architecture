# Skill: Go Database Schema and Migration Parity

## Purpose

Preserve the Rust database contract when creating the Go version. The schema is part of the application behavior.

## Inputs

Use Rust DDL, migrations, seed/reference data, SQL query files, persistence mappings, and the language-neutral spec.

## Required inventory

Document schemas, tables, columns, primary keys, foreign keys, unique constraints, check constraints, indexes, status/reference tables, seed data, timestamp/default behavior, generated IDs/idents, nullable columns, and cascade/restrict behavior.

## Parity rules

The Go version must preserve table relationships, status/reference values, uniqueness, public identifiers, timestamp ownership, active loan semantics, copy status semantics, member status semantics, and route-visible identifiers.

## Migration tooling

Use a Go-friendly tool such as Goose, golang-migrate, Atlas, or a project SQL runner. The tool must not force domain/application to know about migrations.

## Seed data

Reference/status data must be deterministic and match Rust. Application/domain should not depend on magic database integer IDs. Prefer stable status codes/idents for mapping.

## SQL parity

For every Rust SQL query that implements behavior, identify the Go equivalent. The query may be rewritten, but must preserve filtering, joins, locking, result shape, ordering where observable, null handling, and status mapping.

## Generated values

If the database owns a value in Rust, it owns it in Go unless the spec intentionally changes ownership.

## Anti-patterns

Changing schema for ORM convenience, replacing lookup tables with app-only constants without a spec change, exposing DB IDs when Rust does not, or creating Go migrations that silently diverge.

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
