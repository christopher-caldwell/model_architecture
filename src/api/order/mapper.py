from application.schema.commands.port import SchemaValidationAdapter
from domain.order.entity import Order
from domain.schema.entity import RuntimeSchema

from .schemas import OrderRequestBody, OrderResponse


def to_order(
    validator: SchemaValidationAdapter,
    schema: RuntimeSchema,
    body: OrderRequestBody,
) -> Order | None:
    payload: dict[str, object] = body.model_dump(exclude={"entity_version"}, exclude_none=True)
    if not validator.validate(schema, payload):
        return None

    dt_created = body.dt_created or body.created_at
    if not dt_created:
        return None

    return Order(
        id=0,
        ident="",
        customer_id=body.customer_id,
        status=body.status,
        total_amount=body.total_amount,
        currency=body.currency,
        dt_created=dt_created,
        items_count=body.items_count,
    )


def from_order(
    validator: SchemaValidationAdapter,
    schema: RuntimeSchema,
    order: Order,
) -> OrderResponse | None:
    response = OrderResponse(
        id=order.id,
        ident=order.ident,
        customer_id=order.customer_id,
        status=order.status,
        total_amount=order.total_amount,
        currency=order.currency,
        dt_created=order.dt_created,
        items_count=order.items_count,
    )
    response_dict: dict[str, object] = response.model_dump()
    if not validator.validate(schema, response_dict):
        return None
    return response
