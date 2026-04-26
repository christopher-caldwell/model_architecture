interface BookView {
  isbn: string;
  dt_created: Date;
  dt_modified: Date;
  title: string;
  author_name: string;
}

interface BookCopyView {
  barcode: string;
  dt_created: Date;
  dt_modified: Date;
  status: string;
}

interface MemberView {
  ident: string;
  dt_created: Date;
  dt_modified: Date;
  status: string;
  full_name: string;
  max_active_loans: number;
}

interface LoanView {
  ident: string;
  dt_created: Date;
  dt_modified: Date;
  dt_due: Date | null;
  dt_returned: Date | null;
}

export interface BookResponseBody {
  isbn: string;
  dt_created: string;
  dt_modified: string;
  title: string;
  author_name: string;
}

export interface CreateBookRequestBody {
  isbn: string;
  title: string;
  author_name: string;
}

export interface BookCopyResponseBody {
  barcode: string;
  dt_created: string;
  dt_modified: string;
  status: string;
}

export interface CreateBookCopyRequestBody {
  barcode: string;
}

export interface MemberResponseBody {
  ident: string;
  dt_created: string;
  dt_modified: string;
  status: string;
  full_name: string;
  max_active_loans: number;
}

export interface CreateMemberRequestBody {
  full_name: string;
  max_active_loans: number;
}

export interface LoanResponseBody {
  ident: string;
  dt_created: string;
  dt_modified: string;
  dt_due: string | null;
  dt_returned: string | null;
}

export interface CreateLoanRequestBody {
  member_ident: string;
  book_copy_barcode: string;
}

export interface ErrorResponseBody {
  error: string;
}

export interface HealthCheckResponseBody {
  message: string;
}

export function bookResponse(book: BookView): BookResponseBody {
  return {
    isbn: book.isbn,
    dt_created: book.dt_created.toISOString(),
    dt_modified: book.dt_modified.toISOString(),
    title: book.title,
    author_name: book.author_name
  };
}

export function bookCopyResponse(copy: BookCopyView): BookCopyResponseBody {
  return {
    barcode: copy.barcode,
    dt_created: copy.dt_created.toISOString(),
    dt_modified: copy.dt_modified.toISOString(),
    status: copy.status
  };
}

export function memberResponse(member: MemberView): MemberResponseBody {
  return {
    ident: member.ident,
    dt_created: member.dt_created.toISOString(),
    dt_modified: member.dt_modified.toISOString(),
    status: member.status,
    full_name: member.full_name,
    max_active_loans: member.max_active_loans
  };
}

export function loanResponse(loan: LoanView): LoanResponseBody {
  return {
    ident: loan.ident,
    dt_created: loan.dt_created.toISOString(),
    dt_modified: loan.dt_modified.toISOString(),
    dt_due: loan.dt_due?.toISOString() ?? null,
    dt_returned: loan.dt_returned?.toISOString() ?? null
  };
}
