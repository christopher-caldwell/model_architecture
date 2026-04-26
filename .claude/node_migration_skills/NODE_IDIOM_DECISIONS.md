# Node Idiom Decisions

Use strict TypeScript, async/await, explicit package/module boundaries, constructor injection, typed errors or discriminated error mappings, explicit DTO/domain/row mappers, and migration tooling that preserves the Rust SQL schema. Frameworks such as Express/Fastify/Hono/Apollo/Yoga/Mercurius may be discussed but must not own the core. Validation libraries may exist only at the edge.
