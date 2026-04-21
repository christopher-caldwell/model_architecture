from dataclasses import dataclass

EntityName = str
EntityVersion = str
RuntimeSchemaKey = str
SchemaPayload = str


@dataclass(frozen=True)
class SchemaDiscoveryTarget:
    entity_name: EntityName
    versions_to_load: int


@dataclass(frozen=True)
class SchemaDiscoveryConfig:
    entities: tuple["SchemaDiscoveryTarget", ...]


@dataclass(frozen=True)
class NormalizedRegistrySchema:
    entity_name: EntityName
    entity_version: EntityVersion
    schema: SchemaPayload


@dataclass(frozen=True)
class RuntimeSchema:
    entity_name: EntityName
    entity_version: EntityVersion
    schema: SchemaPayload


RuntimeSchemaMap = dict[RuntimeSchemaKey, RuntimeSchema]
