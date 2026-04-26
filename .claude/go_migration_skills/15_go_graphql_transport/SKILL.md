# Skill: Go GraphQL Transport

## Purpose

Implement GraphQL as a second outer adapter. It calls the same application commands and queries as HTTP.

## Library choice

Use a project-approved GraphQL library, likely `gqlgen`. Lock the choice before implementation.

## Responsibilities

GraphQL owns schema definitions, generated resolver interfaces, GraphQL DTOs, resolver methods, auth context extraction, mapping inputs to application inputs, mapping domain models to GraphQL outputs, and GraphQL error mapping.

GraphQL does not own domain rules, repos, SQL, transactions, or application orchestration beyond calling commands/queries.

## Schema parity

For every Rust GraphQL query/mutation, define the Go equivalent: operation name, input shape, output shape, command/query called, auth behavior, and error behavior. Do not invent operations.

## Resolver flow

Resolvers receive GraphQL input, map to application input, call command/query, map domain result to GraphQL output, and return mapped GraphQL errors.

## Generated code boundary

Generated GraphQL types are transport concerns. They must not become domain models, application inputs, or persistence rows.

## Error mapping

Map to edge codes like `BAD_USER_INPUT`, `NOT_FOUND`, `CONFLICT`, `UNAUTHENTICATED`, `FORBIDDEN`, and `INTERNAL`, or to the Rust/spec-defined names.

## Anti-patterns

Resolvers importing persistence, starting UOW directly, leaking generated types inward, letting GraphQL schema drive the domain, duplicating business rules, or changing operation names silently.

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
