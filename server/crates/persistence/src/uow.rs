use anyhow::Context;
use async_trait::async_trait;
use domain::{
    book::{port::BookWriteRepoPort, Book, BookPrepared},
    book_copy::{
        port::BookCopyWriteRepoPort, BookCopy, BookCopyId, BookCopyPrepared, BookCopyStatus,
    },
    loan::{port::LoanWriteRepoPort, Loan, LoanId, LoanPrepared},
    member::{
        port::MemberWriteRepoPort, Member, MemberId, MemberIdent, MemberPrepared, MemberStatus,
    },
    uow::{UnitOfWorkPort, WriteUnitOfWork, WriteUnitOfWorkFactory},
};
use sqlx::{PgPool, Postgres, Transaction};

use crate::{
    book::write_repo::{create_book_with, get_book_by_isbn_for_write_with},
    book_copy::write_repo::{
        create_book_copy_with, get_book_copy_by_barcode_for_update_with,
        update_book_copy_status_with,
    },
    loan::write_repo::{
        count_active_loans_by_member_id_with, create_loan_with, end_loan_with,
        find_active_loan_by_book_copy_id_for_update_with,
    },
    member::write_repo::{
        create_member_with, get_member_by_ident_for_update_with, update_member_status_with,
    },
};

pub struct SqlUnitOfWork {
    tx: Transaction<'static, Postgres>,
}

#[async_trait]
impl UnitOfWorkPort for SqlUnitOfWork {
    fn book(&mut self) -> &mut dyn BookWriteRepoPort {
        self
    }

    fn book_copy(&mut self) -> &mut dyn BookCopyWriteRepoPort {
        self
    }

    fn member(&mut self) -> &mut dyn MemberWriteRepoPort {
        self
    }

    fn loan(&mut self) -> &mut dyn LoanWriteRepoPort {
        self
    }

    async fn commit(self: Box<Self>) -> anyhow::Result<()> {
        self.tx
            .commit()
            .await
            .context("Failed to commit transaction")
    }
}

#[async_trait]
impl BookWriteRepoPort for SqlUnitOfWork {
    async fn create(&mut self, insert: &BookPrepared) -> anyhow::Result<Book> {
        create_book_with(&mut *self.tx, insert).await
    }

    async fn get_by_isbn(&mut self, isbn: &str) -> anyhow::Result<Option<Book>> {
        get_book_by_isbn_for_write_with(&mut *self.tx, isbn).await
    }
}

#[async_trait]
impl BookCopyWriteRepoPort for SqlUnitOfWork {
    async fn create(&mut self, insert: &BookCopyPrepared) -> anyhow::Result<BookCopy> {
        create_book_copy_with(&mut *self.tx, insert).await
    }

    async fn get_by_barcode_for_update(
        &mut self,
        barcode: &str,
    ) -> anyhow::Result<Option<BookCopy>> {
        get_book_copy_by_barcode_for_update_with(&mut *self.tx, barcode).await
    }

    async fn update_status(
        &mut self,
        id: BookCopyId,
        status: BookCopyStatus,
    ) -> anyhow::Result<()> {
        update_book_copy_status_with(&mut *self.tx, id, status).await
    }
}

#[async_trait]
impl MemberWriteRepoPort for SqlUnitOfWork {
    async fn create(&mut self, insert: &MemberPrepared) -> anyhow::Result<Member> {
        create_member_with(&mut *self.tx, insert).await
    }

    async fn get_by_ident_for_update(
        &mut self,
        ident: &MemberIdent,
    ) -> anyhow::Result<Option<Member>> {
        get_member_by_ident_for_update_with(&mut *self.tx, ident).await
    }

    async fn update_status(&mut self, id: MemberId, status: MemberStatus) -> anyhow::Result<()> {
        update_member_status_with(&mut *self.tx, id, status).await
    }
}

#[async_trait]
impl LoanWriteRepoPort for SqlUnitOfWork {
    async fn create(&mut self, insert: &LoanPrepared) -> anyhow::Result<Loan> {
        create_loan_with(&mut *self.tx, insert).await
    }

    async fn end(&mut self, id: LoanId) -> anyhow::Result<()> {
        end_loan_with(&mut *self.tx, id).await
    }

    async fn find_active_by_book_copy_id_for_update(
        &mut self,
        id: BookCopyId,
    ) -> anyhow::Result<Option<Loan>> {
        find_active_loan_by_book_copy_id_for_update_with(&mut *self.tx, id).await
    }

    async fn count_active_by_member_id(&mut self, id: MemberId) -> anyhow::Result<i64> {
        count_active_loans_by_member_id_with(&mut *self.tx, id).await
    }
}

pub struct SqlWriteUnitOfWorkFactory {
    pub pool: PgPool,
}

#[async_trait]
impl WriteUnitOfWorkFactory for SqlWriteUnitOfWorkFactory {
    async fn build(&self) -> anyhow::Result<WriteUnitOfWork> {
        let tx = self
            .pool
            .begin()
            .await
            .context("Failed to begin transaction")?;
        Ok(WriteUnitOfWork::new(Box::new(SqlUnitOfWork { tx })))
    }
}
