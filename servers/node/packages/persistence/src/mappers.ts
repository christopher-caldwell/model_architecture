import {
  bookCopyStatuses,
  memberStatuses,
  type Book,
  type BookCopy,
  type BookCopyStatus,
  type Loan,
  type Member,
  type MemberStatus
} from "@library/domain";

export interface BookRow {
  book_id: number;
  isbn: string;
  dt_created: Date;
  dt_modified: Date;
  title: string;
  author_name: string;
}

export interface BookCopyRow {
  book_copy_id: number;
  barcode: string;
  dt_created: Date;
  dt_modified: Date;
  book_id: number;
  status: string;
}

export interface MemberRow {
  member_id: number;
  member_ident: string;
  dt_created: Date;
  dt_modified: Date;
  status: string;
  full_name: string;
  max_active_loans: number;
}

export interface LoanRow {
  loan_id: number;
  loan_ident: string;
  dt_created: Date;
  dt_modified: Date;
  book_copy_id: number;
  member_id: number;
  dt_due: Date | null;
  dt_returned: Date | null;
}

export function mapBook(row: BookRow): Book {
  return {
    id: row.book_id,
    isbn: row.isbn,
    dt_created: row.dt_created,
    dt_modified: row.dt_modified,
    title: row.title,
    author_name: row.author_name
  };
}

export function mapBookCopy(row: BookCopyRow): BookCopy {
  return {
    id: row.book_copy_id,
    barcode: row.barcode,
    dt_created: row.dt_created,
    dt_modified: row.dt_modified,
    book_id: row.book_id,
    status: parseBookCopyStatus(row.status)
  };
}

export function mapMember(row: MemberRow): Member {
  return {
    id: row.member_id,
    ident: row.member_ident,
    dt_created: row.dt_created,
    dt_modified: row.dt_modified,
    status: parseMemberStatus(row.status),
    full_name: row.full_name,
    max_active_loans: row.max_active_loans
  };
}

export function mapLoan(row: LoanRow): Loan {
  return {
    id: row.loan_id,
    ident: row.loan_ident,
    dt_created: row.dt_created,
    dt_modified: row.dt_modified,
    book_copy_id: row.book_copy_id,
    member_id: row.member_id,
    dt_due: row.dt_due,
    dt_returned: row.dt_returned
  };
}

function parseBookCopyStatus(value: string): BookCopyStatus {
  if (bookCopyStatuses.includes(value as BookCopyStatus)) return value as BookCopyStatus;
  throw new Error(`Invalid book copy status in DB: ${value}`);
}

function parseMemberStatus(value: string): MemberStatus {
  if (memberStatuses.includes(value as MemberStatus)) return value as MemberStatus;
  throw new Error(`Invalid member status in DB: ${value}`);
}
