from fastapi import FastAPI, Request
from fastapi.exceptions import RequestValidationError
from fastapi.responses import JSONResponse

from api.customer.get_handlers import get_customer_handler
from api.customer.post_handlers import create_customer_handler
from api.invoice.get_handlers import get_invoice_handler
from api.invoice.post_handlers import create_invoice_handler
from api.order.get_handlers import get_order_handler
from api.order.post_handlers import create_order_handler
from api.product.get_handlers import get_product_handler
from api.product.post_handlers import create_product_handler
from server.deps import ServerDeps


def new_router(app: FastAPI, deps: ServerDeps) -> FastAPI:
    queries = deps.schema.queries
    validator = deps.schema.validator

    async def on_request_validation_error(
        _request: Request, exc: Exception
    ) -> JSONResponse:
        if not isinstance(exc, RequestValidationError):
            raise exc
        return JSONResponse(
            status_code=400,
            content={"success": False, "error": {"issues": exc.errors()}},
        )

    app.add_exception_handler(RequestValidationError, on_request_validation_error)

    app.include_router(get_product_handler(queries, validator))
    app.include_router(create_product_handler(queries, validator))

    app.include_router(get_order_handler(queries, validator))
    app.include_router(create_order_handler(queries, validator))

    app.include_router(get_customer_handler(queries, validator))
    app.include_router(create_customer_handler(queries, validator))

    app.include_router(get_invoice_handler(queries, validator))
    app.include_router(create_invoice_handler(queries, validator))

    return app
