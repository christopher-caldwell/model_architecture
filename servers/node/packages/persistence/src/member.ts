import type { Pool, PoolClient } from "pg";

import type {
  Member,
  MemberId,
  MemberIdent,
  MemberPrepared,
  MemberReadRepository,
  MemberStatus,
  MemberWriteRepository
} from "@library/domain";

import type { MemberRow } from "./mappers.js";
import { mapMember } from "./mappers.js";
import { sql } from "./sql.js";

interface MemberCreateRow {
  member_id: number;
}

export class MemberReadRepositoryPostgres implements MemberReadRepository {
  constructor(private readonly pool: Pool) {}

  async getById(id: MemberId): Promise<Member | null> {
    const result = await this.pool.query<MemberRow>(sql.member.getById, [id]);
    return result.rows[0] === undefined ? null : mapMember(result.rows[0]);
  }

  async getByIdent(ident: MemberIdent): Promise<Member | null> {
    const result = await this.pool.query<MemberRow>(sql.member.getByIdent, [ident]);
    return result.rows[0] === undefined ? null : mapMember(result.rows[0]);
  }
}

export class MemberWriteRepositoryPostgres implements MemberWriteRepository {
  constructor(private readonly client: PoolClient) {}

  async create(insert: MemberPrepared): Promise<Member> {
    const result = await this.client.query<MemberCreateRow>(sql.member.create, [
      insert.ident,
      insert.status,
      insert.full_name,
      insert.max_active_loans
    ]);
    const created = result.rows[0];
    if (created === undefined) throw new Error("Failed to create member");

    const now = new Date();
    return {
      id: created.member_id,
      ident: insert.ident,
      dt_created: now,
      dt_modified: now,
      status: insert.status,
      full_name: insert.full_name,
      max_active_loans: insert.max_active_loans
    };
  }

  async getByIdentForUpdate(ident: MemberIdent): Promise<Member | null> {
    const result = await this.client.query<MemberRow>(sql.member.getByIdentForUpdate, [ident]);
    return result.rows[0] === undefined ? null : mapMember(result.rows[0]);
  }

  async updateStatus(id: MemberId, status: MemberStatus): Promise<void> {
    await this.client.query(sql.member.updateStatus, [id, status]);
  }
}
