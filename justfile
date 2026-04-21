export PYTHONPYCACHEPREFIX := ".cache/pycache"
export RUFF_CACHE_DIR := ".cache/ruff"
export UV_CACHE_DIR := ".cache/uv"

dev:
    uv run uvicorn --app-dir src server:app --reload

check: typecheck lint

lint:
    uv run ruff check .

lint-fix:
    uv run ruff check . --fix

typecheck:
    uv run basedpyright .

clean:
    rm -rf .cache/pycache .cache/ruff .cache/uv
    find . -type d -name __pycache__ -exec rm -rf {} +
