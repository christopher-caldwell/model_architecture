# Skill: Go Unit of Work

## Purpose

Preserve the Rust write transaction model in idiomatic Go.

## Core shape

Use a UOW factory and UOW interface:

```go
type UnitOfWorkFactory interface {
    New(ctx context.Context) (UnitOfWork, error)
}

type UnitOfWork interface {
    Books() BookWriteRepository
    BookCopies() BookCopyWriteRepository
    Members() MemberWriteRepository
    Loans() LoanWriteRepository
    Commit(ctx context.Context) error
    Rollback(ctx context.Context) error
}
```

Adjust names to the final package layout.

## Command usage pattern

Commands create a UOW, defer rollback, use transaction-backed repos, commit once at the end, and return the result.

Rollback must be safe after commit. Repositories do not commit, rollback, or begin independent transactions.

## Read vs write

Queries use pool-backed read repos. Commands use UOW write repos. If a command reads state it later relies on, use UOW repos and explicit locking methods when needed.

## Row locking

Use explicit names like `GetByIdentForUpdate` rather than hiding locks behind ordinary reads.

## Error behavior

Return the original operation error when a pre-commit operation fails. Return commit errors with context. Do not panic on rollback failure.

## Anti-patterns

Repos starting their own transactions, commands mutating through pool-backed repos, UOW living in handlers/resolvers, application code receiving raw transaction objects, or committing before all domain checks complete.

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
