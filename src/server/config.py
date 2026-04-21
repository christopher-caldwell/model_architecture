import os
from dataclasses import dataclass

from domain.schema.entity import SchemaDiscoveryConfig, SchemaDiscoveryTarget


@dataclass(frozen=True)
class ServerConfig:
    schema_registry_endpoint_url: str
    schema_registry_group_id: str
    schema_discovery_config: SchemaDiscoveryConfig


def required_env(name: str) -> str:
    value = os.environ.get(name)
    if value is None or len(value) == 0:
        raise ValueError(f"{name} must be set")
    return value


def load_server_config() -> ServerConfig:
    return ServerConfig(
        schema_registry_endpoint_url=required_env("SCHEMA_REGISTRY_ENDPOINT_URL"),
        schema_registry_group_id=os.environ.get("SCHEMA_REGISTRY_GROUP_ID")
        or "default",
        schema_discovery_config=SchemaDiscoveryConfig(
            entities=(
                SchemaDiscoveryTarget(entity_name="product", versions_to_load=5),
                SchemaDiscoveryTarget(entity_name="order", versions_to_load=5),
                SchemaDiscoveryTarget(entity_name="customer", versions_to_load=5),
                SchemaDiscoveryTarget(entity_name="invoice", versions_to_load=5),
            )
        ),
    )
