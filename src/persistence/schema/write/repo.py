from collections.abc import Sequence

from typing_extensions import override

from domain.schema.entity import (
    NormalizedRegistrySchema,
    RuntimeSchema,
    RuntimeSchemaKey,
)
from domain.schema.port import RuntimeSchemaWriteRepoPort
from persistence.cache_port import RuntimeSchemaCacheWritePort


class RuntimeSchemaWriteRepoCache(RuntimeSchemaWriteRepoPort):
    def __init__(self, cache: RuntimeSchemaCacheWritePort) -> None:
        self._cache: RuntimeSchemaCacheWritePort = cache

    @override
    def put_many(self, schemas: Sequence[NormalizedRegistrySchema]) -> None:
        cacheable: list[tuple[RuntimeSchemaKey, RuntimeSchema]] = []
        for normalized in schemas:
            key: RuntimeSchemaKey = (
                f"{normalized.entity_name}@{normalized.entity_version}"
            )
            cacheable_item: tuple[RuntimeSchemaKey, RuntimeSchema] = (
                key,
                RuntimeSchema(
                    entity_name=normalized.entity_name,
                    entity_version=normalized.entity_version,
                    schema=normalized.schema,
                ),
            )
            cacheable.append(cacheable_item)
        self._cache.replace_all(cacheable)
