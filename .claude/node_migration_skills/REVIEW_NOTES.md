# Review Notes

## Repository verification
Canonical URL inspected: https://github.com/christopher-caldwell/model_architecture
No redirect was observed in the GitHub pages/API responses used. The inspected repository owner/name remained `christopher-caldwell/model_architecture`.

## Brainstorm pass
Initial required coverage mapped to twenty skills: inventory, language-neutral specification, layout, boundary ownership, domain modeling, domain transitions, errors, repository ports, unit of work, CQRS application layer, persistence, schema parity, bootstrap, HTTP, GraphQL, auth, contract examples, testing, vertical slices, and architecture review.

## Review pass 1: Rust parity
The skill list was checked against inspected Rust facts: crates, domain entities, commands, queries, UoW, HTTP paths, auth claims, and SQL schema. Database schema parity and transport DTO mapping were kept as explicit skills because they are independent failure points.

## Review pass 2: Node idiom fit
The skill list was refined to avoid Rust syntax copying. Traits become ports, Arc/dyn becomes constructor injection from the composition root, Result/thiserror/anyhow becomes typed domain/application errors plus edge mappers, and SQL row mapping stays inside persistence.

## Review pass 3: boundary and anti-hallucination
Each skill template now requires source inspection, a source-parity note, boundary checks, two per-skill review passes, and explicit labeling of placeholders. The skills avoid choosing Express, Fastify, Hono, Apollo, Yoga, Mercurius, pg, Kysely, Drizzle, Prisma, Zod, or Valibot as the architectural source of truth.

## Per-skill review summary
Each of the twenty SKILL.md files was drafted from one common reviewed template, then checked twice for: Rust repo parity, language-neutral architecture, Node idioms, anti-hallucination constraints, boundary ownership, and examples that could invent behavior. No fake examples such as unconfirmed Book shapes were added. Confirmed Book shape is included because it was inspected in the Rust repo.

## Inspected source evidence
- Domain Book entity/logic/ports/errors.
- Domain UnitOfWork port.
- Application catalog, membership, lending commands.
- Application catalog, membership, lending queries.
- HTTP router and schema modules for books, book copies, members, and loans.
- Auth claims.
- SQL DDL at database/sql/library/library/ddl.sql.

## Limitations
The skill pack is designed to guide a future implementation. It intentionally does not implement the Node application. It records confirmed facts where inspected and requires future users of the skills to re-inspect any additional Rust files before producing app-specific code.
