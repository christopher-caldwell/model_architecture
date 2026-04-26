# Plan: Replace Hand-Written SQL Constants With PgTyped

Status: planning only. Do not implement this migration as part of this task.

## Goal

Move persistence from ad hoc SQL string constants plus manual `pool.query<T>` / `client.query<T>` row typing to generated PgTyped query modules.

The architectural goal is not to introduce an ORM. Persistence should continue to be a thin adapter layer behind domain/application repository ports. SQL remains explicit, PostgreSQL-specific, and owned by `@library/persistence`; PgTyped only adds compile-time validation and generated parameter/result types.

PgTyped is a good fit because it generates TypeScript parameter and result types from raw SQL checked against a running PostgreSQL schema, supports SQL files, emits generated query modules, and executes through prepared query objects rather than string interpolation. References:

- https://pgtyped.dev/
- https://github.com/adelsz/pgtyped

## Current State

Current SQL ownership:

- `src/sql.ts` stores every SQL statement in a single nested `sql` object.
- `src/book.ts`, `src/book-copy.ts`, `src/member.ts`, and `src/loan.ts` import those strings and call `pool.query<T>` or `client.query<T>`.
- `src/mappers.ts` defines manual row interfaces (`BookRow`, `BookCopyRow`, `MemberRow`, `LoanRow`) and maps database rows into domain entities.
- `src/uow.ts` controls transactions manually with `BEGIN`, `COMMIT`, and `ROLLBACK`; that should remain hand-written because it is transaction orchestration, not query definition.

Current query groups:

- Book: `create`, `getByIsbn`, `getCatalog`
- Book copy: `create`, `getById`, `getByBarcode`, `getByBarcodeForUpdate`, `updateStatus`
- Member: `create`, `getById`, `getByIdent`, `getByIdentForUpdate`, `updateStatus`
- Loan: `create`, `end`, `findActiveByBookCopyId`, `findActiveByBookCopyIdForUpdate`, `countActiveByMemberId`, `getByMemberIdent`, `getOverdue`

## Non-Goals

- Do not move business logic into SQL or generated query wrappers.
- Do not change domain entities, repository ports, command/query services, HTTP DTOs, or GraphQL DTOs.
- Do not replace PostgreSQL, `pg`, or the existing transaction unit of work.
- Do not collapse read and write repository separation.
- Do not let generated types leak outside `@library/persistence`.

## Target Shape

Add SQL files under a dedicated query directory:

- `src/queries/book.sql`
- `src/queries/book-copy.sql`
- `src/queries/member.sql`
- `src/queries/loan.sql`

Generate query modules beside them or under a generated directory:

- `src/queries/book.queries.ts`
- `src/queries/book-copy.queries.ts`
- `src/queries/member.queries.ts`
- `src/queries/loan.queries.ts`

Repository classes should import generated prepared queries and call `.run(params, poolOrClient)`.

Example target repository shape:

```ts
import { getBookCatalog, getBookByIsbn } from "./queries/book.queries.js";

export class BookReadRepositoryPostgres implements BookReadRepository {
  constructor(private readonly pool: Pool) {}

  async getCatalog(): Promise<Book[]> {
    const rows = await getBookCatalog.run(undefined, this.pool);
    return rows.map(mapBook);
  }

  async getByIsbn(isbn: string): Promise<Book | null> {
    const rows = await getBookByIsbn.run({ isbn }, this.pool);
    return rows[0] === undefined ? null : mapBook(rows[0]);
  }
}
```

Exact generated call signatures may differ based on PgTyped's output for parameterless queries; verify against generated files before editing repositories.

## Migration Steps

1. Add PgTyped dependencies and scripts.

   Add `@pgtyped/cli` as a dev dependency and `@pgtyped/runtime` as a runtime dependency in `packages/persistence/package.json`. Add scripts:

   - `generate:queries`: run PgTyped once.
   - `watch:queries`: run PgTyped in watch mode for local SQL editing.
   - Consider making `typecheck` depend on generated files being present, but do not make it connect to the database unless CI has a predictable database service.

2. Add PgTyped config.

   Create a config file such as `packages/persistence/pgtyped.config.json`.

   Requirements:

   - Source directory should be package-local.
   - Include `src/queries/**/*.sql`.
   - Emit generated files with a stable template such as `{{dir}}/{{name}}.queries.ts`.
   - Read database connection settings from environment variables or an explicit local `DATABASE_URL`; do not hard-code credentials.

3. Split `src/sql.ts` into named SQL files.

   Move each existing SQL statement into the appropriate `.sql` file with PgTyped `@name` annotations.

   Naming should match repository intent, not database table names only:

   - `CreateBook`
   - `GetBookByIsbn`
   - `GetBookCatalog`
   - `CreateBookCopy`
   - `GetBookCopyById`
   - `GetBookCopyByBarcode`
   - `GetBookCopyByBarcodeForUpdate`
   - `UpdateBookCopyStatus`
   - `CreateMember`
   - `GetMemberById`
   - `GetMemberByIdent`
   - `GetMemberByIdentForUpdate`
   - `UpdateMemberStatus`
   - `CreateLoan`
   - `EndLoan`
   - `FindActiveLoanByBookCopyId`
   - `FindActiveLoanByBookCopyIdForUpdate`
   - `CountActiveLoansByMemberId`
   - `GetLoansByMemberIdent`
   - `GetOverdueLoans`

4. Convert positional placeholders to PgTyped named parameters.

   Replace `$1`, `$2`, etc. with named parameters that match Rust/DB field names where they are data fields:

   - `:isbn`
   - `:title`
   - `:author_name`
   - `:book_id`
   - `:status`
   - `:barcode`
   - `:member_id`
   - `:member_ident`

   Keep SQL column names unchanged. Keep domain-facing field names unchanged.

5. Generate query modules.

   Run the generation script against a database created from the library schema.

   Before repository refactors, inspect generated types for:

   - Nullable timestamp behavior for `NULLIF(... sentinel ...) AS dt_due` and `dt_returned`.
   - `COUNT(*)::BIGINT AS count`, which may generate `string` or `bigint` depending on PostgreSQL type parser/generator behavior.
   - Status columns from `struct_type.att_pub_ident`, which should likely remain `string` and continue to be parsed by `parseBookCopyStatus` / `parseMemberStatus`.

6. Refactor repositories one aggregate group at a time.

   Suggested order:

   - Book first because it has the smallest query surface.
   - Book copy next because it exercises joins, updates, and `FOR UPDATE`.
   - Member next because it mirrors book copy and validates status handling.
   - Loan last because it has CTEs, sentinel timestamps, count coercion, and more read/write paths.

   Keep each repository class implementing the same domain port. Generated PgTyped result interfaces should be package-private implementation detail.

7. Reduce manual row interfaces.

   After a repository imports generated query result types, either:

   - Update mappers to accept generated result types, or
   - Keep small local adapter types if generated names are too verbose.

   Preferred direction: keep `mapBook`, `mapBookCopy`, `mapMember`, and `mapLoan` as the only conversion points into domain entities. Do not spread domain object creation across generated query call sites.

8. Remove `src/sql.ts`.

   Delete `src/sql.ts` only after all repositories use generated query modules and typecheck passes.

9. Add verification to CI/developer workflow.

   Minimum verification after migration:

   - `pnpm --filter @library/persistence generate:queries`
   - `pnpm --filter @library/persistence typecheck`
   - Root `pnpm typecheck`
   - Root `pnpm test`
   - Root `pnpm lint:boundaries`

   If CI can run a PostgreSQL schema container, add a check that generated query files are current.

## Boundary Rules During Migration

- `@library/persistence` may depend on PgTyped packages.
- `@library/domain` and `@library/application` must not import PgTyped, generated query modules, PostgreSQL row types, or SQL files.
- HTTP and GraphQL transports must remain unaware of PgTyped.
- Repository methods remain the adapter boundary. They may use generated query types internally but must return domain entities or primitive domain values.

## Risk Areas

- PgTyped generation requires a live database schema. The migration needs a reliable way to spin up or point at the library schema during local development and CI.
- Generated timestamp types must line up with the current `pg` runtime behavior. The current mappers expect JavaScript `Date` values for `dt_created`, `dt_modified`, `dt_due`, and `dt_returned`.
- `COUNT(*)::BIGINT` currently maps through `string` and then `Number(row.count)`. Preserve that behavior unless the runtime parser is deliberately changed.
- `FOR UPDATE` queries must continue using transaction-scoped `PoolClient`, not the read-only pool.
- Create methods currently synthesize `dt_created` / `dt_modified` with `new Date()` after inserts that only return IDs. Decide whether to keep that behavior or change insert queries to return complete rows. Prefer returning complete rows if the database schema provides defaults, but only as a deliberate behavior-preserving refactor with tests.

## Acceptance Criteria

- SQL statements live in PgTyped-consumed SQL files with named queries.
- Repositories execute generated prepared query modules instead of importing `src/sql.ts`.
- Domain/application/transports do not import PgTyped artifacts.
- Existing tests pass.
- Import boundary checker passes.
- Generated files are either committed intentionally or reproducibly generated during build/CI; choose one policy and document it in this package.
