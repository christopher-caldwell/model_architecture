use async_trait::async_trait;

use super::{BookCopy, BookCopyId, BookCopyPrepared, BookCopyStatus};

#[async_trait]
pub trait BookCopyWriteRepoPort: Send {
    async fn create(&mut self, insert: &BookCopyPrepared) -> anyhow::Result<BookCopy>;
    async fn get_by_barcode_for_update(
        &mut self,
        barcode: &str,
    ) -> anyhow::Result<Option<BookCopy>>;
    async fn update_status(&mut self, id: BookCopyId, status: BookCopyStatus)
        -> anyhow::Result<()>;
}

#[async_trait]
pub trait BookCopyReadRepoPort: Send + Sync {
    async fn get_by_id(&self, id: BookCopyId) -> anyhow::Result<Option<BookCopy>>;
    async fn get_by_barcode(&self, barcode: &str) -> anyhow::Result<Option<BookCopy>>;
}
