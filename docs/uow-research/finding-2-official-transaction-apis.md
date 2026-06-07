# Finding 2: Official Rust Transaction APIs

## Isolation Protocol

This round ignored Finding 1 and did not reuse its evidence while deciding. It only searched official Rust database-library APIs and documentation for transaction ownership, commit semantics, rollback semantics, and call-site shape.

Search queries used:

- `docs.rs sqlx Transaction commit self rollback drop transaction Rust`
- `docs.rs sea-orm DatabaseTransaction commit rollback TransactionTrait Rust`
- `docs.rs diesel Connection transaction closure commit rollback Rust`
- `docs.rs tokio-postgres Transaction commit rollback Rust`

## Sources Selected

### SQLx `Transaction`

SQLx 0.8.6, the version family used by this repo, documents `Transaction` as an in-progress transaction/savepoint. It starts from `Pool::begin` or `Connection::begin`, should end with `commit` or `rollback`, and rolls back on drop if neither is called. Its example uses `let mut tx = conn.begin().await?`, executes through `&mut *tx`, then calls `tx.commit().await`. Source: <https://docs.rs/sqlx/0.8.6/sqlx/struct.Transaction.html>

Important facts from this source:

- Transaction is a concrete owned type.
- `commit(self)` consumes it.
- `rollback(self)` consumes it.
- Drop rolls back if still in progress.
- Query execution is mutable because the transaction is a single connection-backed scope.

### SeaORM `DatabaseTransaction`

SeaORM documents `DatabaseTransaction` as a concrete transaction type with `commit(self)` and `rollback(self)`. Its `TransactionTrait` can begin a transaction and return `DatabaseTransaction`, and its closure transaction API commits on success and rolls back on error. Source: <https://docs.rs/sea-orm/latest/sea_orm/struct.DatabaseTransaction.html>

Important facts from this source:

- Transaction is a concrete type.
- Commit and rollback consume the transaction.
- Query/execute methods are available through the transaction object.
- The library also supports closure-managed transaction scope.

### Diesel `Connection::transaction`

Diesel documents `Connection::transaction` as running a closure inside a transaction. It commits if the closure returns `Ok` and rolls back if it returns `Err`. Source: <https://docs.rs/diesel/latest/diesel/connection/trait.Connection.html>

Important facts from this source:

- The transaction boundary is explicit at the call site.
- The transaction scope is represented by a closure over a connection.
- Commit/rollback is centralized by the transaction API.

### `tokio-postgres::Transaction`

`tokio-postgres` documents `Transaction` as a concrete PostgreSQL transaction. It rolls back implicitly on drop, and `commit(self)` consumes the transaction to commit changes. Source: <https://docs.rs/tokio-postgres/latest/tokio_postgres/struct.Transaction.html>

Important facts from this source:

- Transaction is concrete and owned.
- Commit consumes it.
- Drop rolls back.
- Nested transactions are represented as savepoints.

## Sources Rejected

Search results from Reddit and general blog posts were not used for this round because the goal was official Rust database APIs only.

## Independent Decision

The official Rust database APIs favor one of two readable models:

1. a concrete transaction handle that is owned, used mutably, and consumed by commit, or
2. a closure transaction API that makes the transaction scope explicit.

They do not directly judge this app's raw boxed trait-object shape. They do show that the idiomatic user-facing transaction shape is an owned concrete handle or closure scope. The conclusion that boxed trait-object plumbing should be hidden is a local legibility and architecture inference from those APIs.

## Local Recommendation From This Round

The app should expose a concrete domain/application UoW handle to commands:

```rust
let mut uow = self.uow_factory.build().await?;
uow.member().get_by_ident_for_update(&ident).await?;
uow.commit().await?;
```

Internally, that handle can wrap the dynamic persistence implementation. This preserves the official transaction-handle readability without exposing SQLx.
