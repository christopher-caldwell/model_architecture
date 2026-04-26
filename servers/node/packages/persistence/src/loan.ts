import type { Pool, PoolClient } from "pg";

import type {
  BookCopyId,
  Loan,
  LoanId,
  LoanPrepared,
  LoanReadRepository,
  LoanWriteRepository,
  MemberId,
  MemberIdent
} from "@library/domain";

import type { LoanRow } from "./mappers.js";
import { mapLoan } from "./mappers.js";
import { sql } from "./sql.js";

interface CountRow {
  count: string;
}

interface LoanCreateRow {
  loan_id: number;
  loan_ident: string;
}

export class LoanReadRepositoryPostgres implements LoanReadRepository {
  constructor(private readonly pool: Pool) {}

  async getByMemberIdent(ident: MemberIdent): Promise<Loan[]> {
    const result = await this.pool.query<LoanRow>(sql.loan.getByMemberIdent, [ident]);
    return result.rows.map(mapLoan);
  }

  async getOverdue(): Promise<Loan[]> {
    const result = await this.pool.query<LoanRow>(sql.loan.getOverdue);
    return result.rows.map(mapLoan);
  }

  async findActiveByBookCopyId(id: BookCopyId): Promise<Loan | null> {
    const result = await this.pool.query<LoanRow>(sql.loan.findActiveByBookCopyId, [id]);
    return result.rows[0] === undefined ? null : mapLoan(result.rows[0]);
  }

  async countActiveByMemberId(id: MemberId): Promise<number> {
    const result = await this.pool.query<CountRow>(sql.loan.countActiveByMemberId, [id]);
    const row = result.rows[0];
    return row === undefined ? 0 : Number(row.count);
  }
}

export class LoanWriteRepositoryPostgres implements LoanWriteRepository {
  constructor(private readonly client: PoolClient) {}

  async create(insert: LoanPrepared): Promise<Loan> {
    const result = await this.client.query<LoanCreateRow>(sql.loan.create, [
      insert.book_copy_id,
      insert.member_id
    ]);
    const created = result.rows[0];
    if (created === undefined) throw new Error("Failed to create loan");

    const now = new Date();
    return {
      id: created.loan_id,
      ident: created.loan_ident,
      dt_created: now,
      dt_modified: now,
      book_copy_id: insert.book_copy_id,
      member_id: insert.member_id,
      dt_due: null,
      dt_returned: null
    };
  }

  async end(id: LoanId): Promise<void> {
    await this.client.query(sql.loan.end, [id]);
  }

  async findActiveByBookCopyIdForUpdate(id: BookCopyId): Promise<Loan | null> {
    const result = await this.client.query<LoanRow>(sql.loan.findActiveByBookCopyIdForUpdate, [id]);
    return result.rows[0] === undefined ? null : mapLoan(result.rows[0]);
  }

  async countActiveByMemberId(id: MemberId): Promise<number> {
    const result = await this.client.query<CountRow>(sql.loan.countActiveByMemberId, [id]);
    const row = result.rows[0];
    return row === undefined ? 0 : Number(row.count);
  }
}
