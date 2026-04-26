package queries

import (
	"context"
	"fmt"

	"github.com/christophercaldwell/model-architecture/go/internal/domain/book"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
)

type CatalogQueries struct {
	bookRepo     book.ReadRepository
	bookCopyRepo bookcopy.ReadRepository
}

func NewCatalogQueries(bookRepo book.ReadRepository, bookCopyRepo bookcopy.ReadRepository) *CatalogQueries {
	return &CatalogQueries{
		bookRepo:     bookRepo,
		bookCopyRepo: bookCopyRepo,
	}
}

func (q *CatalogQueries) GetBookCatalog(ctx context.Context) ([]book.Book, error) {
	books, err := q.bookRepo.GetCatalog(ctx)
	if err != nil {
		return nil, fmt.Errorf("get book catalog: %w", err)
	}
	return books, nil
}

func (q *CatalogQueries) GetBookByISBN(ctx context.Context, isbn string) (*book.Book, error) {
	b, err := q.bookRepo.GetByISBN(ctx, isbn)
	if err != nil {
		return nil, fmt.Errorf("get book by isbn: %w", err)
	}
	return b, nil
}

func (q *CatalogQueries) GetBookCopyDetails(ctx context.Context, barcode string) (*bookcopy.BookCopy, error) {
	bc, err := q.bookCopyRepo.GetByBarcode(ctx, barcode)
	if err != nil {
		return nil, fmt.Errorf("get book copy details: %w", err)
	}
	return bc, nil
}
