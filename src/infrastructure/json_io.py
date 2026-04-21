"""Typed wrappers for untyped JSON boundaries.

`json.loads`, `requests.Response.json()`, and similar parsers are typed as
returning `Any` in their stubs. Calling them directly taints downstream code
with `Any` (and, under basedpyright's strict mode, trips `reportAny`). The
helpers here contain that `Any` to a single line and expose a properly typed
`object` to the rest of the codebase.
"""

import json

import requests
from typing_extensions import TypeIs


def loads_json(text: str | bytes) -> object:
    return json.loads(text)  # pyright: ignore[reportAny]


def response_json(response: requests.Response) -> object:
    return response.json()  # pyright: ignore[reportAny]


def is_json_object(value: object) -> TypeIs[dict[str, object]]:
    # JSON objects always have string keys per spec; this guard narrows the
    # decoded-but-untyped value to a properly typed mapping without a `cast`.
    return isinstance(value, dict) and all(
        isinstance(k, str) for k in value  # pyright: ignore[reportUnknownVariableType]
    )


def is_json_array(value: object) -> TypeIs[list[object]]:
    # Every Python value is an `object`, so promoting `list` (with unknown
    # element type) to `list[object]` is always sound for decoded JSON.
    return isinstance(value, list)
