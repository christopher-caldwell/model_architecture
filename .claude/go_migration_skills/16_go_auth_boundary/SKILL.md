# Skill: Go Auth Boundary

## Purpose

Preserve the separated auth boundary from the Rust app.

## Package shape

Use `internal/auth` with files such as `claims.go`, `verifier.go`, `jwt.go`, `middleware.go`, and `context.go`.

## Core abstraction

Define a verifier interface:

```go
type Verifier interface {
    Verify(ctx context.Context, token string) (*Claims, error)
}
```

Define claims as auth-layer types. Adjust fields to the Rust/spec reference.

## JWT implementation

JWT libraries are used only inside auth implementation. JWT library types must not leak into domain/application. Bootstrap injects verifier config.

## HTTP integration

Middleware extracts bearer tokens, calls verifier, attaches claims to context, and rejects unauthenticated requests when required.

## GraphQL integration

GraphQL reads claims from request context and maps auth failures to GraphQL errors.

## Application boundary

Application receives only the identity/authorization information needed by a use case. If Rust does not pass actor identity into a command, do not invent it.

## Error mapping

Missing/invalid credentials map to 401 or `UNAUTHENTICATED`. Authenticated-but-forbidden maps to 403 or `FORBIDDEN`.

## Anti-patterns

Domain importing JWT, application parsing bearer tokens, handlers manually decoding JWT when verifier exists, auth verifier doing unrelated persistence, or JWT claims being used as domain entities.

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
