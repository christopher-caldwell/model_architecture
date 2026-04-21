import abc

from domain.schema import EntityName, EntityVersion, RuntimeSchema


class SchemaRegistryQueries(abc.ABC):
    @abc.abstractmethod
    def get_schema(
        self, entity_name: EntityName, entity_version: EntityVersion
    ) -> RuntimeSchema | None:
        pass
