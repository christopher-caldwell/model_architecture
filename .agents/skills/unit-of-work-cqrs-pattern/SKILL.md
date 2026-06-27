---
name: unit-of-work-cqrs-pattern
description: Use when adding or changing write commands, write repository ports, unit-of-work traits, transaction boundaries, command-side reads, SQLx transaction adapters, dependency wiring for commands, or tests around commit behavior in this onion/CQRS Rust project.
---

# Unit Of Work CQRS Pattern

Use this skill when adding or changing write commands, write repository ports, unit-of-work traits, transaction boundaries, command-side reads, SQLx transaction adapters, dependency wiring for commands, or tests around commit behavior.

## Core Rule

Commands inject a factory, not a live transaction.

```text
command method starts
  build one WriteUnitOfWork from WriteUnitOfWorkFactory
  load write-decision state through the UoW
  ask domain objects for business decisions
  persist approved writes through the UoW
  commit once at the end
command method ends
```

Queries do not use the write UoW by default. They use read repository ports directly.

## Current Shape

The current project intentionally keeps repository ports and UoW traits in `domain`.

Important files:

- `server/crates/domain/src/uow.rs`
- `server/crates/domain/src/port_error.rs`
- `server/crates/domain/src/*/port.rs`
- `server/crates/application/src/commands/*.rs`
- `server/crates/persistence/src/uow.rs`
- `server/crates/http_server/src/deps.rs`
- `server/crates/graphql_server/src/deps.rs`

The abstractions are:

```rust
pub trait UnitOfWorkPort: Send {
    fn book(&mut self) -> &mut dyn BookWriteRepoPort;
    fn book_copy(&mut self) -> &mut dyn BookCopyWriteRepoPort;
    fn member(&mut self) -> &mut dyn MemberWriteRepoPort;
    fn loan(&mut self) -> &mut dyn LoanWriteRepoPort;
    async fn commit(self: Box<Self>) -> PortResult<()>;
}

pub struct WriteUnitOfWork {
    inner: Box<dyn UnitOfWorkPort>,
}

pub trait WriteUnitOfWorkFactory: Send + Sync {
    async fn build(&self) -> PortResult<WriteUnitOfWork>;
}
```

`commit(self)` consumes the UoW so application code cannot use the transaction after commit.

## Dependency Injection

Command structs store the factory behind `Arc<dyn WriteUnitOfWorkFactory>`:

```rust
#[derive(Clone)]
pub struct LendingCommands {
    uow_factory: Arc<dyn WriteUnitOfWorkFactory>,
}

impl LendingCommands {
    pub fn new(uow_factory: Arc<dyn WriteUnitOfWorkFactory>) -> Self {
        Self { uow_factory }
    }
}
```

Each command method builds a fresh UoW:

```rust
let mut uow = self.uow_factory.build().await?;
```

Do not inject `PgPool`, `Transaction`, `SqlUnitOfWork`, or persistence repositories into application commands.

## Command Pattern

Use `&mut WriteUnitOfWork` for helper methods that load state inside the transaction:

```rust
async fn get_member_by_ident(
    &self,
    uow: &mut WriteUnitOfWork,
    member_ident: &str,
) -> Result<Member, CommandError> {
    let ident = MemberIdent(member_ident.to_owned());
    uow.member()
        .get_by_ident_for_update(&ident)
        .await?
        .ok_or(MemberError::NotFound.into())
}
```

A write command should:

1. Build one UoW.
2. Load all state needed for write decisions through write repositories.
3. Use locking read methods such as `get_by_*_for_update` when the row participates in a write decision.
4. Ask domain objects to decide through guards and transitions.
5. Call write repository methods.
6. Commit once at the end.
7. Return a domain value without post-write hydration.

Do not call a query from a command to make write decisions. Use the UoW's write repositories.

## Persistence Adapter Pattern

Persistence owns SQLx.

`SqlUnitOfWork` wraps one SQLx transaction and implements all write repo ports:

```rust
pub struct SqlUnitOfWork {
    tx: Transaction<'static, Postgres>,
}

impl UnitOfWorkPort for SqlUnitOfWork {
    fn member(&mut self) -> &mut dyn MemberWriteRepoPort {
        self
    }

    async fn commit(self: Box<Self>) -> PortResult<()> {
        self.tx
            .commit()
            .await
            .context("Failed to commit transaction")
            .map_err(|error| PortError::unit_of_work(error.into_boxed_dyn_error()))
    }
}
```

The factory begins the transaction:

```rust
impl WriteUnitOfWorkFactory for SqlWriteUnitOfWorkFactory {
    async fn build(&self) -> PortResult<WriteUnitOfWork> {
        let tx = self.pool.begin().await?;
        Ok(WriteUnitOfWork::new(Box::new(SqlUnitOfWork { tx })))
    }
}
```

Repository method implementations call SQL helpers with `&mut *self.tx`.

## Reading With Pool Or Transaction

When the same read SQL can run against either a pool or a transaction, share a lower-level helper over a SQLx executor. Keep application-facing read and write ports separate.

Use this rule:

- Same SQL semantics, different handle: share a generic executor helper.
- Different SQL semantics, such as `FOR UPDATE`: use a separate helper and SQL file.
- Shared row-to-domain mapping is fine.
- Do not make application code care whether the handle is a pool or transaction.

Command-side reads that affect writes belong under `sql/<concept>/commands/` and should use write repository methods.

Examples:

- `get_by_ident_for_update`
- `get_by_barcode_for_update`
- `find_active_by_book_copy_id_for_update`

## Error Boundary

Domain ports return `PortResult<T>`.

Persistence may use `anyhow::Context` internally, but converts adapter failures before returning through a port:

```rust
.map_err(|error| PortError::repository(error.into_boxed_dyn_error()))
```

Application errors preserve infrastructure failures:

- commands map `PortError` to `CommandError::Infrastructure`
- queries map `PortError` to `QueryError::Infrastructure`
- HTTP and GraphQL map those at the transport edge

Do not return `anyhow::Result` from domain ports or application query/command boundaries.

## Write Result Rules

Within UoW-backed writes:

- Do not read after write only to hydrate a response.
- Insert SQL returns only generated values the application cannot know.
- Update SQL returns `()` unless a generated value is truly needed.
- Created domain values use one code-side `Utc::now()` for both `dt_created` and `dt_modified`.
- Updated domain values are shaped in the command from the locked pre-update entity, domain-approved value, and code-side `Utc::now()`.

## Composition Roots

HTTP and GraphQL each own their server-specific dependency wiring.

The composition root should:

1. Build read and write pools.
2. Build `SqlWriteUnitOfWorkFactory` from the write pool.
3. Build read repositories from the read pool.
4. Build command structs with the shared UoW factory.
5. Build query structs with read repositories.
6. Store commands and queries in transport state.

Do not construct UoWs, transactions, or persistence adapters in handlers or resolvers.

## Tests

For UoW-sensitive changes, prefer at least one focused test at the lowest useful layer:

- domain tests for guards, transitions, and preparation defaults
- application tests with a fake UoW when the command orchestration is non-trivial
- persistence/integration tests for SQL transaction behavior, locking, and rollback when available

The most useful application fake-UoW test asserts:

- the command builds one UoW
- command-side reads use write repos
- writes happen before commit
- commit happens once
- no post-write read is needed

## Smells

- A command stores a live UoW or transaction instead of a factory.
- Application code imports SQLx or persistence adapter types.
- A handler or resolver opens a transaction.
- A command uses a read repo for a write decision.
- A command commits before all writes for the use case are complete.
- A command re-queries after write only to hydrate the response.
- A write repository checks business rules instead of executing persistence.
- UoW methods expose broad database handles instead of narrow write repository ports.

## Related Skills

- `application-cqrs-use-cases`: use-case choreography and command/query placement.
- `persistence-write-results`: SQL return values, timestamp mapping, and no-read-after-write rules.
- `transport-adapters`: handler/resolver boundaries and transport error mapping.
- `domain-business-encapsulation`: guards, transitions, typed business errors, and domain defaults.
