import { builder, safeResolver } from "../schema-builder.js";
import { catalogTitle, inventoryCopy } from "../types.js";
import {
  AddInventoryCopyInputRef,
  CatalogTitleRef,
  CreateCatalogTitleInputRef,
  InventoryCopyRef
} from "./types.js";

builder.mutationFields((t) => ({
  createBook: t.field({
    type: CatalogTitleRef,
    args: {
      input: t.arg({ type: CreateCatalogTitleInputRef, required: true })
    },
    resolve: safeResolver(async (_parent, args, ctx) => {
      const book = await ctx.deps.catalog.commands.addBook(args.input);
      return catalogTitle(book);
    })
  }),
  addBookCopy: t.field({
    type: InventoryCopyRef,
    args: {
      input: t.arg({ type: AddInventoryCopyInputRef, required: true })
    },
    resolve: safeResolver(async (_parent, args, ctx) => {
      const copy = await ctx.deps.catalog.commands.addBookCopy(args.input);
      return inventoryCopy(copy);
    })
  }),
  markBookCopyLost: t.field({
    type: InventoryCopyRef,
    args: {
      barcode: t.arg.string({ required: true })
    },
    resolve: safeResolver(async (_parent, args, ctx) => {
      const copy = await ctx.deps.catalog.commands.markBookCopyLost(args.barcode);
      return inventoryCopy(copy);
    })
  }),
  markBookCopyFound: t.field({
    type: InventoryCopyRef,
    args: {
      barcode: t.arg.string({ required: true })
    },
    resolve: safeResolver(async (_parent, args, ctx) => {
      const copy = await ctx.deps.catalog.commands.markBookCopyFound(args.barcode);
      return inventoryCopy(copy);
    })
  }),
  sendBookCopyToMaintenance: t.field({
    type: InventoryCopyRef,
    args: {
      barcode: t.arg.string({ required: true })
    },
    resolve: safeResolver(async (_parent, args, ctx) => {
      const copy = await ctx.deps.catalog.commands.sendBookCopyToMaintenance(args.barcode);
      return inventoryCopy(copy);
    })
  }),
  completeBookCopyMaintenance: t.field({
    type: InventoryCopyRef,
    args: {
      barcode: t.arg.string({ required: true })
    },
    resolve: safeResolver(async (_parent, args, ctx) => {
      const copy = await ctx.deps.catalog.commands.completeBookCopyMaintenance(args.barcode);
      return inventoryCopy(copy);
    })
  })
}));
