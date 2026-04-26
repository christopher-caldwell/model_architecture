# Language-Neutral Spec Requirements

The language-neutral specification must be derived from inspected Rust code only. It must capture entities, value objects/IDs, commands, queries, repository ports, UoW semantics, transports, auth claims, persistence schema, fixtures/seeds if present, and tests. It must not contain Rust syntax as a requirement and must not choose TypeScript frameworks.

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
