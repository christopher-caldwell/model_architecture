from dataclasses import dataclass
from typing import Literal

InvoiceStatus = Literal["draft", "issued", "paid", "overdue", "void"]


@dataclass(frozen=True)
class Invoice:
    id: int
    ident: str
    order_id: int
    customer_id: int
    # Decimal float (e.g. 10.99). Was plain integer in prior versions.
    # @version 1.0.0
    # Enforce decimal precision; relevant for typed languages like Go and Java.
    amount: float
    currency: str
    status: InvoiceStatus
    issued_at: str
    # @version 1.2.3
    dt_due: str
