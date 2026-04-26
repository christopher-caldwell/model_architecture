import { builder, safeResolver } from "../schema-builder.js";
import { libraryMember } from "../types.js";
import { LibraryMemberRef } from "./types.js";

builder.queryFields((t) => ({
  member: t.field({
    type: LibraryMemberRef,
    nullable: true,
    args: {
      member_number: t.arg.string({ required: true })
    },
    resolve: safeResolver(async (_parent, args, ctx) => {
      const member = await ctx.deps.membership.queries.getMemberDetails(args.member_number);
      return member === null ? null : libraryMember(member);
    })
  })
}));
