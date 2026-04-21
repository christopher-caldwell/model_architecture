from typing import Annotated

from fastapi import APIRouter, Query
from fastapi.responses import JSONResponse

from api.openapi import ErrorResponse
from application.schema.commands.port import SchemaValidationAdapter
from application.schema.queries.port import SchemaRegistryQueries
from domain.product.entity import Product

from .mapper import from_product
from .schemas import PRODUCTS_PATH, PRODUCTS_TAG, ProductResponse


def get_product_handler(
    queries: SchemaRegistryQueries,
    validator: SchemaValidationAdapter,
) -> APIRouter:
    router = APIRouter(tags=[PRODUCTS_TAG])

    @router.get(
        PRODUCTS_PATH,
        response_model=ProductResponse,
        responses={
            400: {"description": "Missing query param", "model": ErrorResponse},
            422: {"description": "Unknown schema version", "model": ErrorResponse},
            500: {"description": "Response validation failed", "model": ErrorResponse},
        },
    )
    def get_product(v: Annotated[str, Query(examples=["v1.0.0"])]):  # pyright: ignore[reportUnusedFunction]
        entity_version = v

        schema = queries.get_schema("product", entity_version)
        if schema is None:
            return JSONResponse(
                status_code=422,
                content={"error": f"unknown schema version: {entity_version}"},
            )

        fetched_product_from_query = Product(
            id=1,
            ident="product-001",
            name="Sample Product",
            description="A sample product description",
            price=29.99,
            currency="USD",
            category="General",
            sku="SKU-001",
            stock_quantity=100,
        )

        body = from_product(validator, schema, fetched_product_from_query)
        if body is None:
            return JSONResponse(
                status_code=500,
                content={"error": "response failed schema validation"},
            )

        return body

    return router
