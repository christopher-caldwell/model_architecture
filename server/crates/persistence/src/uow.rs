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
    PortError, PortResult,
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

    async fn commit(self: Box<Self>) -> PortResult<()> {
        self.tx
            .commit()
            .await
            .context("Failed to commit transaction")
            .map_err(|error| PortError::unit_of_work(error.into_boxed_dyn_error()))
    }
}

#[async_trait]
impl BookWriteRepoPort for SqlUnitOfWork {
    async fn create(&mut self, insert: &BookPrepared) -> PortResult<Book> {
        create_book_with(&mut *self.tx, insert)
            .await
            .map_err(|error| PortError::repository(error.into_boxed_dyn_error()))
    }

    async fn get_by_isbn(&mut self, isbn: &str) -> PortResult<Option<Book>> {
        get_book_by_isbn_for_write_with(&mut *self.tx, isbn)
            .await
            .map_err(|error| PortError::repository(error.into_boxed_dyn_error()))
    }
}

#[async_trait]
impl BookCopyWriteRepoPort for SqlUnitOfWork {
    async fn create(&mut self, insert: &BookCopyPrepared) -> PortResult<BookCopy> {
        create_book_copy_with(&mut *self.tx, insert)
            .await
            .map_err(|error| PortError::repository(error.into_boxed_dyn_error()))
    }

    async fn get_by_barcode_for_update(&mut self, barcode: &str) -> PortResult<Option<BookCopy>> {
        get_book_copy_by_barcode_for_update_with(&mut *self.tx, barcode)
            .await
            .map_err(|error| PortError::repository(error.into_boxed_dyn_error()))
    }

    async fn update_status(&mut self, id: BookCopyId, status: BookCopyStatus) -> PortResult<()> {
        update_book_copy_status_with(&mut *self.tx, id, status)
            .await
            .map_err(|error| PortError::repository(error.into_boxed_dyn_error()))
    }
}

#[async_trait]
impl MemberWriteRepoPort for SqlUnitOfWork {
    async fn create(&mut self, insert: &MemberPrepared) -> PortResult<Member> {
        create_member_with(&mut *self.tx, insert)
            .await
            .map_err(|error| PortError::repository(error.into_boxed_dyn_error()))
    }

    async fn get_by_ident_for_update(&mut self, ident: &MemberIdent) -> PortResult<Option<Member>> {
        get_member_by_ident_for_update_with(&mut *self.tx, ident)
            .await
            .map_err(|error| PortError::repository(error.into_boxed_dyn_error()))
    }

    async fn update_status(&mut self, id: MemberId, status: MemberStatus) -> PortResult<()> {
        update_member_status_with(&mut *self.tx, id, status)
            .await
            .map_err(|error| PortError::repository(error.into_boxed_dyn_error()))
    }
}

#[async_trait]
impl LoanWriteRepoPort for SqlUnitOfWork {
    async fn create(&mut self, insert: &LoanPrepared) -> PortResult<Loan> {
        create_loan_with(&mut *self.tx, insert)
            .await
            .map_err(|error| PortError::repository(error.into_boxed_dyn_error()))
    }

    async fn end(&mut self, id: LoanId) -> PortResult<()> {
        end_loan_with(&mut *self.tx, id)
            .await
            .map_err(|error| PortError::repository(error.into_boxed_dyn_error()))
    }

    async fn find_active_by_book_copy_id_for_update(
        &mut self,
        id: BookCopyId,
    ) -> PortResult<Option<Loan>> {
        find_active_loan_by_book_copy_id_for_update_with(&mut *self.tx, id)
            .await
            .map_err(|error| PortError::repository(error.into_boxed_dyn_error()))
    }

    async fn count_active_by_member_id(&mut self, id: MemberId) -> PortResult<i64> {
        count_active_loans_by_member_id_with(&mut *self.tx, id)
            .await
            .map_err(|error| PortError::repository(error.into_boxed_dyn_error()))
    }
}

pub struct SqlWriteUnitOfWorkFactory {
    pub pool: PgPool,
}

#[async_trait]
impl WriteUnitOfWorkFactory for SqlWriteUnitOfWorkFactory {
    async fn build(&self) -> PortResult<WriteUnitOfWork> {
        let tx = self
            .pool
            .begin()
            .await
            .context("Failed to begin transaction")
            .map_err(|error| PortError::unit_of_work(error.into_boxed_dyn_error()))?;
        Ok(WriteUnitOfWork::new(Box::new(SqlUnitOfWork { tx })))
    }
}
