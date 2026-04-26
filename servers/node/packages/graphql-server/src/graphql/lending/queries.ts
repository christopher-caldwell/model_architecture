import { builder, safeResolver } from "../../schema-builder.js";
import { loanRecord } from "../../types.js";
import { LoanRecordRef } from "./types.js";

builder.queryFields((t) => ({
  memberLoans: t.field({
    type: [LoanRecordRef],
    args: {
      memberNumber: t.arg.string({ required: true })
    },
    resolve: safeResolver(async (_parent, args, ctx) => {
      const loans = await ctx.deps.lending.queries.getMemberLoans(args.memberNumber);
      return loans.map(loanRecord);
    })
  }),
  overdueLoans: t.field({
    type: [LoanRecordRef],
    resolve: safeResolver(async (_parent, _args, ctx) => {
      const loans = await ctx.deps.lending.queries.getOverdueLoans();
      return loans.map(loanRecord);
    })
  })
}));
