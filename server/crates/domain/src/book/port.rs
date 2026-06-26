use async_trait::async_trait;

use crate::PortResult;

use super::{Book, BookPrepared};

#[async_trait]
pub trait BookWriteRepoPort: Send {
    async fn create(&mut self, insert: &BookPrepared) -> PortResult<Book>;
    async fn get_by_isbn(&mut self, isbn: &str) -> PortResult<Option<Book>>;
}

#[async_trait]
pub trait BookReadRepoPort: Send + Sync {
    async fn get_catalog(&self) -> PortResult<Vec<Book>>;
    async fn get_by_isbn(&self, isbn: &str) -> PortResult<Option<Book>>;
}
