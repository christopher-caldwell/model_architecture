---
name: transport-adapters
description: Use when changing HTTP handlers, GraphQL resolvers, schemas, request/response DTOs, auth extraction, router state, protocol-specific error mapping, or transport dependency wiring in this onion/CQRS Rust project.
---

# Transport Adapters

Use this skill when changing HTTP handlers, GraphQL resolvers, schemas, DTOs, auth extraction, response mapping, or protocol-specific error mapping.

For command transaction and UoW details, read `.agents/skills/unit-of-work-cqrs-pattern/SKILL.md`; transports should only receive already-wired commands and queries.

## Purpose

Transports are replaceable adapters. They expose the application but do not define the application.

HTTP and GraphQL should be able to expose the same use case by calling the same command or query.

## Handler And Resolver Shape

A transport endpoint should:

1. Extract auth/context.
2. Parse path, query string, body, or GraphQL input.
3. Build an application input struct when needed.
4. Call one command or one query for the use case.
5. Map the returned domain object into a transport DTO.
6. Map typed errors into protocol-specific errors.

Keep handlers and resolvers boring.

## What Does Not Belong Here

- Business state checks.
- Query-then-command orchestration.
- Transaction boundaries.
- SQL calls.
- Domain default assignment.
- Reusable workflow branching.
- Direct construction of persistence adapters.

## Error Mapping

Transport layers may map:

- `NotFound` domain errors to 404 or GraphQL not-found style errors.
- Invalid business transitions to 409 or GraphQL conflict style errors.
- Unexpected infrastructure errors to generic internal errors.

They should not invent business meaning. Preserve domain error intent.

## DTO Rules

- Request DTOs belong to the transport.
- Response DTOs belong to the transport.
- Convert domain types to DTOs at the edge.
- Do not leak Axum or async-graphql types into `application` or `domain`.

## Adding A New Transport

1. Define a transport-local `deps.rs` or equivalent composition root.
2. Wire pools, persistence adapters, commands, queries, auth, and transport state there.
3. Define transport-specific schemas.
4. Route each operation to an existing application command or query.
5. Add protocol-specific error mapping.
6. Avoid creating new business workflows in the transport.

## Smells

- Resolver loads a record, checks it, then calls a mutation command.
- Handler has more business branching than parse/map/error handling.
- HTTP and GraphQL have separate implementations of the same use case logic.
- Transport code imports SQLx or persistence repositories directly.
- Transport code assigns domain statuses or enforces loan limits.
