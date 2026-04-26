# Node / TypeScript Server Implementation Plan

This folder is the implementation checkpoint system for recreating the Rust server in Node / TypeScript. The goal is not a line-for-line port. The goal is a faithful interpretation of the same architectural idea in idiomatic TypeScript.

The Rust server demonstrates a clean/onion architecture with CQRS, ports/adapters, a transactional unit of work for writes, separate read repositories for queries, and very thin transports. The Node version must preserve those ideas while avoiding Rust-specific patterns such as newtype wrappers, trait object ergonomics, `Arc<Mutex<_>>`, `Result`-style plumbing, and crate-shaped ceremony that does not belong in TypeScript.

## Core Architectural Rules

Business logic must live only in the domain model and the application command/query handlers. Transports may parse, validate basic request shape, map DTOs, call application services, and map errors. Transports must not decide domain transitions, enforce loan rules, or perform multi-repository workflows.

Commands own mutating use cases. They use a write unit of work, transactional write repositories, and domain behavior to orchestrate changes. Queries own read use cases. They use read repositories directly and do not share command transaction machinery.

Import boundaries should mirror the Rust crate boundaries:

- `domain` imports no workspace package.
- `application` imports `domain` only.
- `auth-core` is independent of application/domain business logic.
- `persistence` imports `domain` and implements domain/application ports.
- `server-bootstrap` imports application, persistence, and auth-core to compose dependencies.
- `http-server` imports server-bootstrap, auth-core types as needed, and its own HTTP DTOs.
- `graphql-server` imports server-bootstrap, auth-core types as needed, and its own GraphQL DTOs.

HTTP and GraphQL types must be separate. A `BookResponseBody` for HTTP and a `CatalogTitle` GraphQL object may expose similar fields, but they are transport-specific contracts and should not be shared as application/domain types.

## Target Shape

Use a Node workspace under `servers/node` with packages:

- `@library/domain`
- `@library/application`
- `@library/auth-core`
- `@library/persistence`
- `@library/server-bootstrap`
- `@library/http-server`
- `@library/graphql-server`

Use TypeScript project references so each package has an explicit build boundary. Add a lightweight boundary-check script so the architecture is enforced without depending on an external linter plugin.

## Implementation Order

1. Establish the workspace, TypeScript configs, package manifests, and import-boundary checker.
2. Implement domain entities, status values, domain errors, repository ports, and unit-of-work ports.
3. Implement application command/query services that preserve the Rust use cases:
   - catalog commands: add book, add copy, lost/found transitions, maintenance transitions
   - lending commands: check out, return, report loaned copy lost
   - membership commands: register, suspend, reactivate
   - catalog queries: book catalog, book by ISBN, copy details
   - lending queries: member loans, overdue loans
   - membership queries: member details
4. Implement auth-core as an auth verification port plus JWT adapter.
5. Implement PostgreSQL persistence adapters:
   - read repos use read-only pool
   - write repos are transaction-scoped
   - unit of work starts a transaction and commits/rolls back once
   - SQL preserves the Rust read/write query split and lock-for-update behavior
6. Implement server-bootstrap as the composition root.
7. Implement HTTP transport with separate HTTP DTOs, thin handlers, auth middleware, and HTTP error mapping.
8. Implement GraphQL transport with separate GraphQL object/input types, thin resolvers, auth middleware, and GraphQL error mapping.
9. Add focused tests where they matter most: domain transition rules, command orchestration with fake repos/UoW, and boundary checker.
10. Run typecheck/tests and record final verification status.

## Checkpoint Protocol

The source of truth for implementation progress is `todo.json`. Each item has a `status` of `todo`, `in_progress`, or `done`. Before starting a task, mark it `in_progress`. When the task is complete, mark it `done`.

After each task is marked `done`, create a memory dump under `memory_dump/` using the task number in the filename. The memory dump should include:

- what changed
- why the change matters architecturally
- files created or modified
- key decisions and tradeoffs
- verification performed
- next task and context needed to continue

After writing each memory dump, compact working context mentally: keep only the Rust architecture facts, current Node structure, outstanding tasks, and decisions needed for the next task. Do not rely on conversation history as the only source of truth.
