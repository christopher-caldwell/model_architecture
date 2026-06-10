use async_trait::async_trait;

use crate::{
    book::port::BookWriteRepoPort, book_copy::port::BookCopyWriteRepoPort,
    loan::port::LoanWriteRepoPort, member::port::MemberWriteRepoPort,
};

#[async_trait]
pub trait UnitOfWorkPort: Send {
    fn book(&mut self) -> &mut dyn BookWriteRepoPort;
    fn book_copy(&mut self) -> &mut dyn BookCopyWriteRepoPort;
    fn member(&mut self) -> &mut dyn MemberWriteRepoPort;
    fn loan(&mut self) -> &mut dyn LoanWriteRepoPort;
    async fn commit(self: Box<Self>) -> anyhow::Result<()>;
}

pub struct WriteUnitOfWork {
    inner: Box<dyn UnitOfWorkPort>,
}

impl WriteUnitOfWork {
    #[must_use]
    pub fn new(inner: Box<dyn UnitOfWorkPort>) -> Self {
        Self { inner }
    }

    pub fn book(&mut self) -> &mut dyn BookWriteRepoPort {
        self.inner.book()
    }

    pub fn book_copy(&mut self) -> &mut dyn BookCopyWriteRepoPort {
        self.inner.book_copy()
    }

    pub fn member(&mut self) -> &mut dyn MemberWriteRepoPort {
        self.inner.member()
    }

    pub fn loan(&mut self) -> &mut dyn LoanWriteRepoPort {
        self.inner.loan()
    }

    pub async fn commit(self) -> anyhow::Result<()> {
        self.inner.commit().await
    }
}

#[async_trait]
pub trait WriteUnitOfWorkFactory: Send + Sync {
    async fn build(&self) -> anyhow::Result<WriteUnitOfWork>;
}
