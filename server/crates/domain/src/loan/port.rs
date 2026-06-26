use async_trait::async_trait;

use crate::{
    book_copy::BookCopyId,
    member::{MemberId, MemberIdent},
    PortResult,
};

use super::{Loan, LoanId, LoanPrepared};

#[async_trait]
pub trait LoanWriteRepoPort: Send {
    async fn create(&mut self, insert: &LoanPrepared) -> PortResult<Loan>;
    async fn end(&mut self, id: LoanId) -> PortResult<()>;
    async fn find_active_by_book_copy_id_for_update(
        &mut self,
        id: BookCopyId,
    ) -> PortResult<Option<Loan>>;
    async fn count_active_by_member_id(&mut self, id: MemberId) -> PortResult<i64>;
}

#[async_trait]
pub trait LoanReadRepoPort: Send + Sync {
    async fn get_by_member_ident(&self, ident: &MemberIdent) -> PortResult<Vec<Loan>>;
    async fn get_overdue(&self) -> PortResult<Vec<Loan>>;
    async fn find_active_by_book_copy_id(&self, id: BookCopyId) -> PortResult<Option<Loan>>;
    async fn count_active_by_member_id(&self, id: MemberId) -> PortResult<i64>;
}
