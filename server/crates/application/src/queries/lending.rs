use domain::{
    loan::{port::LoanReadRepoPort, Loan},
    member::MemberIdent,
};
use std::sync::Arc;

use super::QueryError;

#[derive(Clone)]
pub struct LendingQueries {
    loan_read_repo: Arc<dyn LoanReadRepoPort>,
}

impl LendingQueries {
    #[must_use]
    pub fn new(loan_read_repo: Arc<dyn LoanReadRepoPort>) -> Self {
        Self { loan_read_repo }
    }

    pub async fn get_member_loans(&self, ident: &MemberIdent) -> Result<Vec<Loan>, QueryError> {
        self.loan_read_repo
            .get_by_member_ident(ident)
            .await
            .map_err(QueryError::from)
    }

    pub async fn get_overdue_loans(&self) -> Result<Vec<Loan>, QueryError> {
        self.loan_read_repo
            .get_overdue()
            .await
            .map_err(QueryError::from)
    }
}
