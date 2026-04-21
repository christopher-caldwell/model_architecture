from typing import Annotated

from fastapi import APIRouter, Query
from fastapi.responses import JSONResponse

from api._time import iso_now
from api.openapi import ErrorResponse
from application.schema.commands.port import SchemaValidationAdapter
from application.schema.queries.port import SchemaRegistryQueries
from domain.order.entity import Order

from .mapper import from_order
from .schemas import ORDERS_PATH, ORDERS_TAG, OrderResponse


def get_order_handler(
    queries: SchemaRegistryQueries,
    validator: SchemaValidationAdapter,
) -> APIRouter:
    router = APIRouter(tags=[ORDERS_TAG])

    @router.get(
        ORDERS_PATH,
        response_model=OrderResponse,
        responses={
            400: {"description": "Missing query param", "model": ErrorResponse},
            422: {"description": "Unknown schema version", "model": ErrorResponse},
            500: {"description": "Response validation failed", "model": ErrorResponse},
        },
    )
    def get_order(v: Annotated[str, Query(examples=["v1.0.0"])]):  # pyright: ignore[reportUnusedFunction]
        entity_version = v

        schema = queries.get_schema("order", entity_version)
        if schema is None:
            return JSONResponse(
                status_code=422,
                content={"error": f"unknown schema version: {entity_version}"},
            )

        fetched_order_from_query = Order(
            id=1,
            ident="order-001",
            customer_id=1,
            status="pending",
            total_amount=199.99,
            currency="USD",
            dt_created=iso_now(),
            items_count=3,
        )

        body = from_order(validator, schema, fetched_order_from_query)
        if body is None:
            return JSONResponse(
                status_code=500,
                content={"error": "response failed schema validation"},
            )

        return body

    return router
