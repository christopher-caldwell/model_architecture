package book

import "context"

type ReadRepository interface {
	GetCatalog(ctx context.Context) ([]Book, error)
	GetByISBN(ctx context.Context, isbn string) (*Book, error)
}

type WriteRepository interface {
	Create(ctx context.Context, prepared BookPrepared) (*Book, error)
	GetByISBN(ctx context.Context, isbn string) (*Book, error)
}
