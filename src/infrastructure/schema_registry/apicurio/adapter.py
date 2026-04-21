from collections.abc import Sequence
from dataclasses import dataclass

import requests
from typing_extensions import override

from application.schema.commands.port import SchemaRegistryProviderPort
from domain.schema.entity import (
    NormalizedRegistrySchema,
    SchemaDiscoveryConfig,
    SchemaDiscoveryTarget,
)
from infrastructure.json_io import is_json_array, is_json_object, response_json
from infrastructure.schema_registry.apicurio.types import (
    ApicurioVersionListResponse,
    ApicurioVersionMeta,
    ApicurioVersionState,
    ArtifactRef,
)
from infrastructure.schema_registry.apicurio.util import (
    assert_apicurio_response,
    build_version_content_path,
    build_version_list_path,
    format_artifact_ref,
    get_missing_artifact_error,
    select_effective_versions,
)


@dataclass(frozen=True)
class ApicurioSchemaRegistryProviderShape:
    endpoint_url: str
    group_id: str


class ApicurioSchemaRegistryProvider(SchemaRegistryProviderPort):
    def __init__(self, shape: ApicurioSchemaRegistryProviderShape) -> None:
        self._base_url: str = shape.endpoint_url.rstrip("/") + "/apis/registry/v3"
        self._group_id: str = shape.group_id

    @override
    def fetch_schemas(
        self, config: SchemaDiscoveryConfig
    ) -> Sequence[NormalizedRegistrySchema]:
        schemas: list[NormalizedRegistrySchema] = []
        for target in config.entities:
            schemas.extend(self._fetch_for_target(target))
        return schemas

    # -------------------------------------------------------------------------
    # Step orchestrator
    # -------------------------------------------------------------------------

    def _fetch_for_target(
        self, target: SchemaDiscoveryTarget
    ) -> list[NormalizedRegistrySchema]:
        ref = ArtifactRef(
            group_id=self._group_id,
            artifact_id=target.entity_name,
        )
        selected = select_effective_versions(
            self._load_version_metadata(ref),
            target.versions_to_load,
        )

        results: list[NormalizedRegistrySchema] = []
        for sv in selected:
            results.append(
                NormalizedRegistrySchema(
                    entity_name=target.entity_name,
                    entity_version=sv.semver,
                    schema=self._fetch_version_content(ref, sv.meta),
                )
            )
        return results

    # -------------------------------------------------------------------------
    # Step 2 — load all version metadata for an artifact (paginated)
    # -------------------------------------------------------------------------

    def _load_version_metadata(self, ref: ArtifactRef) -> list[ApicurioVersionMeta]:
        all_versions: list[ApicurioVersionMeta] = []
        offset = 0
        total_count: int | None = None

        while True:
            body = self._fetch_version_page(ref, offset)

            total_count = body.count
            all_versions.extend(body.versions)
            offset += len(body.versions)

            if not (
                len(all_versions) < (total_count or 0) and offset < (total_count or 0)
            ):
                break

        return all_versions

    def _fetch_version_content(
        self, ref: ArtifactRef, meta: ApicurioVersionMeta
    ) -> str:
        if meta.state == "DISABLED":
            raise RuntimeError(
                f'Artifact "{format_artifact_ref(ref)}" '
                f'version "{meta.version}" '
                f"(globalId={meta.global_id}) is in DISABLED state. "
                "Its content cannot be fetched and the startup load "
                "cannot continue."
            )

        return self._fetch_text(
            build_version_content_path(ref, meta.version),
            (
                f'Failed to fetch content for artifact "{format_artifact_ref(ref)}" '
                + f'version "{meta.version}" (globalId={meta.global_id})'
            ),
        )

    def _fetch_version_page(
        self, ref: ArtifactRef, offset: int
    ) -> ApicurioVersionListResponse:
        raw = self._fetch_json(
            build_version_list_path(ref, offset),
            f'Failed to list versions for artifact "{format_artifact_ref(ref)}"',
            get_missing_artifact_error(ref),
        )
        raw_versions = raw.get("versions")
        if not is_json_array(raw_versions):
            raise RuntimeError(
                "Unexpected response shape from Apicurio version list for "
                f'"{format_artifact_ref(ref)}": "versions" field is missing '
                "or not an array."
            )
        count = _require_int(raw, "count", ref)
        versions = [_parse_version_item(item, ref) for item in raw_versions]
        return ApicurioVersionListResponse(count=count, versions=versions)

    def _fetch_json(
        self,
        path: str,
        failure_prefix: str,
        not_found_message: str | None = None,
    ) -> dict[str, object]:
        response = requests.get(f"{self._base_url}{path}")
        assert_apicurio_response(response, failure_prefix, not_found_message)
        parsed = response_json(response)
        if not is_json_object(parsed):
            raise RuntimeError(
                f"{failure_prefix}: expected a JSON object response, "
                + f"got {type(parsed).__name__}."
            )
        return parsed

    def _fetch_text(
        self,
        path: str,
        failure_prefix: str,
        not_found_message: str | None = None,
    ) -> str:
        response = requests.get(f"{self._base_url}{path}")
        assert_apicurio_response(response, failure_prefix, not_found_message)
        return response.text


# ---------------------------------------------------------------------------
# Runtime validators for Apicurio version-list payloads
# ---------------------------------------------------------------------------


def _parse_version_item(raw: object, ref: ArtifactRef) -> ApicurioVersionMeta:
    if not is_json_object(raw):
        raise RuntimeError(
            f'Apicurio version list for "{format_artifact_ref(ref)}" contained '
            + "a non-object version item."
        )
    return ApicurioVersionMeta(
        version=_require_str(raw, "version", ref),
        global_id=_require_int(raw, "globalId", ref),
        content_id=_require_int(raw, "contentId", ref),
        state=_parse_version_state(raw.get("state"), ref),
        created_on=_require_str(raw, "createdOn", ref),
        group_id=raw.get("groupId")
        if isinstance(raw.get("groupId"), str)
        else ref.group_id,
        artifact_id=_require_str(raw, "artifactId", ref),
    )


def _parse_version_state(value: object, ref: ArtifactRef) -> ApicurioVersionState:
    if value == "ENABLED":
        return "ENABLED"
    if value == "DISABLED":
        return "DISABLED"
    if value == "DEPRECATED":
        return "DEPRECATED"
    raise RuntimeError(
        f"Apicurio returned unexpected version state {value!r} "
        + f'for "{format_artifact_ref(ref)}".'
    )


def _require_str(raw: dict[str, object], field: str, ref: ArtifactRef) -> str:
    value = raw.get(field)
    if not isinstance(value, str):
        raise RuntimeError(
            f'Apicurio response for "{format_artifact_ref(ref)}" is missing '
            + f'string field "{field}" (got {type(value).__name__}).'
        )
    return value


def _require_int(raw: dict[str, object], field: str, ref: ArtifactRef) -> int:
    value = raw.get(field)
    # bool is a subclass of int; reject it explicitly so we don't coerce True/False.
    if not isinstance(value, int) or isinstance(value, bool):
        raise RuntimeError(
            f'Apicurio response for "{format_artifact_ref(ref)}" is missing '
            + f'integer field "{field}" (got {type(value).__name__}).'
        )
    return value
