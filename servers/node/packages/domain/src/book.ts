export type BookId = number;

export interface Book {
  id: BookId;
  isbn: string;
  dt_created: Date;
  dt_modified: Date;
  title: string;
  author_name: string;
}

export interface BookCreationPayload {
  isbn: string;
  title: string;
  author_name: string;
}

export interface BookPrepared {
  isbn: string;
  title: string;
  author_name: string;
}

export function prepareBook(payload: BookCreationPayload): BookPrepared {
  return {
    isbn: payload.isbn,
    title: payload.title,
    author_name: payload.author_name
  };
}

export class BookNotFoundError extends Error {
  override readonly name = "BookNotFoundError";

  constructor() {
    super("Book not found");
  }
}

export interface BookWriteRepository {
  create(insert: BookPrepared): Promise<Book>;
  getByIsbn(isbn: string): Promise<Book | null>;
}

export interface BookReadRepository {
  getCatalog(): Promise<Book[]>;
  getByIsbn(isbn: string): Promise<Book | null>;
}
