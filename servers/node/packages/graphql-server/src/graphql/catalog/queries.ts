import { builder, safeResolver } from "../schema-builder.js";
import { catalogTitle, inventoryCopy } from "../types.js";
import { CatalogTitleRef, InventoryCopyRef } from "./types.js";

builder.queryFields((t) => ({
  books: t.field({
    type: [CatalogTitleRef],
    resolve: safeResolver(async (_parent, _args, ctx) => {
      const books = await ctx.deps.catalog.queries.getBookCatalog();
      return books.map(catalogTitle);
    })
  }),
  bookByIsbn: t.field({
    type: CatalogTitleRef,
    nullable: true,
    args: {
      isbn: t.arg.string({ required: true })
    },
    resolve: safeResolver(async (_parent, args, ctx) => {
      const book = await ctx.deps.catalog.queries.getBookByIsbn(args.isbn);
      return book === null ? null : catalogTitle(book);
    })
  }),
  bookCopy: t.field({
    type: InventoryCopyRef,
    nullable: true,
    args: {
      barcode: t.arg.string({ required: true })
    },
    resolve: safeResolver(async (_parent, args, ctx) => {
      const copy = await ctx.deps.catalog.queries.getBookCopyDetails(args.barcode);
      return copy === null ? null : inventoryCopy(copy);
    })
  })
}));
