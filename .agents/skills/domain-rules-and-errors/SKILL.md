# Domain Rules And Errors

Use this skill when adding or changing business decisions, state transitions, entity guards, default state preparation, or domain errors.

## Purpose

The domain layer owns business truth. It should be possible to read the domain module and understand what states exist, which transitions are legal, what defaults are applied, and why a business action can fail.

## What Belongs In Domain

- Entity and ID types.
- Business enums such as member or copy status.
- Creation payloads and prepared payloads.
- Default business values assigned during preparation.
- Guard methods such as `ensure_can_borrow`.
- Transition methods such as `suspend`, `reactivate`, `mark_lost`, `mark_found`.
- Typed business errors such as `LoanLimitReached` or `CannotBeBorrowed`.
- Repository port traits that describe what the core needs.

## Method Shape

Use domain methods to express decisions:

```rust
member.ensure_can_borrow()?;
member.ensure_within_loan_limit(active_loan_count)?;
let next_status = book_copy.mark_lost()?;
```

Prefer this over open-coded `if` statements in commands or transports.

## Error Rules

- Invalid business states return typed domain errors.
- Do not return strings for business failures.
- Do not collapse domain errors into transport status codes in the domain or application layer.
- Let `CommandError` preserve the typed domain error until the transport maps it.

## Preparation Rules

Creation payloads are raw application inputs after transport parsing. Prepared payloads are domain-approved values ready for persistence.

Examples from this project:

- `MemberCreationPayload::prepare` adds generated `MemberIdent` and sets member status to `Active`.
- `BookCopyCreationPayload::prepare` sets copy status to `Active`.
- `LoanCreationPayload::prepare` carries approved IDs into a `LoanPrepared`.

Put default domain state in `prepare`, not in SQL defaults or transport DTOs, when the default is business meaning.

## Timestamp Rule

The domain entity contains timestamps, but business decisions should not depend on database-generated timestamps. Created return entities use one code-side `Utc::now()` for both `dt_created` and `dt_modified`. Updated return entities use `Utc::now()` for `dt_modified` while preserving the previously loaded entity fields.

## Tests

Domain tests should cover:

- Allowed transitions.
- Rejected transitions.
- Guard success and failure.
- Preparation defaults.

These tests are cheap and should be the first line of coverage for business behavior.

## Smells

- A transport or repository checks whether a member can borrow.
- A command duplicates transition rules instead of calling a domain method.
- A status string is compared outside the domain enum conversion boundary.
- A business failure is represented as `anyhow!`, `String`, or a generic 500.
