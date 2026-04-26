package book

import (
	"context"
	"fmt"
	"time"

	domain "github.com/christophercaldwell/model-architecture/go/internal/domain/book"
	"github.com/jackc/pgx/v5"
)

type WriteRepo struct {
	tx pgx.Tx
}

func NewWriteRepo(tx pgx.Tx) *WriteRepo {
	return &WriteRepo{tx: tx}
}

func (r *WriteRepo) Create(ctx context.Context, prepared domain.BookPrepared) (*domain.Book, error) {
	const q = `
		INSERT INTO library.book (isbn, title, author_name)
		VALUES ($1, $2, $3)
		RETURNING book_id`

	var bookID int32
	err := r.tx.QueryRow(ctx, q, prepared.ISBN, prepared.Title, prepared.AuthorName).Scan(&bookID)
	if err != nil {
		return nil, fmt.Errorf("create book: %w", err)
	}
	now := time.Now()
	b := domain.Book{
		ID:         domain.BookID(bookID),
		ISBN:       prepared.ISBN,
		DtCreated:  now,
		DtModified: now,
		Title:      prepared.Title,
		AuthorName: prepared.AuthorName,
	}
	return &b, nil
}

func (r *WriteRepo) GetByISBN(ctx context.Context, isbn string) (*domain.Book, error) {
	const q = `
		SELECT b.book_id, b.isbn, b.dt_created, b.dt_modified, b.title, b.author_name
		FROM library.book b
		WHERE b.isbn = $1`

	var row bookRow
	err := r.tx.QueryRow(ctx, q, isbn).Scan(
		&row.BookID, &row.ISBN, &row.DtCreated, &row.DtModified, &row.Title, &row.AuthorName,
	)
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("fetch book by isbn for write: %w", err)
	}
	b := rowToDomain(row)
	return &b, nil
}

func isNoRows(err error) bool {
	return err == pgx.ErrNoRows
}
