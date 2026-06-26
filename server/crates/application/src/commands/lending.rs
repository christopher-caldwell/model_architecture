use chrono::Utc;
use domain::{
    book_copy::{BookCopy, BookCopyError},
    loan::{Loan, LoanCreationPayload, LoanError},
    member::{Member, MemberError},
    uow::{WriteUnitOfWork, WriteUnitOfWorkFactory},
};
use std::sync::Arc;

#[derive(Clone)]
pub struct LendingCommands {
    uow_factory: Arc<dyn WriteUnitOfWorkFactory>,
}

impl LendingCommands {
    #[must_use]
    pub fn new(uow_factory: Arc<dyn WriteUnitOfWorkFactory>) -> Self {
        Self { uow_factory }
    }

    async fn get_member_by_ident(
        &self,
        uow: &mut WriteUnitOfWork,
        member_ident: &str,
    ) -> Result<Member, super::CommandError> {
        let ident = domain::member::MemberIdent(member_ident.to_owned());
        uow.member()
            .get_by_ident_for_update(&ident)
            .await?
            .ok_or(MemberError::NotFound.into())
    }

    async fn get_book_copy_by_barcode(
        &self,
        uow: &mut WriteUnitOfWork,
        barcode: &str,
    ) -> Result<BookCopy, super::CommandError> {
        uow.book_copy()
            .get_by_barcode_for_update(barcode)
            .await?
            .ok_or(BookCopyError::NotFound.into())
    }

    async fn load_active_loan_for_copy(
        &self,
        uow: &mut WriteUnitOfWork,
        book_copy_id: domain::book_copy::BookCopyId,
    ) -> Result<Option<Loan>, super::CommandError> {
        Ok(uow
            .loan()
            .find_active_by_book_copy_id_for_update(book_copy_id)
            .await?)
    }

    pub async fn check_out_book_copy(
        &self,
        input: super::CheckOutBookCopyInput,
    ) -> Result<Loan, super::CommandError> {
        let mut uow = self.uow_factory.build().await?;
        let member = self
            .get_member_by_ident(&mut uow, &input.member_ident)
            .await?;
        let book_copy = self
            .get_book_copy_by_barcode(&mut uow, &input.book_copy_barcode)
            .await?;

        member.ensure_can_borrow()?;
        book_copy.ensure_can_be_borrowed()?;

        let active_loan_count = uow.loan().count_active_by_member_id(member.id).await?;
        member.ensure_within_loan_limit(active_loan_count)?;

        let active_loan = self
            .load_active_loan_for_copy(&mut uow, book_copy.id)
            .await?;
        if active_loan.is_some() {
            return Err(BookCopyError::CannotBeBorrowed.into());
        }

        let prepared = LoanCreationPayload {
            member_id: member.id,
            book_copy_id: book_copy.id,
        }
        .prepare();
        let result = uow.loan().create(&prepared).await?;
        uow.commit().await?;
        Ok(result)
    }

    pub async fn return_book_copy(&self, barcode: String) -> Result<Loan, super::CommandError> {
        let mut uow = self.uow_factory.build().await?;
        let book_copy = self.get_book_copy_by_barcode(&mut uow, &barcode).await?;
        let loan = self
            .load_active_loan_for_copy(&mut uow, book_copy.id)
            .await?
            .ok_or(LoanError::NoActiveLoanForBookCopy)?;
        loan.ensure_can_be_returned()?;
        uow.loan().end(loan.id).await?;
        uow.commit().await?;
        let now = Utc::now();
        let updated_loan = Loan {
            dt_modified: now,
            dt_returned: Some(now),
            ..loan
        };
        Ok(updated_loan)
    }

    pub async fn report_lost_loaned_book_copy(
        &self,
        barcode: String,
    ) -> Result<BookCopy, super::CommandError> {
        let mut uow = self.uow_factory.build().await?;
        let book_copy = self.get_book_copy_by_barcode(&mut uow, &barcode).await?;
        let lost_status = book_copy.mark_lost()?;
        let loan = self
            .load_active_loan_for_copy(&mut uow, book_copy.id)
            .await?
            .ok_or(LoanError::NoActiveLoanForBookCopy)?;
        loan.ensure_can_be_returned()?;
        uow.loan().end(loan.id).await?;
        uow.book_copy()
            .update_status(book_copy.id, lost_status.clone())
            .await?;
        uow.commit().await?;
        let updated_book_copy = BookCopy {
            status: lost_status,
            dt_modified: Utc::now(),
            ..book_copy
        };
        Ok(updated_book_copy)
    }
}
