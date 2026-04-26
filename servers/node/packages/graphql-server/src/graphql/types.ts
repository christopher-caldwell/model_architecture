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
  status: "active" | "maintenance" | "lost";
}

interface MemberView {
  ident: string;
  dt_created: Date;
  dt_modified: Date;
  status: "active" | "suspended";
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

export interface CatalogTitle {
  isbn: string;
  dt_created: string;
  dt_modified: string;
  title: string;
  author_name: string;
}

export interface InventoryCopy {
  barcode: string;
  dt_created: string;
  dt_modified: string;
  status: InventoryCopyStatus;
}

export interface LibraryMember {
  member_number: string;
  dt_created: string;
  dt_modified: string;
  status: LibraryMemberStatus;
  full_name: string;
  max_active_loans: number;
}

export type InventoryCopyStatus = "ACTIVE" | "MAINTENANCE" | "LOST";
export type LibraryMemberStatus = "ACTIVE" | "SUSPENDED";

export interface LoanRecord {
  loan_number: string;
  dt_created: string;
  dt_modified: string;
  dt_due: string | null;
  dt_returned: string | null;
}

export function catalogTitle(book: BookView): CatalogTitle {
  return {
    isbn: book.isbn,
    dt_created: book.dt_created.toISOString(),
    dt_modified: book.dt_modified.toISOString(),
    title: book.title,
    author_name: book.author_name
  };
}

export function inventoryCopy(copy: BookCopyView): InventoryCopy {
  return {
    barcode: copy.barcode,
    dt_created: copy.dt_created.toISOString(),
    dt_modified: copy.dt_modified.toISOString(),
    status: toGraphqlEnum(copy.status)
  };
}

export function libraryMember(member: MemberView): LibraryMember {
  return {
    member_number: member.ident,
    dt_created: member.dt_created.toISOString(),
    dt_modified: member.dt_modified.toISOString(),
    status: toGraphqlEnum(member.status),
    full_name: member.full_name,
    max_active_loans: member.max_active_loans
  };
}

export function loanRecord(loan: LoanView): LoanRecord {
  return {
    loan_number: loan.ident,
    dt_created: loan.dt_created.toISOString(),
    dt_modified: loan.dt_modified.toISOString(),
    dt_due: loan.dt_due?.toISOString() ?? null,
    dt_returned: loan.dt_returned?.toISOString() ?? null
  };
}

function toGraphqlEnum(status: "active" | "maintenance" | "lost"): InventoryCopyStatus;
function toGraphqlEnum(status: "active" | "suspended"): LibraryMemberStatus;
function toGraphqlEnum(status: string): InventoryCopyStatus | LibraryMemberStatus {
  return status.toUpperCase() as InventoryCopyStatus | LibraryMemberStatus;
}
