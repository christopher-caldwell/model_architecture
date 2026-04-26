# Skill: Go Error Semantics

## Purpose

Translate Rust-style typed errors into idiomatic Go errors without losing architectural meaning.

## Domain errors

Use sentinel errors when callers only need identity:

```go
var ErrMemberSuspended = errors.New("member is suspended")
```

Use custom error types when structured fields are needed. Keep errors inspectable with `errors.Is` and `errors.As`.

## Application errors

Application errors represent use-case failures such as not found, conflict, invalid command, wrapped domain rule failure, or wrapped infrastructure failure. Wrap with `%w`.

```go
return nil, fmt.Errorf("checkout book copy: %w", err)
```

## Persistence errors

Persistence translates expected DB outcomes: no rows to not-found, unique constraint to conflict, unknown status to invalid status, DB unavailable to wrapped infrastructure error. Preserve the original cause.

## Transport errors

HTTP and GraphQL convert application/domain errors into edge-specific responses. HTTP status codes and GraphQL extension codes must not flow inward.

## Mapping guidance

Not found usually maps to HTTP 404 and GraphQL `NOT_FOUND`. Shape validation maps to 400 / `BAD_USER_INPUT`. Auth missing maps to 401 / `UNAUTHENTICATED`. Forbidden maps to 403 / `FORBIDDEN`. Unexpected infrastructure maps to 500 / `INTERNAL`.

## Prohibited patterns

Do not simulate Rust error enums awkwardly. Do not panic for normal business failures. Do not compare `err.Error()`. Do not lose causes by formatting without `%w`. Do not expose raw SQL errors to clients.

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
