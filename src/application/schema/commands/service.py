from typing_extensions import override

from domain.schema import RuntimeSchemaWriteRepoPort, SchemaDiscoveryConfig

from .port import (
    SchemaRegistryCommands,
    SchemaRegistryProviderPort,
    SchemaValidationAdapter,
)


class SchemaCommands(SchemaRegistryCommands):
    def __init__(
        self,
        provider_adapter: SchemaRegistryProviderPort,
        validation_adapter: SchemaValidationAdapter,
        schema_runtime_write_repo: RuntimeSchemaWriteRepoPort,
        discovery_config: SchemaDiscoveryConfig,
    ) -> None:
        self._provider_adapter: SchemaRegistryProviderPort = provider_adapter
        self._validation_adapter: SchemaValidationAdapter = validation_adapter
        self._schema_runtime_write_repo: RuntimeSchemaWriteRepoPort = (
            schema_runtime_write_repo
        )
        self._discovery_config: SchemaDiscoveryConfig = discovery_config

    @override
    def load_schemas(self) -> None:
        normalized = self._provider_adapter.fetch_schemas(self._discovery_config)
        self._validation_adapter.hydrate_validation_cache(normalized)
        self._schema_runtime_write_repo.put_many(normalized)
