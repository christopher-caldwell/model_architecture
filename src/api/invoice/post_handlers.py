from fastapi import APIRouter
from fastapi.responses import JSONResponse

from api.openapi import ErrorResponse
from application.schema.commands.port import SchemaValidationAdapter
from application.schema.queries.port import SchemaRegistryQueries
from domain.invoice.entity import Invoice

from .mapper import from_invoice, to_invoice
from .schemas import INVOICES_PATH, INVOICES_TAG, InvoiceRequestBody, InvoiceResponse


def create_invoice_handler(
    queries: SchemaRegistryQueries,
    validator: SchemaValidationAdapter,
) -> APIRouter:
    router = APIRouter(tags=[INVOICES_TAG])

    @router.post(
        INVOICES_PATH,
        response_model=InvoiceResponse,
        status_code=201,
        responses={
            400: {"description": "Invalid request", "model": ErrorResponse},
            422: {"description": "Unknown schema version", "model": ErrorResponse},
            500: {"description": "Response validation failed", "model": ErrorResponse},
        },
    )
    def create_invoice(body: InvoiceRequestBody):  # pyright: ignore[reportUnusedFunction]
        entity_version = body.entity_version

        schema = queries.get_schema("invoice", entity_version)
        if schema is None:
            return JSONResponse(
                status_code=422,
                content={"error": f"unknown schema version: {entity_version}"},
            )

        invoice = to_invoice(validator, schema, body)
        if invoice is None:
            return JSONResponse(
                status_code=400,
                content={"error": "invalid payload"},
            )

        created_invoice_returned_from_create_command = Invoice(
            id=1,
            ident="invoice-001",
            order_id=invoice.order_id,
            customer_id=invoice.customer_id,
            amount=invoice.amount,
            currency=invoice.currency,
            status=invoice.status,
            issued_at=invoice.issued_at,
            dt_due=invoice.dt_due,
        )

        response = from_invoice(
            validator, schema, created_invoice_returned_from_create_command
        )
        if response is None:
            return JSONResponse(
                status_code=500,
                content={"error": "response failed schema validation"},
            )

        return response

    return router
