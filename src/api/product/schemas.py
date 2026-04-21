from pydantic import BaseModel

PRODUCTS_TAG = "Products"
PRODUCTS_PATH = "/products"


class ProductRequestBody(BaseModel):
    entity_version: str
    id: int
    ident: str
    name: str
    description: str
    price: float
    currency: str
    category: str
    sku: str
    stock_quantity: int


class ProductResponse(BaseModel):
    id: int
    ident: str
    name: str
    description: str
    price: float
    currency: str
    category: str
    sku: str
    stock_quantity: int
