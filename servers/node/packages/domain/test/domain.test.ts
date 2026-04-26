import { describe, expect, it } from "vitest";

import {
  BookCopyCannotBeBorrowedError,
  BookCopyCannotBeReturnedFromLostError,
  BookCopyCannotBeReturnedFromMaintenanceError,
  BookCopyCannotBeSentToMaintenanceError,
  BookCopyCannotMarkLostError,
  LoanCannotBeReturnedError,
  MemberCannotBeReactivatedError,
  MemberCannotBeSuspendedError,
  MemberLoanLimitReachedError,
  completeBookCopyMaintenance,
  ensureBookCopyCanBeBorrowed,
  ensureLoanCanBeReturned,
  ensureMemberWithinLoanLimit,
  markBookCopyFound,
  markBookCopyLost,
  reactivateMember,
  sendBookCopyToMaintenance,
  suspendMember,
  type BookCopy,
  type Loan,
  type Member
} from "@library/domain";

const now = new Date("2026-01-01T00:00:00.000Z");

describe("book copy transitions", () => {
  it("allows valid copy state transitions", () => {
    expect(sendBookCopyToMaintenance(copy("active"))).toBe("maintenance");
    expect(completeBookCopyMaintenance(copy("maintenance"))).toBe("active");
    expect(markBookCopyLost(copy("active"))).toBe("lost");
    expect(markBookCopyLost(copy("maintenance"))).toBe("lost");
    expect(markBookCopyFound(copy("lost"))).toBe("active");
  });

  it("rejects invalid copy transitions", () => {
    expect(() => ensureBookCopyCanBeBorrowed(copy("lost"))).toThrow(BookCopyCannotBeBorrowedError);
    expect(() => sendBookCopyToMaintenance(copy("lost"))).toThrow(BookCopyCannotBeSentToMaintenanceError);
    expect(() => completeBookCopyMaintenance(copy("active"))).toThrow(BookCopyCannotBeReturnedFromMaintenanceError);
    expect(() => markBookCopyLost(copy("lost"))).toThrow(BookCopyCannotMarkLostError);
    expect(() => markBookCopyFound(copy("active"))).toThrow(BookCopyCannotBeReturnedFromLostError);
  });
});

describe("member transitions", () => {
  it("allows valid member transitions and loan limit checks", () => {
    expect(suspendMember(member("active"))).toBe("suspended");
    expect(reactivateMember(member("suspended"))).toBe("active");
    expect(() => ensureMemberWithinLoanLimit(member("active"), 2)).not.toThrow();
  });

  it("rejects invalid member transitions and loan limit checks", () => {
    expect(() => suspendMember(member("suspended"))).toThrow(MemberCannotBeSuspendedError);
    expect(() => reactivateMember(member("active"))).toThrow(MemberCannotBeReactivatedError);
    expect(() => ensureMemberWithinLoanLimit(member("active"), 3)).toThrow(MemberLoanLimitReachedError);
  });
});

describe("loan guard", () => {
  it("rejects already returned loans", () => {
    expect(() => ensureLoanCanBeReturned(loan(new Date("2026-01-02T00:00:00.000Z")))).toThrow(
      LoanCannotBeReturnedError
    );
  });
});

function copy(status: BookCopy["status"]): BookCopy {
  return {
    id: 1,
    barcode: "BC-001",
    dt_created: now,
    dt_modified: now,
    book_id: 1,
    status
  };
}

function member(status: Member["status"]): Member {
  return {
    id: 1,
    ident: "MEM-001",
    dt_created: now,
    dt_modified: now,
    status,
    full_name: "Alice Smith",
    max_active_loans: 3
  };
}

function loan(dt_returned: Date | null): Loan {
  return {
    id: 1,
    ident: "LN-000001",
    dt_created: now,
    dt_modified: now,
    book_copy_id: 1,
    member_id: 1,
    dt_due: null,
    dt_returned
  };
}
