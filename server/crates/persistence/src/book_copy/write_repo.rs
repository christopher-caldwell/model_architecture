use std::sync::Arc;

use anyhow::{Context, Result};
use async_trait::async_trait;
use chrono::Utc;
use domain::book_copy::{
    port::BookCopyWriteRepoPort, BookCopy, BookCopyId, BookCopyPrepared, BookCopyStatus,
};
use sqlx::{Postgres, Transaction};
use tokio::sync::Mutex;

use crate::book_copy::read_repo::BookCopyDbRow;
use crate::book_copy::book_copy_status_ident;

#[derive(sqlx::FromRow)]
pub struct BookCopyPreparedResult {
    pub book_copy_id: i32,
}

pub struct BookCopyWriteRepoTx {
    pub tx: Arc<Mutex<Option<Transaction<'static, Postgres>>>>,
}

#[async_trait]
impl BookCopyWriteRepoPort for BookCopyWriteRepoTx {
    async fn create(&self, insert: &BookCopyPrepared) -> Result<BookCopy> {
        let mut guard = self.tx.lock().await;
        let tx = guard.as_mut().context("Transaction already consumed")?;
        let prepared_result = sqlx::query_file_as!(
            BookCopyPreparedResult,
            "sql/book_copy/commands/create.sql",
            i32::from(insert.book_id.0),
            book_copy_status_ident(&insert.status),
            insert.barcode,
        )
        .fetch_one(&mut **tx)
        .await
        .context("Failed to create book copy")?;

        let now = Utc::now();
        Ok(BookCopy {
            id: BookCopyId(prepared_result.book_copy_id),
            barcode: insert.barcode.clone(),
            dt_created: now,
            dt_modified: now,
            book_id: insert.book_id,
            author_name: String::new(),
            status: insert.status.clone(),
        })
    }

    async fn get_by_barcode_for_update(&self, barcode: &str) -> Result<Option<BookCopy>> {
        let mut guard = self.tx.lock().await;
        let tx = guard.as_mut().context("Transaction already consumed")?;
        let row = sqlx::query_file_as!(
            BookCopyDbRow,
            "sql/book_copy/commands/get_by_barcode_for_update.sql",
            barcode
        )
        .fetch_optional(&mut **tx)
        .await
        .context("Failed to fetch book copy by barcode")?;

        row.map(BookCopy::try_from).transpose()
    }

    async fn update_status(&self, id: BookCopyId, status: BookCopyStatus) -> Result<()> {
        let mut guard = self.tx.lock().await;
        let tx = guard.as_mut().context("Transaction already consumed")?;
        let book_copy_id = i32::try_from(id.0).context("book_copy_id exceeds SQL integer range")?;
        sqlx::query_file!(
            "sql/book_copy/commands/update_status.sql",
            book_copy_id,
            book_copy_status_ident(&status),
        )
        .execute(&mut **tx)
        .await
        .context("Failed to update book copy status")?;

        Ok(())
    }
}
