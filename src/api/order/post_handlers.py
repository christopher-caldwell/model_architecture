from fastapi import APIRouter
from fastapi.responses import JSONResponse

from api.openapi import ErrorResponse
from application.schema.commands.port import SchemaValidationAdapter
from application.schema.queries.port import SchemaRegistryQueries
from domain.order.entity import Order

from .mapper import from_order, to_order
from .schemas import (
    ORDERS_PATH,
    ORDERS_TAG,
    OrderRequestBody,
    OrderResponse,
)


def create_order_handler(
    queries: SchemaRegistryQueries,
    validator: SchemaValidationAdapter,
) -> APIRouter:
    router = APIRouter(tags=[ORDERS_TAG])

    @router.post(
        ORDERS_PATH,
        response_model=OrderResponse,
        status_code=201,
        responses={
            400: {"description": "Invalid request", "model": ErrorResponse},
            422: {"description": "Unknown schema version", "model": ErrorResponse},
            500: {"description": "Response validation failed", "model": ErrorResponse},
        },
    )
    def create_order(body: OrderRequestBody):  # pyright: ignore[reportUnusedFunction]
        entity_version = body.entity_version

        schema = queries.get_schema("order", entity_version)
        if schema is None:
            return JSONResponse(
                status_code=422,
                content={"error": f"unknown schema version: {entity_version}"},
            )

        order = to_order(validator, schema, body)
        if order is None:
            return JSONResponse(
                status_code=400,
                content={"error": "invalid payload"},
            )

        created_order_returned_from_create_command = Order(
            id=1,
            ident="order-001",
            customer_id=order.customer_id,
            status=order.status,
            total_amount=order.total_amount,
            currency=order.currency,
            dt_created=order.dt_created,
            items_count=order.items_count,
        )

        response = from_order(
            validator, schema, created_order_returned_from_create_command
        )
        if response is None:
            return JSONResponse(
                status_code=500,
                content={"error": "response failed schema validation"},
            )

        return response

    return router
