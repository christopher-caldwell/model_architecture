from dataclasses import dataclass


@dataclass(frozen=True)
class Customer:
    id: int
    ident: str
    email: str
    first_name: str
    last_name: str
    # @version 2.0.0
    phone: int
    country: str
    # @version 1.2.3
    dt_created: str
