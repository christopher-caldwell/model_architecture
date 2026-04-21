from collections.abc import AsyncGenerator
from contextlib import asynccontextmanager

from dotenv import load_dotenv
from fastapi import FastAPI

from api.router import new_router
from server.config import load_server_config
from server.deps import build_server_deps

_ = load_dotenv()

_config = load_server_config()
_deps = build_server_deps(_config)


@asynccontextmanager
async def _lifespan(_app: FastAPI) -> AsyncGenerator[None]:
    _deps.schema.commands.load_schemas()
    yield


app = FastAPI(
    title="Schema Registry API",
    version="1.0.0",
    lifespan=_lifespan,
)

_ = new_router(app, _deps)
