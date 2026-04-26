# Skill: Go Boundary Ownership

## Purpose

Use this whenever deciding where a type, function, interface, error, or file belongs.

## Domain owns

Business entities, typed IDs, statuses, creation payloads when they represent domain defaults, prepared structs, business rules, state transitions, domain-specific errors, and pure rule tests.

Domain must not know SQL, Postgres, transactions, HTTP, GraphQL, JSON request bodies, middleware, JWT libraries, config, dependency injection, logging frameworks, or generated code.

## Application owns

Commands, queries, command/query input structs, use-case orchestration, transaction boundaries through UOW, calls to domain rules, repository coordination, application-level errors, and not-found/conflict decisions.

Application must not know concrete Postgres adapters, HTTP status codes, GraphQL error extensions, framework request/response objects, generated GraphQL types, or DB row structs.

## Persistence owns

Concrete Postgres repositories, DB row structs, SQL execution, SQL organization, transaction-backed write repos, pool-backed read repos, row-to-domain mapping, and domain-to-DB value mapping.

Persistence must not own business decisions, route DTOs, GraphQL DTOs, use-case orchestration, or auth decisions.

## HTTP owns

Routes, method/path definitions, request DTOs, response DTOs, path/query/body extraction, transport validation, application calls, response mapping, and HTTP error mapping.

## GraphQL owns

Schema/resolver DTOs, input/output conversion, resolver methods, auth context extraction, application calls, and GraphQL error mapping.

## Bootstrap owns

Config loading, DB pool creation, JWT verifier construction, repo construction, UOW factory construction, command/query construction, and dependency structs passed to transports.

## Auth owns

Verifier interfaces, JWT verification implementation, claims/current-user types, middleware/resolver helpers, and auth-specific errors.

## Decision procedure

Ask: Is it a business rule? Domain. Is it orchestration? Application. Is it SQL or row mapping? Persistence. Is it request/response mapping? Transport. Is it concrete dependency construction? Bootstrap. Is it token verification or claims? Auth.

## Red flags

Handlers checking member suspension, resolvers starting transactions, domain importing DB/HTTP/GraphQL packages, persistence accepting request DTOs, commands returning HTTP responses, GraphQL generated types in application/domain, or bootstrap containing business conditionals.

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
