from pydantic import BaseModel

from domain.order.entity import OrderStatus

ORDERS_TAG = "Orders"
ORDERS_PATH = "/orders"


class OrderRequestBody(BaseModel):
    entity_version: str
    id: int
    ident: str
    customer_id: int
    status: OrderStatus
    total_amount: float
    currency: str
    items_count: int
    created_at: str | None = None
    dt_created: str | None = None


class OrderResponse(BaseModel):
    id: int
    ident: str
    customer_id: int
    status: OrderStatus
    total_amount: float
    currency: str
    dt_created: str
    items_count: int
