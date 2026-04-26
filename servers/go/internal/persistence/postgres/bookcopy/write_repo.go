package bookcopy

import (
	"context"
	"fmt"
	"time"

	domain "github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
	domainbook "github.com/christophercaldwell/model-architecture/go/internal/domain/book"
	"github.com/jackc/pgx/v5"
)

type WriteRepo struct {
	tx pgx.Tx
}

func NewWriteRepo(tx pgx.Tx) *WriteRepo {
	return &WriteRepo{tx: tx}
}

func (r *WriteRepo) Create(ctx context.Context, prepared domain.BookCopyPrepared) (*domain.BookCopy, error) {
	const q = `
		INSERT INTO library.book_copy (book_id, status_id, barcode)
		VALUES (
			$1,
			(SELECT st.struct_type_id FROM library.struct_type st
			 WHERE st.group_name = 'book_copy_status' AND st.att_pub_ident = $2),
			$3
		)
		RETURNING book_copy_id`

	var bookCopyID int32
	err := r.tx.QueryRow(ctx, q, int32(prepared.BookID), string(prepared.Status), prepared.Barcode).Scan(&bookCopyID)
	if err != nil {
		return nil, fmt.Errorf("create book copy: %w", err)
	}
	now := time.Now()
	bc := domain.BookCopy{
		ID:         domain.BookCopyID(bookCopyID),
		Barcode:    prepared.Barcode,
		DtCreated:  now,
		DtModified: now,
		BookID:     domainbook.BookID(prepared.BookID),
		Status:     prepared.Status,
	}
	return &bc, nil
}

func (r *WriteRepo) GetByBarcodeForUpdate(ctx context.Context, barcode string) (*domain.BookCopy, error) {
	const q = `
		SELECT bc.book_copy_id, bc.barcode, bc.dt_created, bc.dt_modified, bc.book_id, st.att_pub_ident
		FROM library.book_copy bc
		JOIN library.struct_type st ON bc.status_id = st.struct_type_id
		WHERE bc.barcode = $1
		FOR UPDATE OF bc`

	var row bookCopyRow
	err := r.tx.QueryRow(ctx, q, barcode).Scan(
		&row.BookCopyID, &row.Barcode, &row.DtCreated, &row.DtModified, &row.BookID, &row.Status,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("fetch book copy by barcode for update: %w", err)
	}
	return rowToDomain(row)
}

func (r *WriteRepo) UpdateStatus(ctx context.Context, id domain.BookCopyID, status domain.BookCopyStatus) error {
	const q = `
		UPDATE library.book_copy bc
		SET status_id = (
			SELECT st.struct_type_id FROM library.struct_type st
			WHERE st.group_name = 'book_copy_status' AND st.att_pub_ident = $2
		)
		WHERE bc.book_copy_id = $1`

	_, err := r.tx.Exec(ctx, q, int32(id), string(status))
	if err != nil {
		return fmt.Errorf("update book copy status: %w", err)
	}
	return nil
}
