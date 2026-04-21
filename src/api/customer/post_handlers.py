from fastapi import APIRouter
from fastapi.responses import JSONResponse

from api._time import iso_now
from api.openapi import ErrorResponse
from application.schema.commands.port import SchemaValidationAdapter
from application.schema.queries.port import SchemaRegistryQueries
from domain.customer.entity import Customer

from .mapper import from_customer, to_customer
from .schemas import (
    CUSTOMERS_PATH,
    CUSTOMERS_TAG,
    CustomerRequestBody,
    CustomerResponse,
)


def create_customer_handler(
    queries: SchemaRegistryQueries,
    validator: SchemaValidationAdapter,
) -> APIRouter:
    router = APIRouter(tags=[CUSTOMERS_TAG])

    @router.post(
        CUSTOMERS_PATH,
        response_model=CustomerResponse,
        status_code=201,
        responses={
            400: {"description": "Invalid request", "model": ErrorResponse},
            422: {"description": "Unknown schema version", "model": ErrorResponse},
            500: {"description": "Response validation failed", "model": ErrorResponse},
        },
    )
    def create_customer(body: CustomerRequestBody):  # pyright: ignore[reportUnusedFunction]
        entity_version = body.entity_version

        schema = queries.get_schema("customer", entity_version)
        if schema is None:
            return JSONResponse(
                status_code=422,
                content={"error": f"unknown schema version: {entity_version}"},
            )

        customer = to_customer(validator, schema, body)
        if customer is None:
            return JSONResponse(
                status_code=400,
                content={"error": "invalid payload"},
            )

        created_customer_returned_from_create_command = Customer(
            id=1,
            ident="customer-001",
            email=customer.email,
            first_name=customer.first_name,
            last_name=customer.last_name,
            phone=customer.phone,
            country=customer.country,
            dt_created=iso_now(),
        )

        response = from_customer(
            validator, schema, created_customer_returned_from_create_command
        )
        if response is None:
            return JSONResponse(
                status_code=500,
                content={"error": "response failed schema validation"},
            )

        return response

    return router
