# Skill: Canonical Rust App Inventory

## Purpose

Use this before any Go migration work. The model must inspect the Rust reference and produce a concrete source-of-truth inventory. No Go code should be written during this skill.

The purpose is anti-hallucination: the model must not begin from a generic Go clean architecture template, a framework tutorial, or memory of a different project.

## When to use

- Starting the Go migration.
- Refreshing the migration after Rust changes.
- Adding a Go feature that already exists in Rust.
- Auditing whether the Go version has drifted.
- Creating or updating the language-neutral app spec.

## Required input sources

Inspect the active Rust implementation, including the workspace/crates, domain files, application commands/queries, repository ports, UOW, persistence adapters, SQL/migrations, HTTP routes, GraphQL schema/resolvers, auth boundary, README/docs, examples, and tests.

## Required output

Produce one inventory document with these sections:

1. Workspace/crate map.
2. Domain entity map.
3. Typed ID/value object map.
4. Status/enum map.
5. Domain rule and transition map.
6. Domain error map.
7. Repository port map.
8. Unit-of-work map.
9. Application command map.
10. Application query map.
11. Persistence table and SQL map.
12. HTTP route map.
13. GraphQL operation map.
14. Auth behavior map.
15. Error translation map.
16. Test/example coverage map.
17. Intentional design choices.
18. Unknowns requiring source inspection.

## Detail required for every inventory item

For every item, record the Rust file path, type/function/method name, layer, inputs, outputs, error behavior, dependencies, and whether Go must reproduce it exactly or reinterpret it idiomatically.

## Domain inventory requirements

Capture books, book copies, members, loans, IDs, identifiers, statuses, creation payloads, prepared structs, domain errors, pure rules, and transitions. For each rule, record the precondition, success result, and failure result.

## Application inventory requirements

For every command and query, capture the input type, output type, state loaded, domain methods called, repository methods called, transaction/UOW behavior, commit point, and errors returned.

## Persistence inventory requirements

Capture table names, DB row shapes, row-to-domain mappings, domain-to-SQL mappings, create/update/read SQL files, status/reference mapping, transaction-backed write repos, and pool-backed read repos.

## Transport inventory requirements

Capture HTTP and GraphQL separately. For each operation, record request/input DTO, response/output DTO, path/query/body fields, application command/query called, auth requirement, error mapping, and example/docs coverage.

## Prohibited behavior

Do not write Go code. Do not redesign the app. Do not invent missing use cases. Do not normalize terminology. Do not replace project-specific concepts with generic names.

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
