# Skill: Go Repository Ports

## Purpose

Define repository interfaces that preserve the Rust architecture's intent while using Go idioms.

## Placement rule

Choose one consistent placement. Preferred: application-facing repository interfaces live in `internal/application` or `internal/application/ports`; pure domain abstractions live in domain only if they are truly domain concepts; concrete implementations live in persistence.

## Port design

Ports expose domain concepts, not persistence concepts:

```go
type MemberReadRepository interface {
    GetByIdent(ctx context.Context, ident member.MemberIdent) (*member.Member, error)
}
```

They must not import DB drivers, HTTP routers, GraphQL packages, or generated types.

## Read/write separation

Read repos are pool-backed and serve queries. Write repos are transaction-backed through UOW and serve commands. Write repos may include `ForUpdate` read methods where command consistency requires locking.

## Method naming

Use explicit names like `GetByIdent`, `GetByISBN`, `GetByBarcode`, `ListActiveByMemberID`, `FindActiveByBookCopyID`, `GetByIdentForUpdate`, `UpdateStatus`, and `Create`. Avoid vague `Save`, `Do`, `Exec`, or generic god-repository methods.

## Context rule

Every repo method accepts `context.Context` as the first parameter.

## Return rules

Return domain models, not DB rows. Return typed IDs/statuses after mapping. Return not-found errors consistently for missing single-object lookups. Return empty slices for successful empty lists.

## Transaction rule

Write repos exposed from UOW use the UOW transaction. They must not independently call `Begin`.

## Anti-patterns

Ports accepting request DTOs, returning row structs, importing DB/framework packages, exposing a raw transaction, or becoming a single all-purpose repository.

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
