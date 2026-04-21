from dataclasses import dataclass
from typing import Literal

OrderStatus = Literal["pending", "confirmed", "shipped", "delivered", "cancelled"]


@dataclass(frozen=True)
class Order:
    id: int
    ident: str
    customer_id: int
    status: OrderStatus
    total_amount: float
    currency: str
    # @version 1.2.3
    dt_created: str
    items_count: int
