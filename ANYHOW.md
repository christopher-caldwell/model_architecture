**Proposal**
I’d fix the `anyhow` boundary issue in two passes: first introduce typed infrastructure errors at the core ports, then add typed query errors at the application boundary. I would not try to model every database failure as business meaning. Most of these are still “unexpected infrastructure failure”; the win is making that explicit without exposing `anyhow`.

**1. Add a Domain Port Error**
Create something like `domain::port::PortError` or `domain::errors::InfrastructureError`:

```rust
#[derive(Debug, thiserror::Error)]
pub enum PortError {
    #[error("repository operation failed")]
    Repository {
        #[source]
        source: Box<dyn std::error::Error + Send + Sync>,
    },

    #[error("unit of work failed")]
    UnitOfWork {
        #[source]
        source: Box<dyn std::error::Error + Send + Sync>,
    },
}
```

Then define a local alias:

```rust
pub type PortResult<T> = Result<T, PortError>;
```

All domain repository and unit-of-work ports would return `PortResult<T>` instead of `anyhow::Result<T>`.

**2. Convert Persistence Errors At The Adapter Boundary**
Persistence would still use `anyhow::Context` internally if useful, but before returning through a domain port it would convert into `PortError`.

For example:

```rust
async fn get_by_ident(&self, ident: &MemberIdent) -> PortResult<Option<Member>> {
    get_member_by_ident_with(&self.pool, ident)
        .await
        .map_err(PortError::repository)
}
```

The persistence layer can keep detailed context, but callers only see a domain-owned infrastructure error type.

**3. Update CommandError**
`CommandError` already separates domain business errors from unexpected failures. I’d replace:

```rust
Unexpected(#[from] anyhow::Error)
```

with something like:

```rust
Infrastructure(#[from] domain::PortError)
```

or, if keeping command terminology:

```rust
UnexpectedInfrastructure(#[from] domain::PortError)
```

That keeps command errors typed and still lets transports map infrastructure failures to 500.

**4. Add QueryError**
Queries currently return `anyhow::Result<T>`. I’d add:

```rust
#[derive(Debug, thiserror::Error)]
pub enum QueryError {
    #[error(transparent)]
    Infrastructure(#[from] domain::PortError),
}
```

Then query methods return:

```rust
Result<Vec<Book>, QueryError>
```

This makes HTTP/GraphQL query error mapping explicit, matching the command side.

**5. Transport Mapping**
HTTP and GraphQL would replace `service_error(anyhow::Error)` / `gql_service_error(anyhow::Error)` with typed query/infrastructure mapping:

```rust
pub fn query_error(error: QueryError) -> ApiError {
    tracing::error!("Unhandled query error: {error:?}");
    internal_server_error()
}
```

Commands and queries remain distinct at the transport edge.

**What I Would Avoid**
I would not create `SqlxError` variants in domain. That would recreate the same leak with a different name.

I would not move repository ports out of domain just to avoid `anyhow`; this project intentionally keeps repository ports in domain.

I would not over-model every infrastructure failure. A small `PortError` is enough unless the application actually branches on specific infrastructure cases, which it should rarely do.
