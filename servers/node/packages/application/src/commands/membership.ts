import {
  MemberNotFoundError,
  prepareMember,
  reactivateMember,
  suspendMember,
  type Member,
  type MemberCreationPayload,
  type WriteUnitOfWork,
  type WriteUnitOfWorkFactory
} from "@library/domain";

import type { IdentGenerator } from "../ports.js";
import type { MemberIdentInput } from "./inputs.js";
import { inWriteUnitOfWork } from "./uow.js";

export class MembershipCommands {
  constructor(
    private readonly uowFactory: WriteUnitOfWorkFactory,
    private readonly identGenerator: IdentGenerator
  ) {}

  async registerMember(payload: MemberCreationPayload): Promise<Member> {
    const prepared = prepareMember(payload, this.identGenerator.generate());

    return inWriteUnitOfWork(this.uowFactory, async (uow) => {
      return uow.members.create(prepared);
    });
  }

  async suspendMember(input: MemberIdentInput): Promise<Member> {
    return this.updateMemberStatus(input.member_ident, suspendMember);
  }

  async reactivateMember(input: MemberIdentInput): Promise<Member> {
    return this.updateMemberStatus(input.member_ident, reactivateMember);
  }

  private async updateMemberStatus(
    member_ident: string,
    transition: (member: Member) => Member["status"]
  ): Promise<Member> {
    return inWriteUnitOfWork(this.uowFactory, async (uow) => {
      const member = await this.getMemberByIdent(uow, member_ident);
      const status = transition(member);

      await uow.members.updateStatus(member.id, status);

      return {
        ...member,
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
}
