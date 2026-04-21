from pydantic import BaseModel, Field


class ErrorResponse(BaseModel):
    error: str


class EntityVersionQuery(BaseModel):
    v: str = Field(examples=["v1.0.0"])
