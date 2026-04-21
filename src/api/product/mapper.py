from application.schema.commands.port import SchemaValidationAdapter
from domain.product.entity import Product
from domain.schema.entity import RuntimeSchema

from .schemas import ProductRequestBody, ProductResponse


def to_product(
    validator: SchemaValidationAdapter,
    schema: RuntimeSchema,
    body: ProductRequestBody,
) -> Product | None:
    payload: dict[str, object] = body.model_dump(exclude={"entity_version"}, exclude_none=True)
    if not validator.validate(schema, payload):
        return None

    return Product(
        id=0,
        ident="",
        name=body.name,
        description=body.description,
        price=body.price,
        currency=body.currency,
        category=body.category,
        sku=body.sku,
        stock_quantity=body.stock_quantity,
    )


def from_product(
    validator: SchemaValidationAdapter,
    schema: RuntimeSchema,
    product: Product,
) -> ProductResponse | None:
    response = ProductResponse(
        id=product.id,
        ident=product.ident,
        name=product.name,
        description=product.description,
        price=product.price,
        currency=product.currency,
        category=product.category,
        sku=product.sku,
        stock_quantity=product.stock_quantity,
    )
    response_dict: dict[str, object] = response.model_dump()
    if not validator.validate(schema, response_dict):
        return None
    return response
