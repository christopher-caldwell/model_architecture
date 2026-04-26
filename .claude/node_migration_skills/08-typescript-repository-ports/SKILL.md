# TypeScript Repository Ports

## Purpose
Guide a model through the typescript repository ports part of migrating the confirmed Rust modular library demo to idiomatic Node / TypeScript without changing behavior or architectural boundaries.

## When to use this skill
Use this when a migration step needs to preserve the Rust app's confirmed shape while translating the architecture into TypeScript. This skill is model-facing: it tells the implementer what to inspect and what to produce; it is not an implementation itself.

## Inputs the model must inspect
- Canonical repository: https://github.com/christopher-caldwell/model_architecture
- Rust crate/module files relevant to this skill.
- Database DDL at `database/sql/library/library/ddl.sql` when persistence, migrations, contracts, or tests are affected.
- HTTP router/schema files when transport contracts are affected.
- GraphQL router/schema files when GraphQL contracts are affected.
- Existing tests and fixtures before creating any test guidance.

## Outputs the model should produce
- A focused TypeScript migration artifact for this skill's concern.
- A short source-parity note listing which Rust files were inspected.
- Boundary notes explaining where produced code belongs and what it must not import.
- Any examples clearly labeled as confirmed Rust app facts or structural placeholders.

Required rules
- Treat the Rust repository as the only behavioral source of truth.
- Inspect the exact Rust files before producing migration output. Do not use memory or previous research.
- Preserve dependency direction: domain inward, application orchestration next, adapters outside, bootstrap wiring outermost.
- Do not let HTTP, GraphQL, auth, database rows, migration tools, OpenAPI, validation libraries, or generated types define domain behavior.
- Keep transport DTOs, persistence rows, and domain models separate.
- Do not invent entities, fields, routes, commands, queries, fixtures, seeds, tables, or framework choices.

Node / TypeScript-specific guidance
- Use TypeScript with strict compiler settings, package.json, tsconfig.json, and Vitest or Jest guidance.
- Prefer plain domain classes/types, branded IDs where helpful, and explicit domain methods/functions.
- Use async/await for all persistence and transport adapter calls.
- Use interfaces for structural ports unless an abstract class is needed for shared protected behavior or runtime instanceof checks.
- Validation libraries may parse edge DTOs only; they are not the domain source of truth.
- Avoid giant models.ts/services.ts files and avoid barrels that hide boundaries.

Rust-to-TypeScript translation cautions
- Translate traits to TypeScript ports, not to framework services.
- Translate Result/thiserror anyhow layering to typed domain/application errors and edge mappers.
- Translate Arc/dyn wiring to explicit constructor injection from the composition root.
- Preserve transaction semantics rather than Rust ownership mechanics.

Anti-hallucination checks
- Every app-specific name must be found in the Rust repo or marked as a structural placeholder.
- Examples must use confirmed names only or say STRUCTURAL PLACEHOLDER.
- If key files cannot be inspected, stop and report that the output cannot be reliable.

Boundary checks
- Domain imports no framework/database/HTTP/GraphQL/auth adapter code.
- Application imports domain and ports only, not concrete adapters.
- Persistence implements ports and owns row mapping.
- Transports map DTOs to application inputs and domain/application errors to protocol responses.

Completion checklist
- Source files inspected and cited in notes.
- Outputs preserve command/query separation, read/write repo separation, UoW transaction boundaries, and DTO/domain/row separation.
- No framework or tool becomes the architecture.
- No unverified app behavior appears in generated code or docs.

## Confirmed Rust app anchors
Confirmed source facts from the Rust repo:
- Rust workspace crates include application, auth_core, domain, graphql_server, http_server, persistence, and server_bootstrap.
- Domain entities confirmed: Book, BookCopy, Loan, Member.
- Domain fields confirmed include Book(id,isbn,dt_created,dt_modified,title,author_name); BookCreationPayload(isbn,title,author_name); Member response fields ident, dt_created, dt_modified, status, full_name, max_active_loans; BookCopy response fields barcode, dt_created, dt_modified,status; Loan response fields ident, dt_created, dt_modified, dt_due, dt_returned.
- Commands confirmed: add_book, add_book_copy, mark_book_copy_lost, mark_book_copy_found, send_book_copy_to_maintenance, complete_book_copy_maintenance, register_member, suspend_member, reactivate_member, check_out_book_copy, return_book_copy, report_lost_loaned_book_copy.
- Queries confirmed: get_book_catalog, get_book_by_isbn, get_book_copy_details, get_member_details, get_member_loans, get_overdue_loans.
- UoW confirmed: write UoW exposes book, book_copy, membership, and loan write repositories and commit; factory builds a UoW.
- HTTP paths confirmed: /books, /books/{isbn}/copies, /book-copies/{barcode}, /book-copies/{barcode}/lost, /book-copies/{barcode}/maintenance, /book-copies/{barcode}/return, /book-copies/{barcode}/report-loss, /members, /members/{ident}, /members/{ident}/suspension, /members/{ident}/loans, /loans, /loans/overdue, plus public health and Swagger/OpenAPI endpoints.
- Database schema confirmed in database/sql/library/library/ddl.sql: schema library; tables struct_type, book, book_copy, member, loan; identity integer/smallint PKs; dt_created/dt_modified defaults; update trigger; unique indexes on struct_type(group_name,att_pub_ident), book(isbn), book_copy(barcode), member(member_ident), loan(loan_ident), and loan(book_copy_id,dt_returned); loan dt_due and dt_returned use sentinel default '9999-01-01'.
- Auth boundary confirmed in auth_core with Claims{sub, exp}; comments say additional JWT fields may exist but are not facts.

## Skill-specific emphasis
For **TypeScript Repository Ports**, preserve the exact confirmed Rust names and behavior where this skill touches them. If this skill needs details not listed above, inspect the repository again and add only confirmed facts.

## Per-skill review record
Before finalizing output made with this skill, perform two review passes:
1. Parity pass: compare every app-specific name, route, field, status, table, and behavior against the Rust repo.
2. Architecture pass: remove any TypeScript convenience that makes an edge framework, ORM, validation library, or database row the owner of core behavior.
