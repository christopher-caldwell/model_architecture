import {
  BookCopyNotFoundError,
  BookNotFoundError,
  completeBookCopyMaintenance,
  markBookCopyFound,
  markBookCopyLost,
  prepareBook,
  prepareBookCopy,
  sendBookCopyToMaintenance,
  type Book,
  type BookCopy,
  type BookCreationPayload,
  type WriteUnitOfWork,
  type WriteUnitOfWorkFactory
} from "@library/domain";

import type { AddBookCopyInput } from "./inputs.js";
import { inWriteUnitOfWork } from "./uow.js";

export class CatalogCommands {
  constructor(private readonly uowFactory: WriteUnitOfWorkFactory) {}

  async addBook(payload: BookCreationPayload): Promise<Book> {
    const prepared = prepareBook(payload);

    return inWriteUnitOfWork(this.uowFactory, async (uow) => {
      return uow.books.create(prepared);
    });
  }

  async addBookCopy(input: AddBookCopyInput): Promise<BookCopy> {
    return inWriteUnitOfWork(this.uowFactory, async (uow) => {
      const book = await this.getBookByIsbn(uow, input.isbn);
      const prepared = prepareBookCopy({
        barcode: input.barcode,
        book_id: book.id
      });

      return uow.bookCopies.create(prepared);
    });
  }

  async markBookCopyLost(barcode: string): Promise<BookCopy> {
    return this.updateBookCopyStatus(barcode, markBookCopyLost);
  }

  async markBookCopyFound(barcode: string): Promise<BookCopy> {
    return this.updateBookCopyStatus(barcode, markBookCopyFound);
  }

  async sendBookCopyToMaintenance(barcode: string): Promise<BookCopy> {
    return this.updateBookCopyStatus(barcode, sendBookCopyToMaintenance);
  }

  async completeBookCopyMaintenance(barcode: string): Promise<BookCopy> {
    return this.updateBookCopyStatus(barcode, completeBookCopyMaintenance);
  }

  private async updateBookCopyStatus(
    barcode: string,
    transition: (copy: BookCopy) => BookCopy["status"]
  ): Promise<BookCopy> {
    return inWriteUnitOfWork(this.uowFactory, async (uow) => {
      const copy = await this.getBookCopyByBarcode(uow, barcode);
      const status = transition(copy);

      await uow.bookCopies.updateStatus(copy.id, status);

      return {
        ...copy,
        status,
        dt_modified: new Date()
      };
    });
  }

  private async getBookByIsbn(uow: WriteUnitOfWork, isbn: string): Promise<Book> {
    const book = await uow.books.getByIsbn(isbn);
    if (book === null) throw new BookNotFoundError();
    return book;
  }

  private async getBookCopyByBarcode(uow: WriteUnitOfWork, barcode: string): Promise<BookCopy> {
    const copy = await uow.bookCopies.getByBarcodeForUpdate(barcode);
    if (copy === null) throw new BookCopyNotFoundError();
    return copy;
  }
}
