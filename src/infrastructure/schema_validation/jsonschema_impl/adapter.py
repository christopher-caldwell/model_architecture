from collections.abc import Sequence

import jsonschema
from typing_extensions import override

from application.schema.commands.port import SchemaValidationAdapter
from domain.schema.entity import NormalizedRegistrySchema, RuntimeSchema
from infrastructure.json_io import is_json_object, loads_json


class JsonSchemaValidationAdapter(SchemaValidationAdapter):
    def __init__(self) -> None:
        self._cache: dict[str, dict[str, object]] = {}
        self._format_checker: jsonschema.FormatChecker = jsonschema.FormatChecker()

    @override
    def validate(
        self, runtime_schema: RuntimeSchema, payload: dict[str, object]
    ) -> bool:
        key = self._key_of(runtime_schema)
        compiled = self._cache.get(key)
        if compiled is None:
            compiled = self._parse_for_compilation(runtime_schema)
            self._cache[key] = compiled
        try:
            jsonschema.validate(
                instance=payload,
                schema=compiled,
                format_checker=self._format_checker,
            )
            return True
        except jsonschema.ValidationError:
            return False

    @override
    def hydrate_validation_cache(
        self, schemas: Sequence[NormalizedRegistrySchema]
    ) -> None:
        next_cache: dict[str, dict[str, object]] = {}
        for schema in schemas:
            compiled = self._parse_for_compilation(schema)
            next_cache[self._key_of(schema)] = compiled
        self._cache.clear()
        self._cache.update(next_cache)

    def _key_of(self, schema: RuntimeSchema | NormalizedRegistrySchema) -> str:
        return f"{schema.entity_name}@{schema.entity_version}"

    def _parse_for_compilation(
        self, schema: RuntimeSchema | NormalizedRegistrySchema
    ) -> dict[str, object]:
        parsed = loads_json(schema.schema)
        if not is_json_object(parsed):
            raise ValueError(
                f"Schema {self._key_of(schema)} must be a JSON object, "
                + f"got {type(parsed).__name__}."
            )

        # Upstream schema files reuse bare version numbers as $id values across
        # different entities, which collides inside jsonschema's resolution.
        # Namespace them per runtime key before compilation.
        if isinstance(parsed.get("$id"), str):
            parsed["$id"] = self._key_of(schema)

        return parsed
