from dataclasses import dataclass

from application.schema.commands.port import (
    SchemaRegistryCommands,
    SchemaValidationAdapter,
)
from application.schema.commands.service import SchemaCommands
from application.schema.queries.port import SchemaRegistryQueries
from application.schema.queries.service import SchemaQueries
from infrastructure.cache.in_memory.adapter import RuntimeSchemaCacheInMemory
from infrastructure.schema_registry.apicurio.adapter import (
    ApicurioSchemaRegistryProvider,
    ApicurioSchemaRegistryProviderShape,
)
from infrastructure.schema_validation.jsonschema_impl.adapter import (
    JsonSchemaValidationAdapter,
)
from persistence.schema.read.repo import RuntimeSchemaReadRepoCache
from persistence.schema.write.repo import RuntimeSchemaWriteRepoCache

from .config import ServerConfig


@dataclass(frozen=True)
class SchemaDeps:
    commands: SchemaRegistryCommands
    queries: SchemaRegistryQueries
    validator: SchemaValidationAdapter


@dataclass(frozen=True)
class ServerDeps:
    schema: SchemaDeps


def build_server_deps(config: ServerConfig) -> ServerDeps:
    # Infrastructure — active provider adapter (Apicurio for v0)
    provider_adapter = ApicurioSchemaRegistryProvider(
        ApicurioSchemaRegistryProviderShape(
            endpoint_url=config.schema_registry_endpoint_url,
            group_id=config.schema_registry_group_id,
        )
    )
    validation_adapter = JsonSchemaValidationAdapter()
    schema_cache = RuntimeSchemaCacheInMemory()
    schema_read_repo = RuntimeSchemaReadRepoCache(schema_cache)
    schema_write_repo = RuntimeSchemaWriteRepoCache(schema_cache)

    # Application — commands write to the store; queries read from it
    schema_commands = SchemaCommands(
        provider_adapter=provider_adapter,
        validation_adapter=validation_adapter,
        schema_runtime_write_repo=schema_write_repo,
        discovery_config=config.schema_discovery_config,
    )
    schema_queries = SchemaQueries(schema_read_repo)

    return ServerDeps(
        schema=SchemaDeps(
            commands=schema_commands,
            queries=schema_queries,
            validator=validation_adapter,
        )
    )
