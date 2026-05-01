---
name: domain-business-encapsulation
description: Preserve this project's domain encapsulation style. Use when Codex adds or changes domain entities, domain logic, business rules, guards, state transitions, creation/preparation payloads, typed domain errors, or domain enums such as MemberStatus and BookCopyStatus. Especially use when moving logic between domain/application/transport layers or when deciding how enum values should cross persistence and transport boundaries.
---

# Domain Business Encapsulation

This project uses a narrow domain style. The domain does not just hold structs and enums; it owns business decisions. Application commands orchestrate use cases, but they must ask domain objects to decide.

## Core Shape

For each domain concept, keep responsibilities split like this:

- `entity.rs`: IDs, entity structs, raw creation payloads, prepared payloads.
- `enums.rs`: domain state enums and their string conversion boundary.
- `errors.rs`: typed business failures.
- `logic.rs`: business predicates, public guards, public transitions, and `prepare` methods.
- `port.rs`: repository traits needed by the core.

Do not put business rules on DTOs, repositories, SQL files, handlers, resolvers, or bootstrap wiring.

## Encapsulation Pattern

Use private boolean predicates for reusable internal facts:

```rust
#[must_use]
fn can_borrow(&self) -> bool {
    self.status == MemberStatus::Active
}
```

Expose public guards for business validation:

```rust
pub fn ensure_can_borrow(&self) -> Result<(), MemberError> {
    if !self.can_borrow() {
        return Err(MemberError::CannotBorrowWhileSuspended);
    }
    Ok(())
}
```

Expose public transitions that return the next domain value, not a mutated entity:

```rust
pub fn suspend(&self) -> Result<MemberStatus, MemberError> {
    if self.status == MemberStatus::Suspended {
        return Err(MemberError::CannotBeSuspended);
    }
    Ok(MemberStatus::Suspended)
}
```

This shape is intentional:

- The private predicate captures the rule in one place.
- The public guard gives commands a readable way to enforce the rule.
- The public transition returns the approved next value.
- The command/repository performs persistence after the domain approves the decision.

## Do Not Expose Raw Decisions

Do not add public `is_*`, `can_*`, or status-check helpers for callers to branch on unless there is a strong read-only presentation need. Prefer `ensure_*` and transition methods that return typed errors.

Avoid this:

```rust
if member.status == MemberStatus::Suspended {
    return Err(...);
}
```

Avoid this:

```rust
if !member.can_borrow() {
    return Err(...);
}
```

Prefer this:

```rust
member.ensure_can_borrow()?;
```

The caller should not reconstruct the domain decision. The caller should ask the domain object.

## Guard Vs Transition

Use a guard when the operation only validates whether a later use-case step may happen:

```rust
member.ensure_can_borrow()?;
book_copy.ensure_can_be_borrowed()?;
loan.ensure_can_be_returned()?;
```

Use a transition when the operation changes a domain state:

```rust
let suspended_status = member.suspend()?;
let active_status = member.reactivate()?;
let lost_status = book_copy.mark_lost()?;
```

Transitions return the new enum/status value. They do not write to the database and do not construct the updated entity. Application commands use the returned value to call the write repository and shape the response.

## Creation Payloads And Prepared Payloads

Use creation payloads for raw application input after transport parsing. Use prepared payloads for domain-approved values ready for persistence.

`prepare` belongs in `logic.rs` and applies domain defaults:

```rust
impl MemberCreationPayload {
    #[must_use]
    pub fn prepare(self, ident: MemberIdent) -> MemberPrepared {
        MemberPrepared {
            ident,
            full_name: self.full_name,
            max_active_loans: self.max_active_loans,
            status: MemberStatus::Active,
        }
    }
}
```

Do not assign domain defaults in HTTP/GraphQL DTOs, SQL defaults, or repositories when the default has business meaning. Examples:

- New members start `Active`.
- New book copies start `Active`.
- Loan creation carries approved member and copy IDs.

## Domain Errors

Every rejected business decision should have a typed error variant in the concept's `errors.rs`.

Examples:

- `MemberError::CannotBorrowWhileSuspended`
- `MemberError::LoanLimitReached`
- `BookCopyError::CannotBeBorrowed`
- `BookCopyError::CannotBeReturnedFromLost`
- `LoanError::CannotBeReturned`

Do not use `anyhow`, strings, HTTP status codes, GraphQL errors, or SQL errors to represent expected business rejections. Those mappings happen outside the domain.

## Domain Enum Handling

Domain enums are the source of truth for business states. Keep them as Rust enums in the domain, not strings.

Use this exact boundary pattern:

```rust
#[derive(Debug, Clone, PartialEq, Eq)]
pub enum MemberStatus {
    Active,
    Suspended,
}

#[derive(thiserror::Error, Debug, Clone, PartialEq, Eq)]
#[error("Unknown member status '{input}'")]
pub struct ParseMemberStatusError {
    input: String,
}

impl std::fmt::Display for MemberStatus {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        let s = match self {
            Self::Active => "active",
            Self::Suspended => "suspended",
        };
        f.write_str(s)
    }
}

impl std::str::FromStr for MemberStatus {
    type Err = ParseMemberStatusError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            "active" => Ok(Self::Active),
            "suspended" => Ok(Self::Suspended),
            _ => Err(ParseMemberStatusError {
                input: s.to_owned(),
            }),
        }
    }
}
```

Rules:

- Use `Display` when an adapter needs to send the enum to persistence.
- Use `FromStr` when an adapter reads a string from persistence.
- Keep parse errors typed and local to the enum.
- Do not compare raw status strings in domain, application, transport, or persistence logic.
- Do not duplicate status string literals outside enum conversion and SQL lookup values.
- Do not use `serde` renaming as the domain's canonical enum boundary unless a transport DTO specifically needs it.

Persistence adapters may convert database rows into domain enums, then fail with context if the database contains an unknown value. Once inside domain/application code, work with enum variants only.

## Application Interaction

Commands should compose domain methods like this:

```rust
let member = self.get_member_by_ident(&*uow, &input.member_ident).await?;
let book_copy = self.get_book_copy_by_barcode(&*uow, &input.book_copy_barcode).await?;

member.ensure_can_borrow()?;
book_copy.ensure_can_be_borrowed()?;
member.ensure_within_loan_limit(active_loan_count as i16)?;

let prepared = LoanCreationPayload {
    member_id: member.id,
    book_copy_id: book_copy.id,
}.prepare();
```

The command loads state, calls domain decisions, persists approved changes, commits, and shapes the result. It should not reimplement the domain condition.

## Tests

Domain tests should live near the logic and cover:

- Every allowed transition.
- Every rejected transition.
- Guard success and failure.
- Private predicate behavior through public guards.
- `prepare` defaults.
- Enum `Display` and `FromStr` conversion for every variant and unknown input when useful.

Use small entity factories in tests, as the existing `member::logic` and `book_copy::logic` modules do.

## Review Checklist

- Is every business decision expressed as a domain method?
- Are private predicates private?
- Do callers use `ensure_*` or transition methods instead of checking status directly?
- Does each expected rejection have a typed domain error?
- Does every state transition return the next enum value?
- Are domain defaults applied in `prepare`?
- Are enum strings confined to `Display`, `FromStr`, and adapter/SQL lookup boundaries?
- Are tests covering both success and rejection paths?
