use anyhow::{Context, Result};
use chrono::Utc;
use domain::{
    book_copy::BookCopyId,
    loan::{Loan, LoanId, LoanIdent, LoanPrepared},
    member::MemberId,
};
use sqlx::PgExecutor;

use crate::loan::read_repo::{CountDbRow, LoanDbRow};

#[derive(sqlx::FromRow)]
pub struct LoanPreparedResult {
    pub loan_id: i32,
    pub loan_ident: String,
}

pub(crate) async fn create_loan_with(
    executor: impl PgExecutor<'_>,
    insert: &LoanPrepared,
) -> Result<Loan> {
    let prepared_result = sqlx::query_file_as!(
        LoanPreparedResult,
        "sql/loan/commands/create.sql",
        insert.book_copy_id.0,
        insert.member_id.0,
    )
    .fetch_one(executor)
    .await
    .context("Failed to create loan")?;

    let now = Utc::now();
    let created_loan = Loan {
        id: LoanId(prepared_result.loan_id),
        ident: LoanIdent(prepared_result.loan_ident),
        dt_created: now,
        dt_modified: now,
        book_copy_id: insert.book_copy_id,
        member_id: insert.member_id,
        dt_due: None,
        dt_returned: None,
    };
    Ok(created_loan)
}

pub(crate) async fn end_loan_with(executor: impl PgExecutor<'_>, id: LoanId) -> Result<()> {
    sqlx::query_file!("sql/loan/commands/end.sql", id.0,)
        .execute(executor)
        .await
        .context("Failed to end loan")?;

    Ok(())
}

pub(crate) async fn find_active_loan_by_book_copy_id_for_update_with(
    executor: impl PgExecutor<'_>,
    id: BookCopyId,
) -> Result<Option<Loan>> {
    let row = sqlx::query_file_as!(
        LoanDbRow,
        "sql/loan/commands/find_active_by_book_copy_id_for_update.sql",
        id.0
    )
    .fetch_optional(executor)
    .await
    .context("Failed to find active loan by book copy id")?;

    row.map(Loan::try_from).transpose()
}

pub(crate) async fn count_active_loans_by_member_id_with(
    executor: impl PgExecutor<'_>,
    id: MemberId,
) -> Result<i64> {
    let row = sqlx::query_file_as!(
        CountDbRow,
        "sql/loan/commands/count_active_by_member_id.sql",
        id.0
    )
    .fetch_one(executor)
    .await
    .context("Failed to count active loans by member id")?;

    Ok(row.count)
}
