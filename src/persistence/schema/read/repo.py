from typing_extensions import override

from domain.schema.entity import (
    EntityName,
    EntityVersion,
    RuntimeSchema,
    RuntimeSchemaKey,
)
from domain.schema.port import RuntimeSchemaReadRepoPort
from persistence.cache_port import RuntimeSchemaCacheReadPort


class RuntimeSchemaReadRepoCache(RuntimeSchemaReadRepoPort):
    def __init__(self, cache: RuntimeSchemaCacheReadPort) -> None:
        self._cache: RuntimeSchemaCacheReadPort = cache

    @override
    def get(
        self, entity_name: EntityName, entity_version: EntityVersion
    ) -> RuntimeSchema | None:
        key: RuntimeSchemaKey = f"{entity_name}@{entity_version}"
        return self._cache.get(key)
