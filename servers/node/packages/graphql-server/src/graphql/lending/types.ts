import { builder } from "../../schema-builder.js";
import type { LoanRecord } from "../../types.js";

export interface StartLoanInput {
  memberNumber: string;
  barcode: string;
}

export const LoanRecordRef = builder.objectRef<LoanRecord>("LoanRecord").implement({
  fields: (t) => ({
    loanNumber: t.exposeString("loanNumber"),
    createdAt: t.exposeString("createdAt"),
    modifiedAt: t.exposeString("modifiedAt"),
    dueAt: t.exposeString("dueAt", { nullable: true }),
    returnedAt: t.exposeString("returnedAt", { nullable: true })
  })
});

type ImplementedInputRef<T extends object> = ReturnType<ReturnType<typeof builder.inputRef<T>>["implement"]>;

export const StartLoanInputRef: ImplementedInputRef<StartLoanInput> = builder
  .inputRef<StartLoanInput>("StartLoanInput")
  .implement({
  fields: (t) => ({
    memberNumber: t.string({ required: true }),
    barcode: t.string({ required: true })
  })
});
