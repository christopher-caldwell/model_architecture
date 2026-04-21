from functools import cmp_to_key
from urllib.parse import quote

import requests
import semver

from infrastructure.schema_registry.apicurio.types import (
    ApicurioVersionMeta,
    ArtifactRef,
    SelectedVersion,
)


def select_effective_versions(
    metas: list[ApicurioVersionMeta],
    versions_to_load: int,
) -> list[SelectedVersion]:
    parsed: list[SelectedVersion] = [
        SelectedVersion(meta=meta, semver=parse_registry_semver(meta.version))
        for meta in metas
    ]

    parsed.sort(key=cmp_to_key(_cmp_versions))

    seen: dict[str, SelectedVersion] = {}
    for sv in parsed:
        if len(seen) >= versions_to_load:
            break
        if sv.semver not in seen:
            seen[sv.semver] = sv

    return list(seen.values())


def _compare_build(a: str, b: str) -> int:
    """Compare two semver strings including build metadata, mirroring npm semver.compareBuild."""  # noqa: E501
    va, vb = semver.Version.parse(a), semver.Version.parse(b)
    base = va.compare(vb)  # spec-correct: major/minor/patch/prerelease per-identifier
    if base != 0:
        return base
    # Build-metadata tiebreaker (npm extension — spec says ignore build):
    # no-build < has-build; else lexicographic comparison
    ab, bb = va.build or "", vb.build or ""
    if ab == bb:
        return 0
    if ab == "":
        return -1
    if bb == "":
        return 1
    return -1 if ab < bb else 1


def _cmp_versions(a: SelectedVersion, b: SelectedVersion) -> int:
    """Descending sort: higher semver first, then higher globalId as tiebreaker."""
    sv = _compare_build(b.semver, a.semver)
    if sv != 0:
        return sv
    return b.meta.global_id - a.meta.global_id


def parse_registry_semver(version: str) -> str:
    # semver.Version.parse is strict; clean leading 'v' or whitespace
    cleaned = version.strip().lstrip("v")
    try:
        v = semver.Version.parse(cleaned)
    except ValueError as err:
        raise ValueError(
            f'Schema registry returned an invalid semver: "{version}". '
            "Startup cannot continue."
        ) from err
    return str(v)


VERSION_PAGE_SIZE = 100


def build_version_list_path(ref: ArtifactRef, offset: int) -> str:
    return (
        f"{_build_artifact_path(ref)}/versions"
        f"?limit={VERSION_PAGE_SIZE}&offset={offset}"
    )


def build_version_content_path(ref: ArtifactRef, version: str) -> str:
    return f"{_build_artifact_path(ref)}/versions/{quote(version, safe='')}/content"


def _build_artifact_path(ref: ArtifactRef) -> str:
    group_id = quote(ref.group_id, safe="")
    artifact_id = quote(ref.artifact_id, safe="")
    return f"/groups/{group_id}/artifacts/{artifact_id}"


def format_artifact_ref(ref: ArtifactRef) -> str:
    return f"{ref.group_id}/{ref.artifact_id}"


def get_missing_artifact_error(ref: ArtifactRef) -> str:
    return (
        "Startup schema hydration failed because artifact "
        f'"{format_artifact_ref(ref)}" '
        f"was not found in Apicurio. "
        f'Check that the artifact exists in group "{ref.group_id}" '
        "and that the registry has been seeded. "
        "If your artifacts live in a different group, set "
        "SCHEMA_REGISTRY_GROUP_ID accordingly."
    )


def assert_apicurio_response(
    response: requests.Response,
    failure_prefix: str,
    not_found_message: str | None = None,
) -> None:
    if response.ok:
        return

    if response.status_code == 404 and not_found_message:
        raise RuntimeError(not_found_message)

    raise RuntimeError(
        f"{failure_prefix}: HTTP {response.status_code} {response.reason}"
    )
