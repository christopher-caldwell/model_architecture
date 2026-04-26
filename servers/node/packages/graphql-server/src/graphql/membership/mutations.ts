import { builder, safeResolver } from "../schema-builder.js";
import { libraryMember } from "../types.js";
import { LibraryMemberRef, RegisterMemberInputRef } from "./types.js";

builder.mutationFields((t) => ({
  registerMember: t.field({
    type: LibraryMemberRef,
    args: {
      input: t.arg({ type: RegisterMemberInputRef, required: true })
    },
    resolve: safeResolver(async (_parent, args, ctx) => {
      const member = await ctx.deps.membership.commands.registerMember(args.input);
      return libraryMember(member);
    })
  }),
  suspendMember: t.field({
    type: LibraryMemberRef,
    args: {
      member_number: t.arg.string({ required: true })
    },
    resolve: safeResolver(async (_parent, args, ctx) => {
      const member = await ctx.deps.membership.commands.suspendMember({ member_ident: args.member_number });
      return libraryMember(member);
    })
  }),
  reactivateMember: t.field({
    type: LibraryMemberRef,
    args: {
      member_number: t.arg.string({ required: true })
    },
    resolve: safeResolver(async (_parent, args, ctx) => {
      const member = await ctx.deps.membership.commands.reactivateMember({ member_ident: args.member_number });
      return libraryMember(member);
    })
  })
}));
