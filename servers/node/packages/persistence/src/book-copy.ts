import type { Pool, PoolClient } from "pg";

import type {
  BookCopy,
  BookCopyId,
  BookCopyPrepared,
  BookCopyReadRepository,
  BookCopyStatus,
  BookCopyWriteRepository
} from "@library/domain";

import type { BookCopyRow } from "./mappers.js";
import { mapBookCopy } from "./mappers.js";
import { sql } from "./sql.js";

interface BookCopyCreateRow {
  book_copy_id: number;
}

export class BookCopyReadRepositoryPostgres implements BookCopyReadRepository {
  constructor(private readonly pool: Pool) {}

  async getById(id: BookCopyId): Promise<BookCopy | null> {
    const result = await this.pool.query<BookCopyRow>(sql.bookCopy.getById, [id]);
    return result.rows[0] === undefined ? null : mapBookCopy(result.rows[0]);
  }

  async getByBarcode(barcode: string): Promise<BookCopy | null> {
    const result = await this.pool.query<BookCopyRow>(sql.bookCopy.getByBarcode, [barcode]);
    return result.rows[0] === undefined ? null : mapBookCopy(result.rows[0]);
  }
}

export class BookCopyWriteRepositoryPostgres implements BookCopyWriteRepository {
  constructor(private readonly client: PoolClient) {}

  async create(insert: BookCopyPrepared): Promise<BookCopy> {
    const result = await this.client.query<BookCopyCreateRow>(sql.bookCopy.create, [
      insert.book_id,
      insert.status,
      insert.barcode
    ]);
    const created = result.rows[0];
    if (created === undefined) throw new Error("Failed to create book copy");

    const now = new Date();
    return {
      id: created.book_copy_id,
      barcode: insert.barcode,
      dt_created: now,
      dt_modified: now,
      book_id: insert.book_id,
      status: insert.status
    };
  }

  async getByBarcodeForUpdate(barcode: string): Promise<BookCopy | null> {
    const result = await this.client.query<BookCopyRow>(sql.bookCopy.getByBarcodeForUpdate, [barcode]);
    return result.rows[0] === undefined ? null : mapBookCopy(result.rows[0]);
  }

  async updateStatus(id: BookCopyId, status: BookCopyStatus): Promise<void> {
    await this.client.query(sql.bookCopy.updateStatus, [id, status]);
  }
}
