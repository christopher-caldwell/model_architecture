# Sub-Agent Validation

## Purpose

Five sub-agents independently rechecked the UoW research conclusion:

> Use a small core-owned `WriteUnitOfWork` facade over the hidden `UnitOfWorkPort` implementation trait, keeping SQLx private to persistence while making command code read as `&mut WriteUnitOfWork`.

They were asked to validate the conclusion, not invent alternatives unless the evidence contradicted it.

## Results

### Agent 1: Official Rust Database APIs

Verdict: supports.

Checked SQLx, SeaORM, Diesel, and tokio-postgres official transaction APIs. Confirmed that the established Rust database APIs favor either an owned concrete transaction handle consumed by commit or an explicit closure-scoped transaction. Noted that `domain-owned` is a local architecture placement, not something the DB APIs prescribe.

Corrections applied:

- Use versioned SQLx 0.8.6 docs where SQLx behavior is cited.
- Use stable SeaORM transaction docs URL.
- Avoid saying official APIs prove boxed trait objects are bad; they support concrete handles/closures, and the boxed-trait hiding is a local legibility inference.

### Agent 2: Clean Architecture / DDD / UoW Literature

Verdict: supported with wording correction.

Checked Fowler, Cosmic Python, Microsoft UoW/persistence ignorance, and Clean Architecture dependency-rule material. Confirmed support for a thin UoW used by application/service code, grouped repository access, explicit commit, rollback safety, and persistence details hidden behind inward-facing abstractions.

Corrections applied:

- Reworded `domain-owned` to `core-owned` / `port-boundary facade currently located in domain`.
- Clarified that the hidden implementation trait is a Rust/Clean Architecture implementation choice, not a pattern-literature requirement.

### Agent 3: SQLx Wrapper Ergonomics

Verdict: validated.

Checked `axum-sqlx-tx`, `axum_sqlx_tx::Tx`, `sqlx-transaction-manager`, and `TransactionContext`. Confirmed that named wrapper/facade types are established responses to SQLx transaction ergonomics, but these crates should not be adopted directly because they either expose SQLx at handler level or are not a Postgres fit.

Corrections applied:

- Qualified `sqlx-transaction-manager` as design evidence only because it currently targets MySQL.
- Clarified automatic commit/rollback wording.

### Agent 4: Rust Clean Architecture / SQLx Community Examples

Verdict: supports.

Checked SQLx transaction/executor docs and community discussions around executor generics, async traits, lifetime friction, and clean architecture with SQLx. Confirmed that exposing SQLx executor generics or raw transaction plumbing to application code is a poor fit for this architecture, and that a wrapper/facade is a reasonable Rust solution.

Corrections applied:

- Use SQLx 0.8.6 docs.
- Clarify that rollback-on-drop is a safety default, while explicit commit remains the command boundary.

### Agent 5: Local Codebase Fit

Verdict: yes, low risk.

Checked local `domain/src/uow.rs`, `persistence/src/uow.rs`, and application command files. Confirmed the facade would remove `uow.as_mut()` helper-call plumbing, keep SQLx private, preserve repository grouping, and require limited churn:

- add `WriteUnitOfWork` in `domain/src/uow.rs`,
- make `WriteUnitOfWorkFactory::build` return it,
- wrap `SqlUnitOfWork` in persistence factory,
- change command helper signatures from `&mut dyn UnitOfWorkPort` to `&mut WriteUnitOfWork`,
- replace `uow.as_mut()` call sites with `&mut uow`.

Correction noted:

- `UnitOfWorkPort` is not literally hidden from the whole codebase because persistence must implement it; it is hidden from application command ergonomics.

## Overall Validation

All five agents supported the conclusion. None recommended a competing design.

The validated recommendation is:

```text
application command
  -> receives a concrete WriteUnitOfWork facade
  -> uses grouped repo ports through uow.member(), uow.loan(), etc.
  -> explicitly commits once

persistence
  -> SqlUnitOfWork implements UnitOfWorkPort
  -> owns SQLx Transaction directly
  -> SQLx executor details remain private
```

