import type { BookWriteRepository } from "./book.js";
import type { BookCopyWriteRepository } from "./book-copy.js";
import type { LoanWriteRepository } from "./loan.js";
import type { MemberWriteRepository } from "./member.js";

export interface WriteUnitOfWork {
  readonly books: BookWriteRepository;
  readonly bookCopies: BookCopyWriteRepository;
  readonly members: MemberWriteRepository;
  readonly loans: LoanWriteRepository;
  commit(): Promise<void>;
  rollback(): Promise<void>;
}

export interface WriteUnitOfWorkFactory {
  build(): Promise<WriteUnitOfWork>;
}
