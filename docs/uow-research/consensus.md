# UoW Research Consensus

## Files Compared

- [Finding 1: Concrete Core UoW Facade](finding-1-concrete-domain-uow-facade.md)
- [Finding 2: Official Rust Transaction APIs](finding-2-official-transaction-apis.md)
- [Finding 3: Clean Architecture And UoW Pattern Sources](finding-3-clean-architecture-and-uow-examples.md)
- [Finding 4: SQLx Wrapper Crates And Ergonomic Transaction Handles](finding-4-wrapper-crates-and-ergonomics.md)

Findings 2 through 4 were rerun as independent research rounds. Each used a separate search set and made its recommendation before comparison.

## Independent Conclusions

Finding 1 concluded that a core-owned concrete `WriteUnitOfWork` facade is the best readable shape for this app. In this repository that port boundary currently lives in the `domain` crate.

Finding 2, using only official Rust database APIs, concluded that established Rust transaction APIs favor a concrete owned transaction handle or closure transaction scope. Applied to this app, that points to a concrete domain/application UoW handle rather than raw boxed trait-object syntax.

Finding 3, using only UoW and clean architecture pattern sources, concluded that UoW should be a thin explicit application/service-layer object that exposes repository access, commits explicitly, and can be faked for tests.

Finding 4, using only SQLx wrapper-crate evidence, concluded that wrapper types are a known solution to SQLx transaction ergonomics, but this app should use its own domain/application-safe wrapper rather than exposing SQLx or moving transactions into transport.

## Agreement

All four findings agree on these points:

1. The application should hold a readable UoW/transaction-scope handle.
2. The UoW should be explicit in command code.
3. Commit should be explicit and should end/consume the UoW.
4. Uncommitted scope should roll back by default as an infrastructure safety behavior.
5. Repository grouping should remain visible through the UoW.
6. SQLx transaction/executor types must remain private to persistence.
7. Raw boxed trait-object plumbing should not leak into command helper signatures if a tiny facade removes it.

## Disagreement Or Constraints

The sources differ on implementation details:

- SQLx and tokio-postgres expose concrete transaction handles.
- Diesel prefers a closure transaction API.
- SeaORM supports both closure and explicit transaction handles.
- `axum-sqlx-tx` places transaction scope at request/transport level.
- UoW literature places transaction scope in the service/application workflow.

For this app, the decisive constraints are:

- commands own write workflows,
- HTTP and GraphQL must share command behavior,
- SQLx must stay in persistence,
- readability is a primary goal,
- repository grouping is still valuable.

Those constraints reject transport-level transaction middleware, application-visible SQLx transactions, and flattened UoW mega-interfaces.

## Consensus Recommendation

Implement a concrete `WriteUnitOfWork` facade at the core port boundary currently housed in `domain`.

Keep `UnitOfWorkPort` as the hidden dynamic implementation trait:

```rust
#[async_trait]
pub trait UnitOfWorkPort: Send {
    fn book(&mut self) -> &mut dyn BookWriteRepoPort;
    fn book_copy(&mut self) -> &mut dyn BookCopyWriteRepoPort;
    fn member(&mut self) -> &mut dyn MemberWriteRepoPort;
    fn loan(&mut self) -> &mut dyn LoanWriteRepoPort;
    async fn commit(self: Box<Self>) -> anyhow::Result<()>;
}
```

Add a concrete facade:

```rust
pub struct WriteUnitOfWork {
    inner: Box<dyn UnitOfWorkPort>,
}

impl WriteUnitOfWork {
    pub fn new(inner: Box<dyn UnitOfWorkPort>) -> Self {
        Self { inner }
    }

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

Change the factory:

```rust
#[async_trait]
pub trait WriteUnitOfWorkFactory: Send + Sync {
    async fn build(&self) -> anyhow::Result<WriteUnitOfWork>;
}
```

Command code becomes:

```rust
let mut uow = self.uow_factory.build().await?;

let member = self
    .get_member_by_ident(&mut uow, &input.member_ident)
    .await?;

let book_copy = self
    .get_book_copy_by_barcode(&mut uow, &input.book_copy_barcode)
    .await?;

uow.loan().create(&prepared).await?;
uow.commit().await?;
```

## Final Decision

The objectively best readable option for this app is:

```text
domain / core port boundary
  WriteUnitOfWork facade
  UnitOfWorkPort hidden implementation trait

application
  commands use &mut WriteUnitOfWork
  commands access grouped repo ports through uow.member(), uow.loan(), etc.

persistence
  SqlUnitOfWork implements UnitOfWorkPort and owns SQLx Transaction directly
```

This keeps the core idea and functionality, improves legibility, and preserves onion architecture boundaries.
