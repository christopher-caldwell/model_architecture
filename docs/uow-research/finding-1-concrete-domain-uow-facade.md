# Finding 1: Concrete Core UoW Facade

## Question

What is the most readable UoW shape for this app after removing shared transaction state?

## Recommendation

Add a small core-owned concrete `WriteUnitOfWork` facade over the existing `Box<dyn UnitOfWorkPort>`.

In this codebase, "core-owned" means the facade belongs at the same port boundary that currently lives in the `domain` crate. This is a local onion-architecture placement decision, not a general DDD rule that UoW must be a domain entity.

The application should work with:

```rust
let mut uow = self.uow_factory.build().await?;

let member = self
    .get_member_by_ident(&mut uow, &input.member_ident)
    .await?;

uow.loan().create(&prepared).await?;
uow.commit().await?;
```

Instead of:

```rust
let mut uow = self.uow_factory.build().await?;

let member = self
    .get_member_by_ident(uow.as_mut(), &input.member_ident)
    .await?;

uow.loan().create(&prepared).await?;
uow.commit().await?;
```

The `WriteUnitOfWork` facade should live in `domain`, because it is only a port-level ownership handle. It must not expose SQLx.

Conceptual shape:

```rust
pub struct WriteUnitOfWork {
    inner: Box<dyn UnitOfWorkPort>,
}

impl WriteUnitOfWork {
    pub fn book(&mut self) -> &mut dyn BookWriteRepoPort {
        self.inner.book()
    }

    pub fn book_copy(&mut self) -> &mut dyn BookCopyWriteRepoPort {
        self.inner.book_copy()
    }

    pub fn member(&mut self) -> &mut dyn MemberWriteRepoPort {
        self.inner.member()
    }

    pub fn loan(&mut self) -> &mut dyn LoanWriteRepoPort {
        self.inner.loan()
    }

    pub async fn commit(self) -> anyhow::Result<()> {
        self.inner.commit().await
    }
}
```

The factory should return this concrete handle:

```rust
async fn build(&self) -> anyhow::Result<WriteUnitOfWork>;
```

## Evidence

SQLx models a transaction as a concrete owned value. Its docs say a transaction starts via `Pool::begin` or `Connection::begin`, should end with `commit` or `rollback`, and rolls back on drop if neither is called. `commit(self)` consumes the transaction. Source: <https://docs.rs/sqlx/0.8.6/sqlx/struct.Transaction.html>

SeaORM exposes the same conceptual model: use a transaction closure, or explicitly `begin`, perform work through the transaction, then `commit`; if the transaction goes out of scope it is rolled back. Source: <https://www.sea-ql.org/SeaORM/docs/advanced-query/transaction/>

`axum-sqlx-tx` is an established Rust transaction wrapper. It provides a concrete `Tx` type around SQLx transaction machinery, exposes a consuming `commit(self)`, and gives ergonomic transaction access. Source: <https://docs.rs/axum-sqlx-tx/latest/axum_sqlx_tx/struct.Tx.html>

## Local Fit

This app needs to keep SQLx private to `persistence`. Therefore the application cannot use SQLx's concrete transaction handle directly.

A core-owned `WriteUnitOfWork` facade gives the application the same readability benefits as an owned transaction handle without leaking the persistence implementation.

It also preserves the current repository grouping:

```rust
uow.member().get_by_ident_for_update(&ident).await?;
uow.book_copy().update_status(id, status).await?;
uow.loan().create(&prepared).await?;
```

That is more legible than flattening every repository method onto the UoW itself.

## Conclusion

The best readable option is not a new transaction abstraction, not a flattened service object, and not raw boxed trait-object call sites. It is a tiny concrete port-boundary facade that wraps the existing dynamic port and gives commands a named UoW value.
