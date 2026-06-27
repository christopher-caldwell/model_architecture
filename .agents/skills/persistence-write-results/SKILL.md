---
name: persistence-write-results
description: Use when touching SQL command files, write repositories, generated IDs, code-side timestamps, SQLx row structs, read/write repository split, status lookup mapping, command-side SQL, or no-read-after-write behavior in this onion/CQRS Rust project.
---

# Persistence Write Results

Use this skill when touching SQL command files, write repositories, generated IDs, timestamps, SQLx row structs, or no-read-after-write behavior. For the transaction-backed UoW adapter shape, also read `.agents/skills/unit-of-work-cqrs-pattern/SKILL.md`.

## Purpose

Persistence is an adapter. It executes SQL and maps between database rows and domain types, but it does not own business decisions.

## Read/Write Split

- Read repositories use pools and SQL under `sql/<concept>/queries`.
- Write repositories use transactions and SQL under `sql/<concept>/commands`.
- Command-side read SQL may exist under `commands` when it is part of a write decision, especially if it uses `FOR UPDATE`.
- Shared SQL helpers may accept a generic SQLx executor when the same SQL semantics work with either a pool or a transaction.

## Write SQL Return Rule

Write SQL should return only generated values the application cannot know:

- Generated primary keys such as `book_id`, `book_copy_id`, `member_id`.
- Generated business identifiers such as `loan_ident` when the database constructs them.

Do not return database timestamps from insert/update SQL just to hydrate response entities.

Good:

```sql
INSERT INTO library.book (...)
VALUES (...)
RETURNING book_id;
```

Avoid:

```sql
RETURNING book_id, dt_created, dt_modified, ...
```

## Create Mapping Rule

For created entities:

1. Execute insert.
2. Read only generated values.
3. Capture one `Utc::now()`.
4. Build the domain entity from generated values plus the prepared input.
5. Set both `dt_created` and `dt_modified` to that `now`.

This keeps writes from depending on a post-write read.

## Update Mapping Rule

For update methods:

- Execute the SQL update.
- Return `()` from the repository unless a generated value is truly needed.
- Let the command construct the returned entity from the locked pre-update entity, domain-approved new value, and `Utc::now()`.

## Locking Reads

If a command makes a decision based on a row it will update or coordinate with a write, use a write repository method backed by command-side SQL with `FOR UPDATE`.

Examples:

- `get_by_ident_for_update`
- `get_by_barcode_for_update`
- `find_active_by_book_copy_id_for_update`

## Status/Dictionary Values

Domain code should use enums. SQL adapters may map lookup-table strings or IDs into those enums. Keep conversion at the adapter boundary and fail with context if the database contains an unknown value.

## Smells

- `UPDATE ... RETURNING *` for response hydration.
- `INSERT ... RETURNING dt_created, dt_modified`.
- A write repository checks whether a business transition is allowed.
- A command performs a select after write only to return current state.
- Read repository methods are used inside a command that needs write consistency.
