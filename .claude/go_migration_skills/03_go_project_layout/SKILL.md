# Skill: Go Project Layout

## Purpose

Create or review the Go project structure for the migrated app. The structure must express the same modular architecture as the Rust app while using idiomatic Go package organization.

## Default layout

Use this default unless the repo explicitly chooses another equivalent shape:

```text
servers/go/
  go.mod
  go.sum
  cmd/
    http_server/
      main.go
    graphql_server/
      main.go
  internal/
    domain/
      book/
      bookcopy/
      member/
      loan/
    application/
      commands/
      queries/
    persistence/
      postgres/
        book/
        bookcopy/
        member/
        loan/
    bootstrap/
    auth/
    transport/
      http/
      graphql/
  migrations/
  sql/
  examples/
  docs/
```

## Rust-to-Go boundary map

- Rust `domain` crate -> Go `internal/domain/...` packages.
- Rust `application` crate -> Go `internal/application/...` packages.
- Rust repo ports -> Go interfaces at the chosen inward boundary.
- Rust `persistence` crate -> Go `internal/persistence/postgres/...`.
- Rust `server_bootstrap` crate -> Go `internal/bootstrap`.
- Rust `http_server` crate -> Go `cmd/http_server` plus `internal/transport/http`.
- Rust `graphql_server` crate -> Go `cmd/graphql_server` plus `internal/transport/graphql`.
- Rust `auth_core` crate -> Go `internal/auth`.

## Naming rules

Use `main.go` only under `cmd/<binary>/`. Do not use `index.go`. Do not create generic `models`, `entities`, `services`, or `utils` dumping grounds. Prefer short lowercase package names that represent domain concepts or architectural roles. Use `bookcopy` as a Go package name, not `book_copy`.

## File organization guidance

Domain packages may use `entity.go`, `ids.go`, `status.go`, `payload.go`, `errors.go`, `logic.go`, and `logic_test.go`.

Persistence packages may use `rows.go`, `mapper.go`, `read_repo.go`, `write_repo.go`, and `sql.go`.

Transport packages may separate `routes.go`, `handlers.go`, `schemas.go` or `dto.go`, `mapper.go`, and `errors.go`.

## Dependency rules

Domain imports only standard library and sibling domain packages when necessary. Application may import domain and port interfaces. Persistence may import domain and port-owning packages. Transports may import application, domain for mapping, auth, and bootstrap dependency types. Bootstrap may import concrete adapters to wire them.

## Layout review questions

Can HTTP be removed without changing domain/application/persistence? Can GraphQL be removed without changing domain/application/persistence? Can Postgres be replaced without changing domain/application? Is any package becoming a dumping ground?

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
