use anyhow::{Context, Result};
use chrono::Utc;
use domain::book::{Book, BookId, BookPrepared};
use sqlx::PgExecutor;

use crate::book::read_repo::BookDbRow;

#[derive(sqlx::FromRow)]
pub struct BookCreateResult {
    pub book_id: i32,
}

pub(crate) async fn create_book_with(
    executor: impl PgExecutor<'_>,
    insert: &BookPrepared,
) -> Result<Book> {
    let created_book = sqlx::query_file_as!(
        BookCreateResult,
        "sql/book/commands/create.sql",
        insert.isbn,
        insert.title,
        insert.author_name,
    )
    .fetch_one(executor)
    .await
    .context("Failed to create book")?;

    let now = Utc::now();
    let created_book = Book {
        id: BookId(created_book.book_id),
        isbn: insert.isbn.clone(),
        dt_created: now,
        dt_modified: now,
        title: insert.title.clone(),
        author_name: insert.author_name.clone(),
    };
    Ok(created_book)
}

pub(crate) async fn get_book_by_isbn_for_write_with(
    executor: impl PgExecutor<'_>,
    isbn: &str,
) -> Result<Option<Book>> {
    let row = sqlx::query_file_as!(BookDbRow, "sql/book/commands/get_by_isbn.sql", isbn)
        .fetch_optional(executor)
        .await
        .context("Failed to fetch book by isbn")?;

    row.map(Book::try_from).transpose()
}
