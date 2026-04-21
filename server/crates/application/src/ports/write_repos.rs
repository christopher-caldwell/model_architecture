use async_trait::async_trait;
use domain::{book::{BookPrepared, Book}, book_copy::{BookCopyPrepared, BookCopy}, loan::{Loan, LoanPrepared}, member::{Member, MemberPrepared}};

#[async_trait]
pub trait BookWriteRepoPort: Send + Sync {
    async fn create(
        &self,
        insert: &BookPrepared,
    ) -> anyhow::Result<Book>;
}

#[async_trait]
pub trait BookCopyWriteRepoPort: Send + Sync {
    async fn create(
        &self,
        insert: &BookCopyPrepared,
    ) -> anyhow::Result<BookCopy>;
    async fn update_status(
        &self,
        book_copy_id: i64,
        status: &str
    ) -> anyhow::Result<BookCopy>;
}

#[async_trait]
pub trait MemberWriteRepoPort: Send + Sync {
    async fn create(
        &self,
        insert: &MemberPrepared,
    ) -> anyhow::Result<Member>;
    async fn update_status(
        &self,
        member_id: i16,
        status: &str,
    ) -> anyhow::Result<Member>;
}

#[async_trait]
pub trait LoanWriteRepoPort: Send + Sync {
    async fn create(
        &self,
        insert: &LoanPrepared,
    ) -> anyhow::Result<Loan>;
    async fn end(
        &self,
        loan_id: i64,
    ) -> anyhow::Result<Loan>;
}
