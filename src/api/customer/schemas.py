from pydantic import BaseModel

CUSTOMERS_TAG = "Customers"
CUSTOMERS_PATH = "/customers"


# Understand that normally these do not have id and ident, but since the proof of
# concep is only one schema per entity, we include them here.
class CustomerRequestBody(BaseModel):
    entity_version: str
    id: int
    ident: str
    email: str
    first_name: str
    last_name: str
    phone: str | int
    country: str
    created_at: str | None = None
    dt_created: str | None = None


class CustomerResponse(BaseModel):
    id: int
    ident: str
    email: str
    first_name: str
    last_name: str
    phone: int
    country: str
    dt_created: str
