# Skill: Go Domain Modeling

## Purpose

Create Go domain entities and value types that preserve the Rust model while staying idiomatic.

## Package style

Use domain packages such as `internal/domain/book`, `internal/domain/bookcopy`, `internal/domain/member`, and `internal/domain/loan`.

## Typed IDs and identifiers

Use named Go types for meaningful IDs and public identifiers:

```go
type MemberID int64
type MemberIdent string
type BookID int64
type ISBN string
type BookCopyID int64
type Barcode string
type LoanID int64
type LoanIdent string
```

Use named types when the value has domain meaning. Keep DB primitive conversion in persistence and HTTP/GraphQL conversion at the edge.

## Status values

Use typed string constants, not raw strings everywhere:

```go
type MemberStatus string

const (
    MemberStatusActive MemberStatus = "active"
    MemberStatusSuspended MemberStatus = "suspended"
)
```

Provide parse/validate helpers when statuses come from persistence or transport. Unknown persisted statuses should return errors.

## Entities

Entities are plain structs with domain types. Do not put DB tags, GraphQL tags, or framework-specific annotations on domain entities. Avoid JSON tags in domain unless intentionally accepted by the project.

## Payloads and prepared structs

Use payload/prepared structs when the Rust domain uses them to separate caller input from domain-prepared state. Creation payloads should not include database-generated values or parent-derived read fields. For example, book copy creation should not include author name if it is derived from the parent book.

## Constructors/preparation

Use straightforward functions for domain defaults. They may return errors if domain preparation can fail. They must not call repositories.

## Time fields

Domain entities may carry timestamps as returned state. Persistence maps database timestamps into domain values. Application should not fabricate persisted timestamps.

## Anti-patterns

Raw string statuses, ORM models as domain entities, generated GraphQL types as domain entities, domain structs with DB tags, `any` for normal domain values, or a `models` package mixing rows and entities.

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
