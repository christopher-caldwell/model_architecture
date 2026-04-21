from dataclasses import dataclass


@dataclass(frozen=True)
class Product:
    id: int
    ident: str
    name: str
    description: str
    price: float
    currency: str
    category: str
    sku: str
    stock_quantity: int
