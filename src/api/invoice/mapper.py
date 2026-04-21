from application.schema.commands.port import SchemaValidationAdapter
from domain.invoice.entity import Invoice
from domain.schema.entity import RuntimeSchema

from .schemas import InvoiceRequestBody, InvoiceResponse


def to_invoice(
    validator: SchemaValidationAdapter,
    schema: RuntimeSchema,
    body: InvoiceRequestBody,
) -> Invoice | None:
    payload: dict[str, object] = body.model_dump(exclude={"entity_version"}, exclude_none=True)
    if not validator.validate(schema, payload):
        return None

    dt_due = body.dt_due or body.due_at
    if not dt_due:
        return None

    return Invoice(
        id=0,
        ident="",
        order_id=body.order_id,
        customer_id=body.customer_id,
        amount=body.amount,
        currency=body.currency,
        status=body.status,
        issued_at=body.issued_at,
        dt_due=dt_due,
    )


def from_invoice(
    validator: SchemaValidationAdapter,
    schema: RuntimeSchema,
    invoice: Invoice,
) -> InvoiceResponse | None:
    response = InvoiceResponse(
        id=invoice.id,
        ident=invoice.ident,
        order_id=invoice.order_id,
        customer_id=invoice.customer_id,
        amount=invoice.amount,
        currency=invoice.currency,
        status=invoice.status,
        issued_at=invoice.issued_at,
        dt_due=invoice.dt_due,
    )
    response_dict: dict[str, object] = response.model_dump()
    if not validator.validate(schema, response_dict):
        return None
    return response
