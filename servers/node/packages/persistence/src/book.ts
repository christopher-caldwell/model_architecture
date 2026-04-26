import type { Pool, PoolClient } from "pg";

import type { Book, BookPrepared, BookReadRepository, BookWriteRepository } from "@library/domain";

import type { BookRow } from "./mappers.js";
import { mapBook } from "./mappers.js";
import { sql } from "./sql.js";

interface BookCreateRow {
  book_id: number;
}

export class BookReadRepositoryPostgres implements BookReadRepository {
  constructor(private readonly pool: Pool) {}

  async getCatalog(): Promise<Book[]> {
    const result = await this.pool.query<BookRow>(sql.book.getCatalog);
    return result.rows.map(mapBook);
  }

  async getByIsbn(isbn: string): Promise<Book | null> {
    const result = await this.pool.query<BookRow>(sql.book.getByIsbn, [isbn]);
    return result.rows[0] === undefined ? null : mapBook(result.rows[0]);
  }
}

export class BookWriteRepositoryPostgres implements BookWriteRepository {
  constructor(private readonly client: PoolClient) {}

  async create(insert: BookPrepared): Promise<Book> {
    const result = await this.client.query<BookCreateRow>(sql.book.create, [
      insert.isbn,
      insert.title,
      insert.author_name
    ]);
    const created = result.rows[0];
    if (created === undefined) throw new Error("Failed to create book");

    const now = new Date();
    return {
      id: created.book_id,
      isbn: insert.isbn,
      dt_created: now,
      dt_modified: now,
      title: insert.title,
      author_name: insert.author_name
    };
  }

  async getByIsbn(isbn: string): Promise<Book | null> {
    const result = await this.client.query<BookRow>(sql.book.getByIsbn, [isbn]);
    return result.rows[0] === undefined ? null : mapBook(result.rows[0]);
  }
}
