from typing import Literal

from pydantic import BaseModel

INVOICES_TAG = "Invoices"
INVOICES_PATH = "/invoices"

InvoiceStatusValue = Literal["draft", "issued", "paid", "overdue", "void"]


class InvoiceRequestBody(BaseModel):
    entity_version: str
    id: int
    ident: str
    order_id: int
    customer_id: int
    amount: float
    currency: str
    status: InvoiceStatusValue
    issued_at: str
    due_at: str | None = None
    dt_due: str | None = None


class InvoiceResponse(BaseModel):
    id: int
    ident: str
    order_id: int
    customer_id: int
    amount: float
    currency: str
    status: InvoiceStatusValue
    issued_at: str
    dt_due: str
