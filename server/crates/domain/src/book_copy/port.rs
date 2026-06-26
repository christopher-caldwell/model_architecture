use async_trait::async_trait;

use crate::PortResult;

use super::{BookCopy, BookCopyId, BookCopyPrepared, BookCopyStatus};

#[async_trait]
pub trait BookCopyWriteRepoPort: Send {
    async fn create(&mut self, insert: &BookCopyPrepared) -> PortResult<BookCopy>;
    async fn get_by_barcode_for_update(&mut self, barcode: &str) -> PortResult<Option<BookCopy>>;
    async fn update_status(&mut self, id: BookCopyId, status: BookCopyStatus) -> PortResult<()>;
}

#[async_trait]
pub trait BookCopyReadRepoPort: Send + Sync {
    async fn get_by_id(&self, id: BookCopyId) -> PortResult<Option<BookCopy>>;
    async fn get_by_barcode(&self, barcode: &str) -> PortResult<Option<BookCopy>>;
}
