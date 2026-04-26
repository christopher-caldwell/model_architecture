import type { Pool, PoolClient } from "pg";

import type {
  BookCopyWriteRepository,
  BookWriteRepository,
  LoanWriteRepository,
  MemberWriteRepository,
  WriteUnitOfWork,
  WriteUnitOfWorkFactory
} from "@library/domain";

import { BookWriteRepositoryPostgres } from "./book.js";
import { BookCopyWriteRepositoryPostgres } from "./book-copy.js";
import { LoanWriteRepositoryPostgres } from "./loan.js";
import { MemberWriteRepositoryPostgres } from "./member.js";

export class SqlWriteUnitOfWork implements WriteUnitOfWork {
  readonly books: BookWriteRepository;
  readonly bookCopies: BookCopyWriteRepository;
  readonly members: MemberWriteRepository;
  readonly loans: LoanWriteRepository;

  private consumed = false;

  constructor(private readonly client: PoolClient) {
    this.books = new BookWriteRepositoryPostgres(client);
    this.bookCopies = new BookCopyWriteRepositoryPostgres(client);
    this.members = new MemberWriteRepositoryPostgres(client);
    this.loans = new LoanWriteRepositoryPostgres(client);
  }

  async commit(): Promise<void> {
    await this.finish("COMMIT");
  }

  async rollback(): Promise<void> {
    await this.finish("ROLLBACK");
  }

  private async finish(statement: "COMMIT" | "ROLLBACK"): Promise<void> {
    if (this.consumed) return;
    this.consumed = true;

    try {
      await this.client.query(statement);
    } finally {
      this.client.release();
    }
  }
}

export class SqlWriteUnitOfWorkFactory implements WriteUnitOfWorkFactory {
  constructor(private readonly pool: Pool) {}

  async build(): Promise<WriteUnitOfWork> {
    const client = await this.pool.connect();

    try {
      await client.query("BEGIN");
      return new SqlWriteUnitOfWork(client);
    } catch (error) {
      client.release();
      throw error;
    }
  }
}
