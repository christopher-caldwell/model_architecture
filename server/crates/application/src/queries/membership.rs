use domain::member::{port::MemberReadRepoPort, Member, MemberIdent};
use std::sync::Arc;

use super::QueryError;

#[derive(Clone)]
pub struct MembershipQueries {
    member_read_repo: Arc<dyn MemberReadRepoPort>,
}

impl MembershipQueries {
    #[must_use]
    pub fn new(member_read_repo: Arc<dyn MemberReadRepoPort>) -> Self {
        Self { member_read_repo }
    }

    pub async fn get_member_details(
        &self,
        ident: &MemberIdent,
    ) -> Result<Option<Member>, QueryError> {
        self.member_read_repo
            .get_by_ident(ident)
            .await
            .map_err(QueryError::from)
    }
}
