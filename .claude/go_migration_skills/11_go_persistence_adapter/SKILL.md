# Skill: Go Persistence Adapter

## Purpose

Implement concrete Postgres persistence. Persistence adapts database rows and SQL operations to repository ports; it does not define the business model.

## Implementation choice

Lock one style before coding: either `pgx` with hand-written SQL or `sqlc` generated query methods with hand-written domain mapping. Do not mix casually. Unless instructed otherwise, prefer explicit `pgx` plus SQL for readability close to the Rust SQLx style.

## Package shape

Use packages under `internal/persistence/postgres`, with entity subpackages for rows, mappers, read repos, write repos, and SQL.

## Row structs

Rows are persistence-only, usually unexported, and never leave persistence packages.

## Mappers

Map rows into domain models explicitly. Validate status/reference values during mapping. Unknown DB status strings return errors; they do not silently default.

## Create operations

Create operations must return actual persisted rows using `RETURNING`. Do not reconstruct domain entities from input plus `time.Now()`. Let the database own generated IDs, timestamps, and DB defaults. Return joined fields when the domain return model requires them.

## Read repos

Read repos are pool-backed. They execute read SQL, scan rows, map rows to domain models, and return consistent not-found errors.

## Write repos

Write repos are transaction-backed. They use the UOW transaction, execute write SQL, return mapped domain models, and never commit/rollback/begin.

## SQL organization

Keep SQL out of domain/application. Use SQL constants or external `.sql` files consistently. Avoid dynamic SQL unless justified.

## Error handling

Wrap DB errors with operation context. Translate no-row and expected constraints. Preserve original causes with `%w`.

## Anti-patterns

ORM active record models as domain entities, `map[string]any` scans for normal queries, persistence returning HTTP codes, persistence accepting transport DTOs, or create methods returning non-persisted guesses.

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
