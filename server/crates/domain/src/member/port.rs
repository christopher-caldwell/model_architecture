use async_trait::async_trait;

use super::{Member, MemberId, MemberIdent, MemberPrepared, MemberStatus};

#[async_trait]
pub trait MemberWriteRepoPort: Send {
    async fn create(&mut self, insert: &MemberPrepared) -> anyhow::Result<Member>;
    async fn get_by_ident_for_update(
        &mut self,
        ident: &MemberIdent,
    ) -> anyhow::Result<Option<Member>>;
    async fn update_status(&mut self, id: MemberId, status: MemberStatus) -> anyhow::Result<()>;
}

#[async_trait]
pub trait MemberReadRepoPort: Send + Sync {
    async fn get_by_id(&self, id: MemberId) -> anyhow::Result<Option<Member>>;
    async fn get_by_ident(&self, ident: &MemberIdent) -> anyhow::Result<Option<Member>>;
}
