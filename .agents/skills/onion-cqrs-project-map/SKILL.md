# Onion CQRS Project Map

Use this skill when orienting to the project, adding a new feature, moving behavior between layers, or adapting this structure to a new application.

## Core Idea

This project demonstrates that the inner application model is independent of its outer delivery mechanisms. HTTP and GraphQL are interchangeable adapters over the same domain and application code.

Dependency direction:

```text
transport -> bootstrap -> application -> domain
                 persistence -> domain
```

Nothing in `domain` depends on transport, persistence, framework, database, runtime, or DI concerns.

## Layer Responsibilities

- `domain`: business concepts, invariants, state transitions, typed business errors, repository ports, unit-of-work traits.
- `application`: CQRS use cases. Commands orchestrate writes and transactions. Queries orchestrate reads. This is where multi-step workflows live.
- `persistence`: adapter that implements repository ports with SQLx, SQL files, read pools, and transactional write repositories.
- `server_bootstrap`: composition root. It wires pools, adapters, command/query structs, auth, and shared dependencies.
- `http_server`: Axum transport only.
- `graphql_server`: async-graphql transport only.
- `auth_core`: auth boundary and JWT adapter.
- `database/sql`: demo schema and seed data.

## Decision Rules

- Business policy belongs in `domain`.
- Use-case choreography belongs in `application`.
- IO implementation belongs in adapters.
- Request parsing and response mapping belong in transports.
- Dependency wiring belongs in bootstrap.

## Common File Pattern

Domain modules usually have:

- `entity.rs` for IDs, entities, creation payloads, prepared payloads.
- `enums.rs` for domain states.
- `errors.rs` for business failures.
- `logic.rs` for guards, transitions, and prepare methods.
- `port.rs` for read/write repository traits.

Application modules usually split by use-case area:

- `commands/*.rs` for write use cases.
- `queries/*.rs` for read use cases.
- `ports/*.rs` for application-owned external service ports.

Persistence mirrors domain concepts and separates SQL:

- `src/<concept>/read_repo.rs`
- `src/<concept>/write_repo.rs`
- `sql/<concept>/queries/*.sql`
- `sql/<concept>/commands/*.sql`

## When Adding A New Capability

1. Model the domain concept and invariants first.
2. Add domain errors for invalid business states.
3. Add or adjust ports.
4. Implement a command or query in `application`.
5. Implement persistence adapters and SQL.
6. Wire dependencies in `server_bootstrap`.
7. Add transport endpoints that call the application layer.
8. Add tests at the lowest useful layer. Domain logic tests should cover state transitions, guards, and preparation defaults.

## Portable Architecture Notes

This structure is intended to be adapted to other applications:

- Replace `http_server` and `graphql_server` with any transport.
- Replace `persistence` with any adapter that satisfies the same ports.
- Keep business rules in `domain`.
- Keep use-case workflow in `application`.
- Keep dependency wiring in `server_bootstrap` or the equivalent composition root.

## Review Checklist

- Could the same use case be exposed through HTTP and GraphQL with no duplicated business workflow?
- Does the domain make every reusable business decision?
- Are transport handlers free of query-then-command orchestration?
- Are command writes transactional and committed once?
- Is there any read after write used only to hydrate a response?
- Does SQL return only generated values required by the application?
- Are timestamps for returned created/updated entities produced with `Utc::now()` in code rather than selected back from the database?
- Are typed domain errors preserved until the transport-specific error mapper?

## Smells

- A handler or resolver branches on business state.
- HTTP and GraphQL implement the same workflow separately.
- A repository decides whether a business action is allowed.
- A command does a write and then selects the same row just to return it.
- A transport calls a query and then a command to perform one use case.
