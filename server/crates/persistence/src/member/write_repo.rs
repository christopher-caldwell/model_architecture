use anyhow::{Context, Result};
use chrono::Utc;
use domain::member::{Member, MemberId, MemberIdent, MemberPrepared, MemberStatus};
use sqlx::PgExecutor;

use crate::member::read_repo::MemberDbRow;

#[derive(sqlx::FromRow)]
pub struct MemberPreparedResult {
    pub member_id: i32,
}

pub(crate) async fn create_member_with(
    executor: impl PgExecutor<'_>,
    insert: &MemberPrepared,
) -> Result<Member> {
    let prepared_result = sqlx::query_file_as!(
        MemberPreparedResult,
        "sql/member/commands/create.sql",
        insert.ident.0.as_str(),
        insert.status.to_string(),
        insert.full_name,
        insert.max_active_loans,
    )
    .fetch_one(executor)
    .await
    .context("Failed to create member")?;

    let now = Utc::now();
    let created_member = Member {
        id: MemberId(prepared_result.member_id),
        ident: insert.ident.clone(),
        dt_created: now,
        dt_modified: now,
        status: insert.status.clone(),
        full_name: insert.full_name.clone(),
        max_active_loans: insert.max_active_loans,
    };
    Ok(created_member)
}

pub(crate) async fn get_member_by_ident_for_update_with(
    executor: impl PgExecutor<'_>,
    ident: &MemberIdent,
) -> Result<Option<Member>> {
    let row = sqlx::query_file_as!(
        MemberDbRow,
        "sql/member/commands/get_by_ident_for_update.sql",
        ident.0.as_str()
    )
    .fetch_optional(executor)
    .await
    .context("Failed to fetch member by ident")?;

    row.map(Member::try_from).transpose()
}

pub(crate) async fn update_member_status_with(
    executor: impl PgExecutor<'_>,
    id: MemberId,
    status: MemberStatus,
) -> Result<()> {
    sqlx::query_file!(
        "sql/member/commands/update_status.sql",
        id.0,
        status.to_string(),
    )
    .execute(executor)
    .await
    .context("Failed to update member status")?;

    Ok(())
}
