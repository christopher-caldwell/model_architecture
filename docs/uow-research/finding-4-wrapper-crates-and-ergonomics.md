# Finding 4: SQLx Wrapper Crates And Ergonomic Transaction Handles

## Isolation Protocol

This round ignored Findings 1, 2, and 3 and searched only for crates that wrap SQLx transactions for ergonomic or safety reasons. It did not use general UoW literature or official ORM transaction APIs to decide.

Search queries used:

- `site:docs.rs sqlx transaction wrapper automatic rollback commit Rust crate`
- `site:docs.rs axum sqlx transaction extractor commit rollback wrapper Tx Rust`
- `site:docs.rs "TransactionContext" "SQLx" "commit" "rollback" Rust`
- `site:docs.rs "&mut Tx" implements "sqlx::Executor"`

## Sources Selected

### `axum-sqlx-tx`

`axum-sqlx-tx` provides a concrete `Tx` extractor for request-bound SQLx transactions. The crate docs say the transaction begins when first used, is stored in request extensions, and is committed or rolled back based on response status. It explicitly says using the extractor instead of direct SQLx transactions means users cannot forget to commit under that request model. Source: <https://docs.rs/axum-sqlx-tx>

Important facts from this source:

- The library introduces a named concrete transaction wrapper.
- The wrapper improves call-site ergonomics.
- The wrapper centralizes transaction resolution.
- It exposes SQLx because it is designed for handler-level SQLx use.

### `axum_sqlx_tx::Tx`

The `Tx` type implements `&mut Tx` as `sqlx::Executor`, implements `Deref`/`DerefMut` to SQLx `Transaction`, and exposes `commit(self)`. Source: <https://docs.rs/axum-sqlx-tx/latest/axum_sqlx_tx/struct.Tx.html>

Important facts from this source:

- A concrete wrapper can hide awkward executor/reborrow mechanics.
- Commit consumes the wrapper.
- Ergonomics are improved by passing a named handle, not a raw inner transaction everywhere.

### `sqlx-transaction-manager`

`sqlx-transaction-manager` describes itself as a type-safe SQLx transaction wrapper with automatic rollback, compile-time transaction boundaries, an ergonomic `with_transaction` function, nested transaction support, and zero-runtime-overhead wrapping over SQLx. It is design evidence only for this project because it currently targets MySQL, while this app uses PostgreSQL. Source: <https://docs.rs/sqlx-transaction-manager/latest/sqlx_transaction_manager/>

Important facts from this source:

- Wrapper types are a known solution to SQLx transaction ergonomics.
- The wrapper has explicit transaction state and `commit`.
- The goal is safety and readability around SQLx's lower-level transaction mechanics.
- It is not a practical dependency recommendation for this app's Postgres backend.

### `TransactionContext`

`TransactionContext` wraps SQLx's `Transaction`, rolls back on drop if `commit()` is not called, consumes itself on commit, and exposes `as_executor()` for query execution. In the manual context API, commit is explicit; automatic resolution only applies to rollback-on-drop or the crate's closure helper. Source: <https://docs.rs/sqlx-transaction-manager/latest/sqlx_transaction_manager/context/struct.TransactionContext.html>

Important facts from this source:

- Wrapper around SQLx transaction is explicit.
- Consuming commit prevents reuse.
- `as_executor()` hides the lower-level executor detail from most call sites.

## Sources Rejected

`aide-axum-sqlx-tx` was not used as a primary source because it is a compatibility/re-export wrapper over `axum-sqlx-tx`, not an independent transaction design.

## Independent Decision

The ergonomic-wrapper ecosystem independently supports creating a named wrapper around SQLx transaction machinery when raw SQLx transaction handling creates call-site noise.

However, those crates expose SQLx because they are infrastructure or transport integration crates. This app has stricter onion boundaries, so the same wrapper idea should be adapted at the domain/application boundary rather than copied directly.

## Local Recommendation From This Round

Use a project-owned `WriteUnitOfWork` facade rather than adding one of these crates. The facade should do for application command code what these wrappers do for SQLx users: hide plumbing and present a clear transaction/UoW handle.

The command shape should be:

```rust
let mut uow = self.uow_factory.build().await?;
uow.member().get_by_ident_for_update(&ident).await?;
uow.commit().await?;
```

Not:

```rust
self.get_member_by_ident(uow.as_mut(), &ident).await?;
```
