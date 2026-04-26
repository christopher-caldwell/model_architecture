# Skill: Go CQRS Application

## Purpose

Implement Go application commands and queries in the same architectural image as the Rust app.

## Commands

Commands are write use cases. A command receives an input struct, starts a UOW, loads required state, calls domain rules/transitions, prepares new domain values, persists through write repos, commits once, returns domain models, and wraps errors with use-case context.

Application input structs are not HTTP bodies or GraphQL inputs. Transports map DTOs into application inputs.

## Queries

Queries are read use cases. A query receives typed input, calls read repo ports, returns domain models or slices, avoids UOW unless a consistent read transaction is explicitly needed, and avoids transport DTOs.

## Application services

Use grouped command/query structs that mirror the Rust app's use-case groups, such as catalog, lending, and membership, if those are still the Rust reference groups.

## Error behavior

Application may translate lower-level errors into not-found, conflict, invalid command, or wrapped infrastructure errors. It must not return HTTP responses or GraphQL extension codes.

## Dependency rules

Application may depend on domain, standard library, and port interfaces. It must not depend on concrete persistence, HTTP routers, GraphQL packages, DB drivers, or bootstrap.

## Transaction rules

All multi-step writes use UOW. Queries do not mutate. Commands should not call other commands unless the project explicitly adopts command composition.

## Anti-patterns

Handlers performing half the command, commands accepting request bodies, queries returning generated GraphQL types, application opening DB connections, or application duplicating domain rules as ad-hoc conditionals.

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
