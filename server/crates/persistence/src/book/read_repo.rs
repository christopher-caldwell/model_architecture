use anyhow::{Context, Result};
use async_trait::async_trait;
use chrono::{DateTime, Utc};
use domain::{
    book::{port::BookReadRepoPort, Book, BookId},
    PortError, PortResult,
};
use sqlx::PgPool;

#[derive(sqlx::FromRow)]
pub struct BookDbRow {
    pub book_id: i32,
    pub isbn: String,
    pub dt_created: DateTime<Utc>,
    pub dt_modified: DateTime<Utc>,
    pub title: String,
    pub author_name: String,
}

impl TryFrom<BookDbRow> for Book {
    type Error = anyhow::Error;

    fn try_from(value: BookDbRow) -> Result<Self> {
        Ok(Self {
            id: BookId(value.book_id),
            isbn: value.isbn,
            dt_created: value.dt_created,
            dt_modified: value.dt_modified,
            title: value.title,
            author_name: value.author_name,
        })
    }
}

pub struct BookReadRepoSql {
    pub pool: PgPool,
}

#[async_trait]
impl BookReadRepoPort for BookReadRepoSql {
    async fn get_catalog(&self) -> PortResult<Vec<Book>> {
        let rows = sqlx::query_file_as!(BookDbRow, "sql/book/queries/get_catalog.sql")
            .fetch_all(&self.pool)
            .await
            .context("Failed to fetch book catalog")
            .map_err(|error| PortError::repository(error.into_boxed_dyn_error()))?;

        rows.into_iter()
            .map(Book::try_from)
            .collect::<Result<Vec<_>>>()
            .map_err(|error| PortError::repository(error.into_boxed_dyn_error()))
    }

    async fn get_by_isbn(&self, isbn: &str) -> PortResult<Option<Book>> {
        let row = sqlx::query_file_as!(BookDbRow, "sql/book/queries/get_by_isbn.sql", isbn)
            .fetch_optional(&self.pool)
            .await
            .context("Failed to fetch book by isbn")
            .map_err(|error| PortError::repository(error.into_boxed_dyn_error()))?;

        row.map(Book::try_from)
            .transpose()
            .map_err(|error| PortError::repository(error.into_boxed_dyn_error()))
    }
}
