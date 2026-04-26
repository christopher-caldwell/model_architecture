package bookcopy

import (
	"context"
	"fmt"

	domain "github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReadRepo struct {
	pool *pgxpool.Pool
}

func NewReadRepo(pool *pgxpool.Pool) *ReadRepo {
	return &ReadRepo{pool: pool}
}

const selectBookCopyCols = `
	bc.book_copy_id, bc.barcode, bc.dt_created, bc.dt_modified, bc.book_id, st.att_pub_ident`

const bookCopyJoin = `
	FROM library.book_copy bc
	JOIN library.struct_type st ON bc.status_id = st.struct_type_id`

func scanRow(row pgx.Row) (*domain.BookCopy, error) {
	var r bookCopyRow
	if err := row.Scan(&r.BookCopyID, &r.Barcode, &r.DtCreated, &r.DtModified, &r.BookID, &r.Status); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan book copy row: %w", err)
	}
	return rowToDomain(r)
}

func (r *ReadRepo) GetByID(ctx context.Context, id domain.BookCopyID) (*domain.BookCopy, error) {
	q := `SELECT ` + selectBookCopyCols + bookCopyJoin + ` WHERE bc.book_copy_id = $1`
	bc, err := scanRow(r.pool.QueryRow(ctx, q, int32(id)))
	if err != nil {
		return nil, fmt.Errorf("fetch book copy by id: %w", err)
	}
	return bc, nil
}

func (r *ReadRepo) GetByBarcode(ctx context.Context, barcode string) (*domain.BookCopy, error) {
	q := `SELECT ` + selectBookCopyCols + bookCopyJoin + ` WHERE bc.barcode = $1`
	bc, err := scanRow(r.pool.QueryRow(ctx, q, barcode))
	if err != nil {
		return nil, fmt.Errorf("fetch book copy by barcode: %w", err)
	}
	return bc, nil
}
