import { builder } from "../schema-builder.js";
import type { LibraryMember } from "../types.js";

export interface RegisterMemberInput {
  full_name: string;
  max_active_loans: number;
}

export const LibraryMemberStatusRef = builder.enumType("LibraryMemberStatus", {
  values: ["ACTIVE", "SUSPENDED"] as const
});

export const LibraryMemberRef = builder.objectRef<LibraryMember>("LibraryMember").implement({
  fields: (t) => ({
    member_number: t.exposeString("member_number"),
    dt_created: t.exposeString("dt_created"),
    dt_modified: t.exposeString("dt_modified"),
    status: t.field({
      type: LibraryMemberStatusRef,
      resolve: (member) => member.status
    }),
    full_name: t.exposeString("full_name"),
    max_active_loans: t.exposeInt("max_active_loans")
  })
});

type ImplementedInputRef<T extends object> = ReturnType<ReturnType<typeof builder.inputRef<T>>["implement"]>;

export const RegisterMemberInputRef: ImplementedInputRef<RegisterMemberInput> = builder
  .inputRef<RegisterMemberInput>("RegisterMemberInput")
  .implement({
  fields: (t) => ({
    full_name: t.string({ required: true }),
    max_active_loans: t.int({ required: true })
  })
});
