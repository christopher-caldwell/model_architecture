from datetime import datetime, timedelta, timezone
from typing import Annotated

from fastapi import APIRouter, Query
from fastapi.responses import JSONResponse

from api._time import iso_format, iso_now
from api.openapi import ErrorResponse
from application.schema.commands.port import SchemaValidationAdapter
from application.schema.queries.port import SchemaRegistryQueries
from domain.invoice.entity import Invoice

from .mapper import from_invoice
from .schemas import INVOICES_PATH, INVOICES_TAG, InvoiceResponse


def get_invoice_handler(
    queries: SchemaRegistryQueries,
    validator: SchemaValidationAdapter,
) -> APIRouter:
    router = APIRouter(tags=[INVOICES_TAG])

    @router.get(
        INVOICES_PATH,
        response_model=InvoiceResponse,
        responses={
            400: {"description": "Missing query param", "model": ErrorResponse},
            422: {"description": "Unknown schema version", "model": ErrorResponse},
            500: {"description": "Response validation failed", "model": ErrorResponse},
        },
    )
    def get_invoice(v: Annotated[str, Query(examples=["v1.0.0"])]):  # pyright: ignore[reportUnusedFunction]
        entity_version = v

        schema = queries.get_schema("invoice", entity_version)
        if schema is None:
            return JSONResponse(
                status_code=422,
                content={"error": f"unknown schema version: {entity_version}"},
            )

        now = datetime.now(timezone.utc)
        fetched_invoice_from_query = Invoice(
            id=1,
            ident="invoice-001",
            order_id=1,
            customer_id=1,
            amount=199.99,
            currency="USD",
            status="issued",
            issued_at=iso_now(),
            dt_due=iso_format(now + timedelta(days=30)),
        )

        body = from_invoice(validator, schema, fetched_invoice_from_query)
        if body is None:
            return JSONResponse(
                status_code=500,
                content={"error": "response failed schema validation"},
            )

        return body

    return router
