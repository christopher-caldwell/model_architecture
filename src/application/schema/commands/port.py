import abc
from collections.abc import Sequence

from domain.schema import NormalizedRegistrySchema, RuntimeSchema, SchemaDiscoveryConfig


class SchemaRegistryProviderPort(abc.ABC):
    @abc.abstractmethod
    def fetch_schemas(
        self, config: SchemaDiscoveryConfig
    ) -> Sequence[NormalizedRegistrySchema]:
        pass


class SchemaValidationAdapter(abc.ABC):
    @abc.abstractmethod
    def hydrate_validation_cache(
        self, schemas: Sequence[NormalizedRegistrySchema]
    ) -> None:
        pass

    @abc.abstractmethod
    def validate(
        self, runtime_schema: RuntimeSchema, payload: dict[str, object]
    ) -> bool:
        pass


class SchemaRegistryCommands(abc.ABC):
    @abc.abstractmethod
    def load_schemas(self) -> None:
        pass
