export interface AddBookCopyInput {
  isbn: string;
  barcode: string;
}

export interface CheckOutBookCopyInput {
  member_ident: string;
  book_copy_barcode: string;
}

export interface MemberIdentInput {
  member_ident: string;
}
