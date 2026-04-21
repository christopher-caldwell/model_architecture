from typing_extensions import override

from domain.schema import (
    EntityName,
    EntityVersion,
    RuntimeSchema,
    RuntimeSchemaReadRepoPort,
)

from .port import SchemaRegistryQueries


class SchemaQueries(SchemaRegistryQueries):
    def __init__(self, read_repo: RuntimeSchemaReadRepoPort) -> None:
        self._read_repo: RuntimeSchemaReadRepoPort = read_repo

    @override
    def get_schema(
        self, entity_name: EntityName, entity_version: EntityVersion
    ) -> RuntimeSchema | None:
        return self._read_repo.get(entity_name, entity_version)
