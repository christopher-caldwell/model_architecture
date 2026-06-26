use domain::{
    book::{port::BookReadRepoPort, Book},
    book_copy::{port::BookCopyReadRepoPort, BookCopy},
};
use std::sync::Arc;

use super::QueryError;

#[derive(Clone)]
pub struct CatalogQueries {
    book_read_repo: Arc<dyn BookReadRepoPort>,
    book_copy_read_repo: Arc<dyn BookCopyReadRepoPort>,
}

impl CatalogQueries {
    #[must_use]
    pub fn new(
        book_read_repo: Arc<dyn BookReadRepoPort>,
        book_copy_read_repo: Arc<dyn BookCopyReadRepoPort>,
    ) -> Self {
        Self {
            book_read_repo,
            book_copy_read_repo,
        }
    }

    pub async fn get_book_catalog(&self) -> Result<Vec<Book>, QueryError> {
        self.book_read_repo
            .get_catalog()
            .await
            .map_err(QueryError::from)
    }

    pub async fn get_book_by_isbn(&self, isbn: &str) -> Result<Option<Book>, QueryError> {
        self.book_read_repo
            .get_by_isbn(isbn)
            .await
            .map_err(QueryError::from)
    }

    pub async fn get_book_copy_details(
        &self,
        barcode: &str,
    ) -> Result<Option<BookCopy>, QueryError> {
        self.book_copy_read_repo
            .get_by_barcode(barcode)
            .await
            .map_err(QueryError::from)
    }
}
