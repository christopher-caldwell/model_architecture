use anyhow::{Context, Result};
use chrono::Utc;
use domain::book_copy::{BookCopy, BookCopyId, BookCopyPrepared, BookCopyStatus};
use sqlx::PgExecutor;

use crate::book_copy::read_repo::BookCopyDbRow;

#[derive(sqlx::FromRow)]
pub struct BookCopyCreateResult {
    pub book_copy_id: i32,
}

pub(crate) async fn create_book_copy_with(
    executor: impl PgExecutor<'_>,
    insert: &BookCopyPrepared,
) -> Result<BookCopy> {
    let created_book_copy = sqlx::query_file_as!(
        BookCopyCreateResult,
        "sql/book_copy/commands/create.sql",
        insert.book_id.0,
        insert.status.to_string(),
        insert.barcode,
    )
    .fetch_one(executor)
    .await
    .context("Failed to create book copy")?;

    let now = Utc::now();
    let result = BookCopy {
        id: BookCopyId(created_book_copy.book_copy_id),
        barcode: insert.barcode.clone(),
        dt_created: now,
        dt_modified: now,
        book_id: insert.book_id,
        status: insert.status.clone(),
    };
    Ok(result)
}

pub(crate) async fn get_book_copy_by_barcode_for_update_with(
    executor: impl PgExecutor<'_>,
    barcode: &str,
) -> Result<Option<BookCopy>> {
    let row = sqlx::query_file_as!(
        BookCopyDbRow,
        "sql/book_copy/commands/get_by_barcode_for_update.sql",
        barcode
    )
    .fetch_optional(executor)
    .await
    .context("Failed to fetch book copy by barcode")?;

    row.map(BookCopy::try_from).transpose()
}

pub(crate) async fn update_book_copy_status_with(
    executor: impl PgExecutor<'_>,
    book_copy_id: BookCopyId,
    status: BookCopyStatus,
) -> Result<()> {
    sqlx::query_file!(
        "sql/book_copy/commands/update_status.sql",
        book_copy_id.0,
        status.to_string(),
    )
    .execute(executor)
    .await
    .context("Failed to update book copy status")?;

    Ok(())
}
