# Node Migration Skills Manifest

Repository inspected at canonical URL; no redirect was observed in the GitHub pages/API responses used.

## Skills
- `01-canonical-rust-app-inventory/SKILL.md` — Canonical Rust App Inventory
- `02-language-neutral-app-specification/SKILL.md` — Language-Neutral App Specification
- `03-node-typescript-project-layout/SKILL.md` — Node / TypeScript Project Layout
- `04-node-boundary-ownership/SKILL.md` — Node Boundary Ownership
- `05-typescript-domain-modeling/SKILL.md` — TypeScript Domain Modeling
- `06-typescript-domain-rules-and-transitions/SKILL.md` — TypeScript Domain Rules and Transitions
- `07-typescript-error-semantics/SKILL.md` — TypeScript Error Semantics
- `08-typescript-repository-ports/SKILL.md` — TypeScript Repository Ports
- `09-typescript-unit-of-work/SKILL.md` — TypeScript Unit of Work
- `10-typescript-cqrs-application-layer/SKILL.md` — TypeScript CQRS Application Layer
- `11-node-persistence-adapter/SKILL.md` — Node Persistence Adapter
- `12-database-schema-and-migration-parity/SKILL.md` — Database Schema and Migration Parity
- `13-node-bootstrap-composition-root/SKILL.md` — Node Bootstrap / Composition Root
- `14-node-http-transport/SKILL.md` — Node HTTP Transport
- `15-node-graphql-transport/SKILL.md` — Node GraphQL Transport
- `16-node-auth-boundary/SKILL.md` — Node Auth Boundary
- `17-node-contract-examples/SKILL.md` — Node Contract Examples
- `18-node-testing-and-conformance/SKILL.md` — Node Testing and Conformance
- `19-node-vertical-slice-generation/SKILL.md` — Node Vertical Slice Generation
- `20-node-architecture-review/SKILL.md` — Node Architecture Review

## Source summary
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
