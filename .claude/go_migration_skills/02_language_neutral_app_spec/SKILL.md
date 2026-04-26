# Skill: Language-Neutral App Spec

## Purpose

Convert the Rust inventory into a language-neutral application specification. This spec defines the app behavior for Go, Node, Python, and future languages without binding it to Rust crates or Go packages.

## When to use

Use after the canonical Rust inventory and before writing Go code. Use again whenever Rust changes.

## Required output sections

1. Application purpose.
2. Core concepts.
3. Business rules.
4. Commands.
5. Queries.
6. Persistence model.
7. HTTP contract.
8. GraphQL contract.
9. Auth contract.
10. Error contract.
11. Conformance expectations.
12. Explicit unknowns.

## Core concept requirements

Define Book, Book Copy, Member, Loan, statuses, public identifiers, internal IDs, and externally visible identity fields. Describe what each concept means, not how Rust stores it.

## Business rule format

Write every rule in Given/When/Then form.

Example:

```text
Rule: A suspended member cannot borrow a book copy.
Given a member with suspended status,
when checkout is attempted,
then the command fails with a domain/application error that maps to a client-facing error at the transport edge.
```

## Command format

For each command, describe name, purpose, input, state loaded, rules applied, persistence changes, output, and failure cases.

## Query format

For each query, describe name, purpose, input, read source, output shape, and not-found behavior.

## Contract parity

HTTP routes and GraphQL operations must be recorded as contracts. The Go version may use idiomatic implementation details, but it must not silently change contract shape.

## Error contract

Separate domain failures, application failures, infrastructure failures, authentication failures, authorization failures, and transport mapping.

## Translation rules

The spec is not a code-generation prompt. It must preserve behavior, not Rust syntax. It must explicitly say which details require exact parity and which may be language-idiomatic. It must flag missing information instead of guessing.

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
