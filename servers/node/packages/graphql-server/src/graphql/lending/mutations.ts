import { InventoryCopyRef } from "../../catalog/types.js";
import { builder, safeResolver } from "../../schema-builder.js";
import { inventoryCopy, loanRecord } from "../../types.js";
import { LoanRecordRef, StartLoanInputRef } from "./types.js";

builder.mutationFields((t) => ({
  checkOutBookCopy: t.field({
    type: LoanRecordRef,
    args: {
      input: t.arg({ type: StartLoanInputRef, required: true })
    },
    resolve: safeResolver(async (_parent, args, ctx) => {
      const loan = await ctx.deps.lending.commands.checkOutBookCopy({
        memberIdent: args.input.memberNumber,
        bookCopyBarcode: args.input.barcode
      });
      return loanRecord(loan);
    })
  }),
  returnBookCopy: t.field({
    type: LoanRecordRef,
    args: {
      barcode: t.arg.string({ required: true })
    },
    resolve: safeResolver(async (_parent, args, ctx) => {
      const loan = await ctx.deps.lending.commands.returnBookCopy(args.barcode);
      return loanRecord(loan);
    })
  }),
  reportLostLoanedBookCopy: t.field({
    type: InventoryCopyRef,
    args: {
      barcode: t.arg.string({ required: true })
    },
    resolve: safeResolver(async (_parent, args, ctx) => {
      const copy = await ctx.deps.lending.commands.reportLostLoanedBookCopy(args.barcode);
      return inventoryCopy(copy);
    })
  })
}));
