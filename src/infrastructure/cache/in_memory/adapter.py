from collections.abc import Sequence

from typing_extensions import override

from domain.schema.entity import RuntimeSchema, RuntimeSchemaKey
from persistence.cache_port import (
    RuntimeSchemaCacheReadPort,
    RuntimeSchemaCacheWritePort,
)


class RuntimeSchemaCacheInMemory(
    RuntimeSchemaCacheReadPort, RuntimeSchemaCacheWritePort
):
    def __init__(self) -> None:
        self._cache: dict[RuntimeSchemaKey, RuntimeSchema] = {}

    @override
    def get(self, key: RuntimeSchemaKey) -> RuntimeSchema | None:
        return self._cache.get(key, None)

    @override
    def replace_all(
        self, schemas: Sequence[tuple[RuntimeSchemaKey, RuntimeSchema]]
    ) -> None:
        self._cache.clear()
        for key, value in schemas:
            self._cache[key] = value
