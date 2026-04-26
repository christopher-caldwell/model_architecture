export type MemberId = number;
export type MemberIdent = string;

export const memberStatuses = ["active", "suspended"] as const;
export type MemberStatus = (typeof memberStatuses)[number];

export interface Member {
  id: MemberId;
  ident: MemberIdent;
  dt_created: Date;
  dt_modified: Date;
  status: MemberStatus;
  full_name: string;
  max_active_loans: number;
}

export interface MemberCreationPayload {
  full_name: string;
  max_active_loans: number;
}

export interface MemberPrepared {
  ident: MemberIdent;
  full_name: string;
  max_active_loans: number;
  status: MemberStatus;
}

export function prepareMember(payload: MemberCreationPayload, ident: MemberIdent): MemberPrepared {
  return {
    ident,
    full_name: payload.full_name,
    max_active_loans: payload.max_active_loans,
    status: "active"
  };
}

export function suspendMember(member: Member): MemberStatus {
  if (member.status === "suspended") {
    throw new MemberCannotBeSuspendedError();
  }

  return "suspended";
}

export function reactivateMember(member: Member): MemberStatus {
  if (member.status !== "suspended") {
    throw new MemberCannotBeReactivatedError();
  }

  return "active";
}

export function ensureMemberCanBorrow(member: Member): void {
  if (member.status !== "active") {
    throw new MemberCannotBorrowWhileSuspendedError();
  }
}

export function ensureMemberWithinLoanLimit(member: Member, activeLoanCount: number): void {
  if (activeLoanCount >= member.max_active_loans) {
    throw new MemberLoanLimitReachedError();
  }
}

export class MemberNotFoundError extends Error {
  override readonly name = "MemberNotFoundError";

  constructor() {
    super("Member not found");
  }
}

export class MemberCannotBeSuspendedError extends Error {
  override readonly name = "MemberCannotBeSuspendedError";

  constructor() {
    super("Member is already suspended");
  }
}

export class MemberCannotBeReactivatedError extends Error {
  override readonly name = "MemberCannotBeReactivatedError";

  constructor() {
    super("Member is not currently suspended");
  }
}

export class MemberCannotBorrowWhileSuspendedError extends Error {
  override readonly name = "MemberCannotBorrowWhileSuspendedError";

  constructor() {
    super("Member is suspended and cannot borrow new books");
  }
}

export class MemberLoanLimitReachedError extends Error {
  override readonly name = "MemberLoanLimitReachedError";

  constructor() {
    super("Member has reached the maximum number of active loans");
  }
}

export interface MemberWriteRepository {
  create(insert: MemberPrepared): Promise<Member>;
  getByIdentForUpdate(ident: MemberIdent): Promise<Member | null>;
  updateStatus(id: MemberId, status: MemberStatus): Promise<void>;
}

export interface MemberReadRepository {
  getById(id: MemberId): Promise<Member | null>;
  getByIdent(ident: MemberIdent): Promise<Member | null>;
}
