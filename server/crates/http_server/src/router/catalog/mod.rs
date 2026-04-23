pub mod book_copies;
pub mod books;

pub use book_copies::{
    add_book_copy, complete_book_copy_maintenance, get_book_copy_details, mark_book_copy_found,
    mark_book_copy_lost, report_lost_loaned_book_copy, return_book_copy,
    send_book_copy_to_maintenance, BOOK_COPIES_PATH, BOOK_COPY_BY_ID_PATH, BOOK_COPY_LOSS_PATH,
    BOOK_COPY_LOSS_REPORTS_PATH, BOOK_COPY_MAINTENANCE_PATH, BOOK_COPY_RETURNS_PATH,
};
pub use books::{add_book, get_book_catalog, BOOKS_PATH};

use utoipa::OpenApi;

#[derive(OpenApi)]
#[openapi(
    nest(
        (path = books::BOOKS_PATH, api = books::BooksApi),
        (path = book_copies::BOOK_COPIES_PATH, api = book_copies::BookCopiesApi)
    )
)]
pub struct CatalogApi;
