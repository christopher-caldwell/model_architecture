from abc import ABC, abstractmethod
from collections.abc import Sequence

from .entity import EntityName, EntityVersion, NormalizedRegistrySchema, RuntimeSchema


class RuntimeSchemaReadRepoPort(ABC):
    @abstractmethod
    def get(
        self, entity_name: EntityName, entity_version: EntityVersion
    ) -> RuntimeSchema | None: ...


class RuntimeSchemaWriteRepoPort(ABC):
    @abstractmethod
    def put_many(self, schemas: Sequence[NormalizedRegistrySchema]) -> None: ...
