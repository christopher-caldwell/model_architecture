from dataclasses import dataclass
from typing import Literal

# ---------------------------------------------------------------------------
# Apicurio v3 API shapes (private to this module)
# ---------------------------------------------------------------------------

ApicurioVersionState = Literal["ENABLED", "DISABLED", "DEPRECATED"]


@dataclass(frozen=True)
class ApicurioVersionMeta:
    version: str
    global_id: int
    content_id: int
    state: ApicurioVersionState
    created_on: str
    group_id: str
    artifact_id: str


@dataclass(frozen=True)
class ApicurioVersionListResponse:
    count: int
    versions: list["ApicurioVersionMeta"]


# ---------------------------------------------------------------------------
# Artifact reference — internal coordinate pair used by all fetch helpers
# ---------------------------------------------------------------------------


@dataclass(frozen=True)
class ArtifactRef:
    group_id: str
    artifact_id: str


# ---------------------------------------------------------------------------
# Selected version — carries both raw metadata and its parsed semver
# ---------------------------------------------------------------------------


@dataclass(frozen=True)
class SelectedVersion:
    meta: ApicurioVersionMeta
    semver: str
