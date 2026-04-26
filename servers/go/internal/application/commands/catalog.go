package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/christophercaldwell/model-architecture/go/internal/application"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/book"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
)

type AddBookCopyInput struct {
	ISBN    string
	Barcode string
}

type CatalogCommands struct {
	factory application.UnitOfWorkFactory
}

func NewCatalogCommands(factory application.UnitOfWorkFactory) *CatalogCommands {
	return &CatalogCommands{factory: factory}
}

func (c *CatalogCommands) getBookByISBN(ctx context.Context, uow application.UnitOfWork, isbn string) (*book.Book, error) {
	b, err := uow.Books().GetByISBN(ctx, isbn)
	if err != nil {
		return nil, fmt.Errorf("load book for write: %w", err)
	}
	if b == nil {
		return nil, book.ErrNotFound
	}
	return b, nil
}

func (c *CatalogCommands) getBookCopyByBarcode(ctx context.Context, uow application.UnitOfWork, barcode string) (*bookcopy.BookCopy, error) {
	bc, err := uow.BookCopies().GetByBarcodeForUpdate(ctx, barcode)
	if err != nil {
		return nil, fmt.Errorf("load book copy for write: %w", err)
	}
	if bc == nil {
		return nil, bookcopy.ErrNotFound
	}
	return bc, nil
}

func (c *CatalogCommands) AddBook(ctx context.Context, payload book.BookCreationPayload) (*book.Book, error) {
	prepared := payload.Prepare()
	uow, err := c.factory.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("build unit of work: %w", err)
	}
	defer uow.Rollback(ctx) //nolint:errcheck

	result, err := uow.Books().Create(ctx, prepared)
	if err != nil {
		return nil, fmt.Errorf("add book: %w", err)
	}
	if err := uow.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}

func (c *CatalogCommands) AddBookCopy(ctx context.Context, input AddBookCopyInput) (*bookcopy.BookCopy, error) {
	uow, err := c.factory.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("build unit of work: %w", err)
	}
	defer uow.Rollback(ctx) //nolint:errcheck

	b, err := c.getBookByISBN(ctx, uow, input.ISBN)
	if err != nil {
		return nil, err
	}
	prepared := bookcopy.BookCopyCreationPayload{
		Barcode: input.Barcode,
		BookID:  b.ID,
	}.Prepare()
	result, err := uow.BookCopies().Create(ctx, prepared)
	if err != nil {
		return nil, fmt.Errorf("add book copy: %w", err)
	}
	if err := uow.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}

func (c *CatalogCommands) MarkBookCopyLost(ctx context.Context, barcode string) (*bookcopy.BookCopy, error) {
	uow, err := c.factory.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("build unit of work: %w", err)
	}
	defer uow.Rollback(ctx) //nolint:errcheck

	bc, err := c.getBookCopyByBarcode(ctx, uow, barcode)
	if err != nil {
		return nil, err
	}
	newStatus, err := bc.MarkLost()
	if err != nil {
		return nil, err
	}
	if err := uow.BookCopies().UpdateStatus(ctx, bc.ID, newStatus); err != nil {
		return nil, fmt.Errorf("mark book copy lost: %w", err)
	}
	if err := uow.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	updated := *bc
	updated.Status = newStatus
	updated.DtModified = time.Now()
	return &updated, nil
}

func (c *CatalogCommands) MarkBookCopyFound(ctx context.Context, barcode string) (*bookcopy.BookCopy, error) {
	uow, err := c.factory.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("build unit of work: %w", err)
	}
	defer uow.Rollback(ctx) //nolint:errcheck

	bc, err := c.getBookCopyByBarcode(ctx, uow, barcode)
	if err != nil {
		return nil, err
	}
	newStatus, err := bc.MarkFound()
	if err != nil {
		return nil, err
	}
	if err := uow.BookCopies().UpdateStatus(ctx, bc.ID, newStatus); err != nil {
		return nil, fmt.Errorf("mark book copy found: %w", err)
	}
	if err := uow.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	updated := *bc
	updated.Status = newStatus
	updated.DtModified = time.Now()
	return &updated, nil
}

func (c *CatalogCommands) SendBookCopyToMaintenance(ctx context.Context, barcode string) (*bookcopy.BookCopy, error) {
	uow, err := c.factory.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("build unit of work: %w", err)
	}
	defer uow.Rollback(ctx) //nolint:errcheck

	bc, err := c.getBookCopyByBarcode(ctx, uow, barcode)
	if err != nil {
		return nil, err
	}
	newStatus, err := bc.SendToMaintenance()
	if err != nil {
		return nil, err
	}
	if err := uow.BookCopies().UpdateStatus(ctx, bc.ID, newStatus); err != nil {
		return nil, fmt.Errorf("send book copy to maintenance: %w", err)
	}
	if err := uow.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	updated := *bc
	updated.Status = newStatus
	updated.DtModified = time.Now()
	return &updated, nil
}

func (c *CatalogCommands) CompleteBookCopyMaintenance(ctx context.Context, barcode string) (*bookcopy.BookCopy, error) {
	uow, err := c.factory.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("build unit of work: %w", err)
	}
	defer uow.Rollback(ctx) //nolint:errcheck

	bc, err := c.getBookCopyByBarcode(ctx, uow, barcode)
	if err != nil {
		return nil, err
	}
	newStatus, err := bc.CompleteMaintenance()
	if err != nil {
		return nil, err
	}
	if err := uow.BookCopies().UpdateStatus(ctx, bc.ID, newStatus); err != nil {
		return nil, fmt.Errorf("complete book copy maintenance: %w", err)
	}
	if err := uow.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	updated := *bc
	updated.Status = newStatus
	updated.DtModified = time.Now()
	return &updated, nil
}
