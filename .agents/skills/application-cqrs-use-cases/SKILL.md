# Application CQRS Use Cases

Use this skill when adding or changing commands, queries, use-case inputs, unit-of-work orchestration, transaction boundaries, or command/query composition.

## Purpose

The application layer owns use-case choreography. It decides which domain objects must be loaded, which domain decisions to ask for, which repositories to call, and when to commit.

Business workflow logic belongs here instead of in HTTP or GraphQL. Reusable business decisions still belong in domain methods and errors.

## Command Rules

Commands live in `server/crates/application/src/commands`.

A command should:

1. Build a write unit of work.
2. Load all records needed for write decisions through write repositories.
3. Use command-side locking reads such as `get_by_*_for_update` when the row participates in a write decision.
4. Ask domain objects to decide through guards and transitions.
5. Call write repository methods.
6. Commit once at the end.
7. Return a domain value without post-write hydration.

Commands may read during a write workflow, but those reads belong inside the same unit of work when they affect the write.

## Query Rules

Queries live in `server/crates/application/src/queries`.

A query should:

- Use read repository ports.
- Avoid transactions unless a read use case explicitly requires one.
- Return read results.
- Never mutate state.
- Never make transport decisions.

## No Transport Composition

Do not implement a use case by composing a query and command in HTTP or GraphQL.

Wrong shape:

```text
handler -> query
handler -> if not found return
handler -> command
```

Correct shape:

```text
handler -> one command
command -> load needed state
command -> domain decisions
command -> write
```

If a workflow requires branching, locking, or business interpretation, it belongs in an application command or query.

## No Read After Write

Do not run a select after an insert/update merely to hydrate the response. Use generated values returned by the write SQL plus the input/preloaded domain state to build the returned domain object.

Examples:

- Create: repository gets generated ID, then constructs the entity from prepared input plus `Utc::now()`.
- Update: command loads entity for update, asks domain for the new value, writes it, commits, then returns `{ new_value, dt_modified: Utc::now(), ..loaded_entity }`.

## Error Handling

- Domain errors flow into `CommandError`.
- Infrastructure failures get context and become unexpected command errors.
- Do not map command errors to HTTP status codes or GraphQL extensions in application code.

## Smells

- Command commits before all writes for a use case are complete.
- Command re-queries the updated row after commit.
- Query performs a write or calls a command.
- Handler or resolver calls more than one application use case to perform one business action.
- Application code contains SQLx, Axum, async-graphql, or DTO types.
