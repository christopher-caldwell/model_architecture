from typing import Annotated

from fastapi import APIRouter, Query
from fastapi.responses import JSONResponse

from api._time import iso_now
from api.openapi import ErrorResponse
from application.schema.commands.port import SchemaValidationAdapter
from application.schema.queries.port import SchemaRegistryQueries
from domain.customer.entity import Customer

from .mapper import from_customer
from .schemas import CUSTOMERS_PATH, CUSTOMERS_TAG, CustomerResponse


def get_customer_handler(
    queries: SchemaRegistryQueries,
    validator: SchemaValidationAdapter,
) -> APIRouter:
    router = APIRouter(tags=[CUSTOMERS_TAG])

    @router.get(
        CUSTOMERS_PATH,
        response_model=CustomerResponse,
        responses={
            400: {"description": "Missing query param", "model": ErrorResponse},
            422: {"description": "Unknown schema version", "model": ErrorResponse},
            500: {"description": "Response validation failed", "model": ErrorResponse},
        },
    )
    def get_customer(v: Annotated[str, Query(examples=["v1.0.0"])]):  # pyright: ignore[reportUnusedFunction]
        entity_version = v

        schema = queries.get_schema("customer", entity_version)
        if schema is None:
            return JSONResponse(
                status_code=422,
                content={"error": f"unknown schema version: {entity_version}"},
            )

        fetched_customer_from_query = Customer(
            id=1,
            ident="customer-001",
            email="customer@example.com",
            first_name="Jane",
            last_name="Smith",
            phone=15550100,
            country="US",
            dt_created=iso_now(),
        )

        body = from_customer(validator, schema, fetched_customer_from_query)
        if body is None:
            return JSONResponse(
                status_code=500,
                content={"error": "response failed schema validation"},
            )

        return body

    return router
