from fastapi import APIRouter
from fastapi.responses import JSONResponse

from api.openapi import ErrorResponse
from application.schema.commands.port import SchemaValidationAdapter
from application.schema.queries.port import SchemaRegistryQueries
from domain.product.entity import Product

from .mapper import from_product, to_product
from .schemas import (
    PRODUCTS_PATH,
    PRODUCTS_TAG,
    ProductRequestBody,
    ProductResponse,
)


def create_product_handler(
    queries: SchemaRegistryQueries,
    validator: SchemaValidationAdapter,
) -> APIRouter:
    router = APIRouter(tags=[PRODUCTS_TAG])

    @router.post(
        PRODUCTS_PATH,
        response_model=ProductResponse,
        status_code=201,
        responses={
            400: {"description": "Invalid request", "model": ErrorResponse},
            422: {"description": "Unknown schema version", "model": ErrorResponse},
            500: {"description": "Response validation failed", "model": ErrorResponse},
        },
    )
    def create_product(body: ProductRequestBody):  # pyright: ignore[reportUnusedFunction]
        entity_version = body.entity_version

        schema = queries.get_schema("product", entity_version)
        if schema is None:
            return JSONResponse(
                status_code=422,
                content={"error": f"unknown schema version: {entity_version}"},
            )

        product = to_product(validator, schema, body)
        if product is None:
            return JSONResponse(
                status_code=400,
                content={"error": "invalid payload"},
            )

        created_product_returned_from_create_command = Product(
            id=1,
            ident="product-001",
            name=product.name,
            description=product.description,
            price=product.price,
            currency=product.currency,
            category=product.category,
            sku=product.sku,
            stock_quantity=product.stock_quantity,
        )

        response = from_product(
            validator, schema, created_product_returned_from_create_command
        )
        if response is None:
            return JSONResponse(
                status_code=500,
                content={"error": "response failed schema validation"},
            )

        return response

    return router
