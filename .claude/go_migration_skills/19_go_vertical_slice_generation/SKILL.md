# Skill: Go Vertical Slice Generation

## Purpose

Add new features/use cases to the Go implementation after the initial migration without drifting edge-first.

## Required order

1. Update language-neutral behavior spec.
2. Add/update domain types and rules.
3. Add domain tests.
4. Add/update repository ports.
5. Add application command/query.
6. Add application tests with fakes.
7. Add persistence SQL/mapping.
8. Add persistence tests where appropriate.
9. Wire bootstrap dependencies.
10. Add HTTP route/DTO/handler.
11. Add GraphQL schema/resolver.
12. Update examples/docs.
13. Run architecture review.

## Mini-plan before code

State what behavior exists in Rust/spec, what domain rule is involved, whether it is command/query/both, which repos are needed, whether UOW is required, which tables/queries are touched, which HTTP route exposes it, which GraphQL operation exposes it, and which tests prove it.

## Domain-first rule

Any business behavior must be visible in domain code before transport code is added.

## Persistence rule

Persistence implements only what ports require. Do not add generic CRUD methods for hypothetical future use.

## Done criteria

A slice is done only when domain behavior and tests exist, application behavior and tests exist, persistence maps DB rows correctly, bootstrap wiring is complete, HTTP/GraphQL exposure is complete where required, examples/docs match, and architecture review passes.

## Anti-patterns

Route first development, generic repository bloat, GraphQL generated models defining domain shape, raw string statuses, DB changes without mapping/spec updates, or manual route testing replacing automated tests.

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
