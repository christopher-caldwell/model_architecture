import type {
  Book,
  BookCopy,
  BookCopyReadRepository,
  BookReadRepository
} from "@library/domain";

export class CatalogQueries {
  constructor(
    private readonly bookReadRepository: BookReadRepository,
    private readonly bookCopyReadRepository: BookCopyReadRepository
  ) {}

  getBookCatalog(): Promise<Book[]> {
    return this.bookReadRepository.getCatalog();
  }

  getBookByIsbn(isbn: string): Promise<Book | null> {
    return this.bookReadRepository.getByIsbn(isbn);
  }

  getBookCopyDetails(barcode: string): Promise<BookCopy | null> {
    return this.bookCopyReadRepository.getByBarcode(barcode);
  }
}
