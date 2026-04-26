import { describe, expect, it } from "vitest";

import { LendingCommands } from "@library/application";
import {
  BookCopyCannotBeBorrowedError,
  type Book,
  type BookCopy,
  type BookCopyStatus,
  type BookPrepared,
  type BookWriteRepository,
  type BookCopyPrepared,
  type BookCopyWriteRepository,
  type Loan,
  type LoanId,
  type LoanPrepared,
  type LoanWriteRepository,
  type Member,
  type MemberId,
  type MemberIdent,
  type MemberPrepared,
  type MemberStatus,
  type MemberWriteRepository,
  type WriteUnitOfWork,
  type WriteUnitOfWorkFactory
} from "@library/domain";

const now = new Date("2026-01-01T00:00:00.000Z");

describe("LendingCommands", () => {
  it("checks out a borrowable copy for an active member inside a unit of work", async () => {
    const uow = new FakeUnitOfWork();
    const commands = new LendingCommands(new FakeUnitOfWorkFactory(uow));

    const loan = await commands.checkOutBookCopy({
      member_ident: "MEM-001",
      book_copy_barcode: "BC-001"
    });

    expect(loan.ident).toBe("LN-000001");
    expect(uow.committed).toBe(true);
    expect(uow.rolledBack).toBe(false);
    expect(uow.loans.created).toHaveLength(1);
  });

  it("rolls back when a copy already has an active loan", async () => {
    const uow = new FakeUnitOfWork();
    uow.loans.activeLoan = sampleLoan();
    const commands = new LendingCommands(new FakeUnitOfWorkFactory(uow));

    await expect(
      commands.checkOutBookCopy({
        member_ident: "MEM-001",
        book_copy_barcode: "BC-001"
      })
    ).rejects.toThrow(BookCopyCannotBeBorrowedError);

    expect(uow.committed).toBe(false);
    expect(uow.rolledBack).toBe(true);
    expect(uow.loans.created).toHaveLength(0);
  });
});

class FakeUnitOfWorkFactory implements WriteUnitOfWorkFactory {
  constructor(private readonly uow: WriteUnitOfWork) {}

  async build(): Promise<WriteUnitOfWork> {
    return this.uow;
  }
}

class FakeUnitOfWork implements WriteUnitOfWork {
  readonly books = new FakeBookWriteRepository();
  readonly bookCopies = new FakeBookCopyWriteRepository();
  readonly members = new FakeMemberWriteRepository();
  readonly loans = new FakeLoanWriteRepository();
  committed = false;
  rolledBack = false;

  async commit(): Promise<void> {
    this.committed = true;
  }

  async rollback(): Promise<void> {
    this.rolledBack = true;
  }
}

class FakeBookWriteRepository implements BookWriteRepository {
  async create(_insert: BookPrepared): Promise<Book> {
    throw new Error("not used");
  }

  async getByIsbn(_isbn: string): Promise<Book | null> {
    throw new Error("not used");
  }
}

class FakeBookCopyWriteRepository implements BookCopyWriteRepository {
  copy: BookCopy = {
    id: 1,
    barcode: "BC-001",
    dt_created: now,
    dt_modified: now,
    book_id: 1,
    status: "active"
  };

  async create(_insert: BookCopyPrepared): Promise<BookCopy> {
    throw new Error("not used");
  }

  async getByBarcodeForUpdate(_barcode: string): Promise<BookCopy | null> {
    return this.copy;
  }

  async updateStatus(_id: number, status: BookCopyStatus): Promise<void> {
    this.copy = {
      ...this.copy,
      status
    };
  }
}

class FakeMemberWriteRepository implements MemberWriteRepository {
  member: Member = {
    id: 1,
    ident: "MEM-001",
    dt_created: now,
    dt_modified: now,
    status: "active",
    full_name: "Alice Smith",
    max_active_loans: 3
  };

  async create(_insert: MemberPrepared): Promise<Member> {
    throw new Error("not used");
  }

  async getByIdentForUpdate(_ident: MemberIdent): Promise<Member | null> {
    return this.member;
  }

  async updateStatus(_id: MemberId, status: MemberStatus): Promise<void> {
    this.member = {
      ...this.member,
      status
    };
  }
}

class FakeLoanWriteRepository implements LoanWriteRepository {
  created: LoanPrepared[] = [];
  activeLoan: Loan | null = null;

  async create(insert: LoanPrepared): Promise<Loan> {
    this.created.push(insert);
    return sampleLoan();
  }

  async end(_id: LoanId): Promise<void> {}

  async findActiveByBookCopyIdForUpdate(_id: number): Promise<Loan | null> {
    return this.activeLoan;
  }

  async countActiveByMemberId(_id: MemberId): Promise<number> {
    return 0;
  }
}

function sampleLoan(): Loan {
  return {
    id: 1,
    ident: "LN-000001",
    dt_created: now,
    dt_modified: now,
    book_copy_id: 1,
    member_id: 1,
    dt_due: null,
    dt_returned: null
  };
}
