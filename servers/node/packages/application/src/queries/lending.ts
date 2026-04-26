import type { Loan, LoanReadRepository, MemberIdent } from "@library/domain";

export class LendingQueries {
  constructor(private readonly loanReadRepository: LoanReadRepository) {}

  getMemberLoans(ident: MemberIdent): Promise<Loan[]> {
    return this.loanReadRepository.getByMemberIdent(ident);
  }

  getOverdueLoans(): Promise<Loan[]> {
    return this.loanReadRepository.getOverdue();
  }
}
