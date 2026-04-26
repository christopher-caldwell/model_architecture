import type { BookId } from "./book.js";

export type BookCopyId = number;

export const bookCopyStatuses = ["active", "maintenance", "lost"] as const;
export type BookCopyStatus = (typeof bookCopyStatuses)[number];

export interface BookCopy {
  id: BookCopyId;
  barcode: string;
  dt_created: Date;
  dt_modified: Date;
  book_id: BookId;
  status: BookCopyStatus;
}

export interface BookCopyCreationPayload {
  barcode: string;
  book_id: BookId;
}

export interface BookCopyPrepared {
  barcode: string;
  book_id: BookId;
  status: BookCopyStatus;
}

export function prepareBookCopy(payload: BookCopyCreationPayload): BookCopyPrepared {
  return {
    barcode: payload.barcode,
    book_id: payload.book_id,
    status: "active"
  };
}

export function ensureBookCopyCanBeBorrowed(copy: BookCopy): void {
  if (copy.status !== "active") {
    throw new BookCopyCannotBeBorrowedError();
  }
}

export function sendBookCopyToMaintenance(copy: BookCopy): BookCopyStatus {
  if (copy.status !== "active") {
    throw new BookCopyCannotBeSentToMaintenanceError();
  }

  return "maintenance";
}

export function completeBookCopyMaintenance(copy: BookCopy): BookCopyStatus {
  if (copy.status !== "maintenance") {
    throw new BookCopyCannotBeReturnedFromMaintenanceError();
  }

  return "active";
}

export function markBookCopyLost(copy: BookCopy): BookCopyStatus {
  if (copy.status === "lost") {
    throw new BookCopyCannotMarkLostError();
  }

  return "lost";
}

export function markBookCopyFound(copy: BookCopy): BookCopyStatus {
  if (copy.status !== "lost") {
    throw new BookCopyCannotBeReturnedFromLostError();
  }

  return "active";
}

export class BookCopyNotFoundError extends Error {
  override readonly name = "BookCopyNotFoundError";

  constructor() {
    super("Book copy not found");
  }
}

export class BookCopyCannotBeBorrowedError extends Error {
  override readonly name = "BookCopyCannotBeBorrowedError";

  constructor() {
    super("Book cannot currently be borrowed");
  }
}

export class BookCopyCannotBeSentToMaintenanceError extends Error {
  override readonly name = "BookCopyCannotBeSentToMaintenanceError";

  constructor() {
    super("Book is not active and cannot be sent to maintenance");
  }
}

export class BookCopyCannotBeReturnedFromMaintenanceError extends Error {
  override readonly name = "BookCopyCannotBeReturnedFromMaintenanceError";

  constructor() {
    super("Book is not currently in maintenance, and therefore cannot be returned");
  }
}

export class BookCopyCannotMarkLostError extends Error {
  override readonly name = "BookCopyCannotMarkLostError";

  constructor() {
    super("Book is already marked lost");
  }
}

export class BookCopyCannotBeReturnedFromLostError extends Error {
  override readonly name = "BookCopyCannotBeReturnedFromLostError";

  constructor() {
    super("Book is not currently lost, and cannot be returned from lost");
  }
}

export interface BookCopyWriteRepository {
  create(insert: BookCopyPrepared): Promise<BookCopy>;
  getByBarcodeForUpdate(barcode: string): Promise<BookCopy | null>;
  updateStatus(id: BookCopyId, status: BookCopyStatus): Promise<void>;
}

export interface BookCopyReadRepository {
  getById(id: BookCopyId): Promise<BookCopy | null>;
  getByBarcode(barcode: string): Promise<BookCopy | null>;
}
