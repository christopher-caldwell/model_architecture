import { builder } from "../schema-builder.js";
import type { CatalogTitle, InventoryCopy } from "../types.js";

export interface CreateCatalogTitleInput {
  isbn: string;
  title: string;
  author_name: string;
}

export interface AddInventoryCopyInput {
  isbn: string;
  barcode: string;
}

export const InventoryCopyStatusRef = builder.enumType("InventoryCopyStatus", {
  values: ["ACTIVE", "MAINTENANCE", "LOST"] as const
});

export const CatalogTitleRef = builder.objectRef<CatalogTitle>("CatalogTitle").implement({
  fields: (t) => ({
    isbn: t.exposeString("isbn"),
    dt_created: t.exposeString("dt_created"),
    dt_modified: t.exposeString("dt_modified"),
    title: t.exposeString("title"),
    author_name: t.exposeString("author_name")
  })
});

export const InventoryCopyRef = builder.objectRef<InventoryCopy>("InventoryCopy").implement({
  fields: (t) => ({
    barcode: t.exposeString("barcode"),
    dt_created: t.exposeString("dt_created"),
    dt_modified: t.exposeString("dt_modified"),
    status: t.field({
      type: InventoryCopyStatusRef,
      resolve: (copy) => copy.status
    })
  })
});

type ImplementedInputRef<T extends object> = ReturnType<ReturnType<typeof builder.inputRef<T>>["implement"]>;

export const CreateCatalogTitleInputRef: ImplementedInputRef<CreateCatalogTitleInput> = builder
  .inputRef<CreateCatalogTitleInput>("CreateCatalogTitleInput")
  .implement({
  fields: (t) => ({
    isbn: t.string({ required: true }),
    title: t.string({ required: true }),
    author_name: t.string({ required: true })
  })
});

export const AddInventoryCopyInputRef: ImplementedInputRef<AddInventoryCopyInput> = builder
  .inputRef<AddInventoryCopyInput>("AddInventoryCopyInput")
  .implement({
  fields: (t) => ({
    isbn: t.string({ required: true }),
    barcode: t.string({ required: true })
  })
});
