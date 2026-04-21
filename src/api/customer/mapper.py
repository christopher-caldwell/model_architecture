from application.schema.commands.port import SchemaValidationAdapter
from domain.customer.entity import Customer
from domain.schema.entity import RuntimeSchema

from .schemas import CustomerRequestBody, CustomerResponse


def to_customer(
    validator: SchemaValidationAdapter,
    schema: RuntimeSchema,
    body: CustomerRequestBody,
) -> Customer | None:
    payload: dict[str, object] = body.model_dump(exclude={"entity_version"}, exclude_none=True)
    if not validator.validate(schema, payload):
        print("Oopsie: failed validation")
        return None

    dt_created = body.dt_created or body.created_at
    if not dt_created:
        print("Oopsie: no dt_created")
        return None

    return Customer(
        id=0,
        ident="",
        email=body.email,
        first_name=body.first_name,
        last_name=body.last_name,
        phone=int(body.phone),
        country=body.country,
        dt_created=dt_created,
    )


def from_customer(
    validator: SchemaValidationAdapter,
    schema: RuntimeSchema,
    customer: Customer,
) -> CustomerResponse | None:
    response = CustomerResponse(
        id=customer.id,
        ident=customer.ident,
        email=customer.email,
        first_name=customer.first_name,
        last_name=customer.last_name,
        phone=customer.phone,
        country=customer.country,
        dt_created=customer.dt_created,
    )
    response_dict: dict[str, object] = response.model_dump()
    if not validator.validate(schema, response_dict):
        return None
    return response
