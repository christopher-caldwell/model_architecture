# Skill: Go Domain Rules and Transitions

## Purpose

Implement business behavior in the Go domain layer. Preserve the Rust intent: rules are visible, pure, and tested.

## Rule style

Prefer methods or pure functions on domain types:

```go
func (m Member) EnsureCanBorrow(activeLoanCount int) error {
    if m.Status == MemberStatusSuspended {
        return ErrMemberSuspended
    }
    if activeLoanCount >= m.MaxActiveLoans {
        return ErrMemberLoanLimitReached
    }
    return nil
}
```

## Transition style

For state changes, prefer explicit transition functions:

```go
func (m Member) Suspend() (MemberStatus, error) {
    if m.Status == MemberStatusSuspended {
        return "", ErrMemberAlreadySuspended
    }
    return MemberStatusSuspended, nil
}
```

The application command calls the transition and persists the returned next state.

## Expected transition/rule areas

Represent member suspension/reactivation, member borrowing eligibility, book copy circulation/maintenance/lost behavior, checkout eligibility, active loan limits, and loan return behavior according to the language-neutral spec.

## Naming rules

Name methods according to what they actually check. If active copy status means physically circulatable, avoid `CanBeBorrowed` unless active-loan state is also checked. Prefer names like `EnsureCirculatable` or `CanEnterCirculation` when loan state is checked separately.

## Error rules

Domain rules return domain errors. They do not return HTTP status codes, GraphQL errors, SQL errors, or stringly typed codes as the primary contract.

## Test rules

Every non-trivial rule gets table-driven tests. Tests should assert with `errors.Is` or `errors.As`, not string matching.

## Anti-patterns

Handlers implementing business checks, repositories deciding rule validity, SQL hiding domain transitions, application setting arbitrary statuses without asking domain, or domain methods starting transactions.

## Non-negotiable guardrails

- Treat the Rust implementation as the behavioral reference, not syntax to copy.
- Preserve the same app, use cases, route contracts, GraphQL operations, persistence semantics, and architectural boundaries.
- Use idiomatic Go. Do not write Rust-in-Go clothing.
- Do not invent framework patterns, service locators, active-record models, generic `models` dumping grounds, or `index.go` files.
- Do not move business rules into HTTP handlers, GraphQL resolvers, SQL adapters, bootstrap wiring, generated code, or tests.
- Do not let transport DTOs, database rows, ORM models, SQL types, driver types, or generated GraphQL types cross inward into domain or application code.
- When a Rust construct does not translate directly, preserve the intent using normal Go mechanisms.
- Make uncertainty explicit. If the Rust reference does not contain a behavior, do not invent it.

## Required completion review

Before considering work from this skill complete, run two additional passes.

### Review pass 1: parity and boundary check

- Confirm the result still matches the Rust reference behavior.
- Confirm no dependency points inward incorrectly.
- Confirm no transport, persistence, framework, or generated-code concern leaked into domain/application.
- Confirm no generic Go tutorial pattern replaced the architecture of this project.
- Confirm every generated or changed artifact is justified by the Rust inventory or language-neutral spec.

### Review pass 2: Go idiom and drift check

- Confirm the Go code is idiomatic for package naming, errors, context, interfaces, transactions, and tests.
- Confirm the model did not create Rust-shaped Go abstractions where Go has a clearer native idiom.
- Confirm names are explicit and project-specific, not vague names like `models`, `common`, `manager`, `service2`, or `index`.
- Confirm all changed examples/docs/tests still agree with the generated code.
- Confirm unknowns were documented instead of guessed.
