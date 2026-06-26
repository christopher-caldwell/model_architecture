use chrono::Utc;
use domain::{
    book::{Book, BookCreationPayload, BookError},
    book_copy::{BookCopy, BookCopyCreationPayload, BookCopyError},
    uow::{WriteUnitOfWork, WriteUnitOfWorkFactory},
};
use std::sync::Arc;

#[derive(Clone)]
pub struct CatalogCommands {
    uow_factory: Arc<dyn WriteUnitOfWorkFactory>,
}

impl CatalogCommands {
    #[must_use]
    pub fn new(uow_factory: Arc<dyn WriteUnitOfWorkFactory>) -> Self {
        Self { uow_factory }
    }

    async fn get_book_by_isbn(
        &self,
        uow: &mut WriteUnitOfWork,
        isbn: &str,
    ) -> Result<Book, super::CommandError> {
        uow.book()
            .get_by_isbn(isbn)
            .await?
            .ok_or(BookError::NotFound.into())
    }

    async fn get_book_copy_by_barcode(
        &self,
        uow: &mut WriteUnitOfWork,
        barcode: &str,
    ) -> Result<BookCopy, super::CommandError> {
        uow.book_copy()
            .get_by_barcode_for_update(barcode)
            .await?
            .ok_or(BookCopyError::NotFound.into())
    }

    pub async fn add_book(
        &self,
        payload: BookCreationPayload,
    ) -> Result<Book, super::CommandError> {
        let prepared = payload.prepare();
        let mut uow = self.uow_factory.build().await?;
        let result = uow.book().create(&prepared).await?;
        uow.commit().await?;
        Ok(result)
    }

    pub async fn add_book_copy(
        &self,
        input: super::AddBookCopyInput,
    ) -> Result<BookCopy, super::CommandError> {
        let mut uow = self.uow_factory.build().await?;
        let book = self.get_book_by_isbn(&mut uow, &input.isbn).await?;
        let prepared = BookCopyCreationPayload {
            barcode: input.barcode,
            book_id: book.id,
        }
        .prepare();
        let result = uow.book_copy().create(&prepared).await?;
        uow.commit().await?;
        Ok(result)
    }

    pub async fn mark_book_copy_lost(
        &self,
        barcode: String,
    ) -> Result<BookCopy, super::CommandError> {
        let mut uow = self.uow_factory.build().await?;
        let book_copy = self.get_book_copy_by_barcode(&mut uow, &barcode).await?;
        let lost_status = book_copy.mark_lost()?;
        uow.book_copy()
            .update_status(book_copy.id, lost_status.clone())
            .await?;
        uow.commit().await?;
        let updated_copy = BookCopy {
            status: lost_status,
            dt_modified: Utc::now(),
            ..book_copy
        };
        Ok(updated_copy)
    }

    pub async fn mark_book_copy_found(
        &self,
        barcode: String,
    ) -> Result<BookCopy, super::CommandError> {
        let mut uow = self.uow_factory.build().await?;
        let book_copy = self.get_book_copy_by_barcode(&mut uow, &barcode).await?;
        let found_status = book_copy.mark_found()?;
        uow.book_copy()
            .update_status(book_copy.id, found_status.clone())
            .await?;
        uow.commit().await?;
        let updated_copy = BookCopy {
            status: found_status,
            dt_modified: Utc::now(),
            ..book_copy
        };
        Ok(updated_copy)
    }

    pub async fn send_book_copy_to_maintenance(
        &self,
        barcode: String,
    ) -> Result<BookCopy, super::CommandError> {
        let mut uow = self.uow_factory.build().await?;
        let book_copy = self.get_book_copy_by_barcode(&mut uow, &barcode).await?;
        let maintenance_status = book_copy.send_to_maintenance()?;
        uow.book_copy()
            .update_status(book_copy.id, maintenance_status.clone())
            .await?;
        uow.commit().await?;
        let updated_copy = BookCopy {
            status: maintenance_status,
            dt_modified: Utc::now(),
            ..book_copy
        };
        Ok(updated_copy)
    }

    pub async fn complete_book_copy_maintenance(
        &self,
        barcode: String,
    ) -> Result<BookCopy, super::CommandError> {
        let mut uow = self.uow_factory.build().await?;
        let book_copy = self.get_book_copy_by_barcode(&mut uow, &barcode).await?;
        let active_status = book_copy.complete_maintenance()?;
        uow.book_copy()
            .update_status(book_copy.id, active_status.clone())
            .await?;
        uow.commit().await?;
        let updated_copy = BookCopy {
            status: active_status,
            dt_modified: Utc::now(),
            ..book_copy
        };
        Ok(updated_copy)
    }
}
