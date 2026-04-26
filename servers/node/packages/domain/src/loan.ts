import type { BookCopyId } from "./book-copy.js";
import type { MemberId, MemberIdent } from "./member.js";

export type LoanId = number;
export type LoanIdent = string;

export interface Loan {
  id: LoanId;
  ident: LoanIdent;
  dt_created: Date;
  dt_modified: Date;
  book_copy_id: BookCopyId;
  member_id: MemberId;
  dt_due: Date | null;
  dt_returned: Date | null;
}

export interface LoanCreationPayload {
  member_id: MemberId;
  book_copy_id: BookCopyId;
}

export interface LoanPrepared {
  member_id: MemberId;
  book_copy_id: BookCopyId;
}

export function prepareLoan(payload: LoanCreationPayload): LoanPrepared {
  return {
    member_id: payload.member_id,
    book_copy_id: payload.book_copy_id
  };
}

export function ensureLoanCanBeReturned(loan: Loan): void {
  if (loan.dt_returned !== null) {
    throw new LoanCannotBeReturnedError();
  }
}

export class NoActiveLoanForBookCopyError extends Error {
  override readonly name = "NoActiveLoanForBookCopyError";

  constructor() {
    super("Book copy does not have an active loan");
  }
}

export class LoanCannotBeReturnedError extends Error {
  override readonly name = "LoanCannotBeReturnedError";

  constructor() {
    super("Loan has already been returned");
  }
}

export interface LoanWriteRepository {
  create(insert: LoanPrepared): Promise<Loan>;
  end(id: LoanId): Promise<void>;
  findActiveByBookCopyIdForUpdate(id: BookCopyId): Promise<Loan | null>;
  countActiveByMemberId(id: MemberId): Promise<number>;
}

export interface LoanReadRepository {
  getByMemberIdent(ident: MemberIdent): Promise<Loan[]>;
  getOverdue(): Promise<Loan[]>;
  findActiveByBookCopyId(id: BookCopyId): Promise<Loan | null>;
  countActiveByMemberId(id: MemberId): Promise<number>;
}
