from abc import ABC, abstractmethod
from collections.abc import Sequence

from domain.schema.entity import RuntimeSchema, RuntimeSchemaKey


class RuntimeSchemaCacheReadPort(ABC):
    @abstractmethod
    def get(self, key: RuntimeSchemaKey) -> RuntimeSchema | None: ...


class RuntimeSchemaCacheWritePort(ABC):
    @abstractmethod
    def replace_all(
        self, schemas: Sequence[tuple[RuntimeSchemaKey, RuntimeSchema]]
    ) -> None: ...
