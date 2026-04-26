import {
  BookCopyCannotBeBorrowedError,
  BookCopyNotFoundError,
  MemberNotFoundError,
  NoActiveLoanForBookCopyError,
  ensureBookCopyCanBeBorrowed,
  ensureLoanCanBeReturned,
  ensureMemberCanBorrow,
  ensureMemberWithinLoanLimit,
  markBookCopyLost,
  prepareLoan,
  type BookCopy,
  type Loan,
  type Member,
  type WriteUnitOfWork,
  type WriteUnitOfWorkFactory
} from "@library/domain";

import type { CheckOutBookCopyInput } from "./inputs.js";
import { inWriteUnitOfWork } from "./uow.js";

export class LendingCommands {
  constructor(private readonly uowFactory: WriteUnitOfWorkFactory) {}

  async checkOutBookCopy(input: CheckOutBookCopyInput): Promise<Loan> {
    return inWriteUnitOfWork(this.uowFactory, async (uow) => {
      const member = await this.getMemberByIdent(uow, input.member_ident);
      const copy = await this.getBookCopyByBarcode(uow, input.book_copy_barcode);

      ensureMemberCanBorrow(member);
      ensureBookCopyCanBeBorrowed(copy);

      const activeLoanCount = await uow.loans.countActiveByMemberId(member.id);
      ensureMemberWithinLoanLimit(member, activeLoanCount);

      const activeLoan = await uow.loans.findActiveByBookCopyIdForUpdate(copy.id);
      if (activeLoan !== null) {
        throw new BookCopyCannotBeBorrowedError();
      }

      return uow.loans.create(
        prepareLoan({
          member_id: member.id,
          book_copy_id: copy.id
        })
      );
    });
  }

  async returnBookCopy(barcode: string): Promise<Loan> {
    return inWriteUnitOfWork(this.uowFactory, async (uow) => {
      const copy = await this.getBookCopyByBarcode(uow, barcode);
      const loan = await this.getActiveLoanForCopy(uow, copy.id);

      ensureLoanCanBeReturned(loan);
      await uow.loans.end(loan.id);

      const now = new Date();
      return {
        ...loan,
        dt_modified: now,
        dt_returned: now
      };
    });
  }

  async reportLostLoanedBookCopy(barcode: string): Promise<BookCopy> {
    return inWriteUnitOfWork(this.uowFactory, async (uow) => {
      const copy = await this.getBookCopyByBarcode(uow, barcode);
      const status = markBookCopyLost(copy);
      const loan = await this.getActiveLoanForCopy(uow, copy.id);

      ensureLoanCanBeReturned(loan);
      await uow.loans.end(loan.id);
      await uow.bookCopies.updateStatus(copy.id, status);

      return {
        ...copy,
        status,
        dt_modified: new Date()
      };
    });
  }

  private async getMemberByIdent(uow: WriteUnitOfWork, ident: string): Promise<Member> {
    const member = await uow.members.getByIdentForUpdate(ident);
    if (member === null) throw new MemberNotFoundError();
    return member;
  }

  private async getBookCopyByBarcode(uow: WriteUnitOfWork, barcode: string): Promise<BookCopy> {
    const copy = await uow.bookCopies.getByBarcodeForUpdate(barcode);
    if (copy === null) throw new BookCopyNotFoundError();
    return copy;
  }

  private async getActiveLoanForCopy(
    uow: WriteUnitOfWork,
    book_copy_id: BookCopy["id"]
  ): Promise<Loan> {
    const loan = await uow.loans.findActiveByBookCopyIdForUpdate(book_copy_id);
    if (loan === null) throw new NoActiveLoanForBookCopyError();
    return loan;
  }
}
