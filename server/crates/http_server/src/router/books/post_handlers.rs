use axum::{
    extract::{Path, State},
    http::StatusCode,
    Json,
};
use domain::book_copy::BookCopyCreationPayload;
use server_bootstrap::ServerDeps;

use crate::router::{
    auth::AuthUser,
    book_copies::schemas::{BookCopyResponseBody, CreateBookCopyRequestBody},
    books::schemas::{BookResponseBody, CreateBookRequestBody, BOOKS_TAG},
    errors::{not_found, service_error, ApiError},
};

#[utoipa::path(
    post,
    path = "",
    tag = BOOKS_TAG,
    request_body = CreateBookRequestBody,
    responses(
        (status = 201, description = "Book created", body = BookResponseBody),
        (status = 401, description = "Unauthorized"),
        (status = 403, description = "Forbidden"),
        (status = 500, description = "Internal server error", body = crate::router::errors::ErrorResponseBody)
    ),
    security(("bearer_auth" = []))
)]
pub async fn add_book(
    AuthUser(_claims): AuthUser,
    State(deps): State<ServerDeps>,
    Json(body): Json<CreateBookRequestBody>,
) -> Result<(StatusCode, Json<BookResponseBody>), ApiError> {
    let add_book_result = deps.catalog.commands.add_book(body.into()).await;

    let book_response = match add_book_result {
        Ok(book) => Json(BookResponseBody::from(book)),
        Err(error) => return Err(service_error(error)),
    };

    Ok((StatusCode::CREATED, book_response))
}

#[utoipa::path(
    post,
    path = "/{isbn}/copies",
    tag = BOOKS_TAG,
    params(
        ("isbn" = String, Path, description = "ISBN identifier for the book")
    ),
    request_body = CreateBookCopyRequestBody,
    responses(
        (status = 201, description = "Book copy created", body = BookCopyResponseBody),
        (status = 401, description = "Unauthorized"),
        (status = 403, description = "Forbidden"),
        (status = 404, description = "Book not found", body = crate::router::errors::ErrorResponseBody),
        (status = 500, description = "Internal server error", body = crate::router::errors::ErrorResponseBody)
    ),
    security(("bearer_auth" = []))
)]
pub async fn add_book_copy(
    AuthUser(_claims): AuthUser,
    State(deps): State<ServerDeps>,
    Path(isbn): Path<String>,
    Json(body): Json<CreateBookCopyRequestBody>,
) -> Result<(StatusCode, Json<BookCopyResponseBody>), ApiError> {
    let book_result = deps.catalog.queries.get_book_by_isbn(&isbn).await;

    let book = match book_result {
        Ok(Some(book)) => book,
        Ok(None) => return Err(not_found("Book not found")),
        Err(error) => return Err(service_error(error)),
    };

    let payload = BookCopyCreationPayload {
        barcode: body.barcode,
        author_name: body.author_name,
        book_id: book.id,
    };

    let add_book_copy_result = deps.catalog.commands.add_book_copy(payload).await;

    let book_copy_response = match add_book_copy_result {
        Ok(book_copy) => Json(BookCopyResponseBody::from(book_copy)),
        Err(error) => return Err(service_error(error)),
    };

    Ok((StatusCode::CREATED, book_copy_response))
}
