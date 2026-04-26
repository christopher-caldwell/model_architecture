import type { Member, MemberIdent, MemberReadRepository } from "@library/domain";

export class MembershipQueries {
  constructor(private readonly memberReadRepository: MemberReadRepository) {}

  getMemberDetails(ident: MemberIdent): Promise<Member | null> {
    return this.memberReadRepository.getByIdent(ident);
  }
}
