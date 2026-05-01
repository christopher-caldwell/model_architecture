# Agent Guide

This repository is a clean/onion architecture + CQRS reference implementation. The application core is meant to be stable while outer layers are replaceable; HTTP and GraphQL already expose the same internal commands and queries.

Dependencies point inward:

```text
transport -> server_bootstrap -> application -> domain
                       persistence -> domain
```

## Always-Relevant Rules

- Business workflow belongs in `application` commands and queries, not in HTTP handlers or GraphQL resolvers.
- Reusable business decisions and invariants belong in `domain` methods and typed domain errors.
- Transports parse input, call one command or query for the use case, map output, and map errors.
- Do not compose business workflows in transports. Query-then-branch-then-command in a handler/resolver is wrong.
- Writes must not read after write just to hydrate a response.
- Write SQL returns only generated values the application cannot know, such as IDs or generated identifiers. Do not return database timestamps for hydration.
- Created entities use one code-side `Utc::now()` for both `dt_created` and `dt_modified` in the returned domain value.
- Ports stay abstract and narrow; implementations belong in adapters.

## Skills To Read Before Changes

- `.agents/skills/onion-cqrs-project-map/SKILL.md`: layer map, crate roles, feature flow, and review checklist.
- `.agents/skills/domain-business-encapsulation/SKILL.md`: exact domain encapsulation style, guard/transition shape, and enum handling.
- `.agents/skills/domain-rules-and-errors/SKILL.md`: business decisions, guards, transitions, defaults, and domain errors.
- `.agents/skills/application-cqrs-use-cases/SKILL.md`: commands, queries, orchestration, transactions, and unit-of-work behavior.
- `.agents/skills/persistence-write-results/SKILL.md`: repositories, SQL command files, generated IDs, timestamps, and no-read-after-write rules.
- `.agents/skills/transport-adapters/SKILL.md`: HTTP, GraphQL, schemas, request parsing, response DTOs, and error mapping.
