package bookcopy

import "context"

type ReadRepository interface {
	GetByID(ctx context.Context, id BookCopyID) (*BookCopy, error)
	GetByBarcode(ctx context.Context, barcode string) (*BookCopy, error)
}

type WriteRepository interface {
	Create(ctx context.Context, prepared BookCopyPrepared) (*BookCopy, error)
	GetByBarcodeForUpdate(ctx context.Context, barcode string) (*BookCopy, error)
	UpdateStatus(ctx context.Context, id BookCopyID, status BookCopyStatus) error
}
