# Finding 3: Clean Architecture And UoW Pattern Sources

## Isolation Protocol

This round ignored Findings 1 and 2 and searched only for Unit of Work pattern literature and clean architecture/DDD examples. It did not use Rust database API docs to decide.

Search queries used:

- `Martin Fowler Unit of Work maintains list commit rollback pattern`
- `Cosmic Python unit of work repository pattern commit rollback`
- `Enterprise Application Architecture Unit of Work pattern transaction commit rollback repositories`
- `DDD unit of work repository pattern transaction clean architecture`

## Sources Selected

### Martin Fowler: Unit of Work

Fowler defines Unit of Work as maintaining objects affected by a business transaction and coordinating writes and concurrency resolution. Source: <https://martinfowler.com/eaaCatalog/unitOfWork.html>

Important facts from this source:

- UoW is tied to a business transaction, not a table-specific repository.
- Its job is to coordinate a set of changes.
- It exists to make one coherent write boundary.

### Cosmic Python: Unit of Work and Repository

Cosmic Python explicitly ties Repository, Service Layer, and Unit of Work. Its abstract UoW exposes repositories and has `commit()` and `rollback()`. It argues for explicit commits because the default behavior should be safe: no changes unless the use case explicitly commits. Source: <https://www.cosmicpython.com/book/chapter_06_uow.html>

Important facts from this source:

- UoW provides repository access.
- Application/service code uses UoW explicitly.
- Explicit commit is preferred for safety and reasoning.
- Fake UoWs make application logic testable without concrete persistence.

### Cosmic Python production-pattern article

The Cosmic Python article describes UoW as a thin layer over the persistence session with explicit commit and rollback points. It also shows repositories exposed through the UoW for convenient access. Source: <https://www.cosmicpython.com/blog/2017-09-08-repository-and-unit-of-work-pattern-in-python.html>

Important facts from this source:

- UoW is thin.
- UoW owns the persistence session.
- Repositories can be accessed from the UoW.
- Command handlers start and commit the UoW.
- The abstraction allows tests without concrete implementation details.

### Microsoft: Unit of Work and Persistence Ignorance

Microsoft's archived UoW article references Fowler's definition and shows a processor obtaining a UoW from a factory, executing commands against it, committing on success, and rolling back on exception. Source: <https://learn.microsoft.com/en-us/archive/msdn-magazine/2009/may/the-unit-of-work-pattern-and-persistence-ignorance>

Important facts from this source:

- A UoW factory starts a new UoW.
- UoW is passed to operations that need coordinated persistence.
- Commit and rollback belong to the UoW boundary.

## Sources Rejected

Reddit discussion and generic PyPI packages were not used to decide this round. They were too informal or too implementation-specific compared with Fowler/Cosmic Python/Microsoft.

## Independent Decision

The pattern literature independently supports a UoW object that:

- is thin,
- is explicit in application/service code,
- exposes repository access,
- owns or scopes the persistence session,
- commits explicitly,
- rolls back by default or on error,
- can be faked in tests.

It does not support flattening all repository methods into a large service-style object if repository grouping is meaningful. Cosmic Python explicitly shows repositories on the UoW as a convenience.

## Local Recommendation From This Round

The best fit for this app is a concrete `WriteUnitOfWork` object at the domain/application boundary. It should expose grouped repository accessors and explicit `commit`.

That directly matches the pattern sources:

```rust
let mut uow = self.uow_factory.build().await?;
let member = uow.member().get_by_ident_for_update(&ident).await?;
uow.loan().create(&prepared).await?;
uow.commit().await?;
```

The implementation can still be a boxed trait object hidden behind that concrete facade.

