package book

import (
	"context"
	"fmt"

	domain "github.com/christophercaldwell/model-architecture/go/internal/domain/book"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReadRepo struct {
	pool *pgxpool.Pool
}

func NewReadRepo(pool *pgxpool.Pool) *ReadRepo {
	return &ReadRepo{pool: pool}
}

func (r *ReadRepo) GetCatalog(ctx context.Context) ([]domain.Book, error) {
	const q = `
		SELECT b.book_id, b.isbn, b.dt_created, b.dt_modified, b.title, b.author_name
		FROM library.book b
		ORDER BY b.title, b.book_id`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("fetch book catalog: %w", err)
	}
	defer rows.Close()

	var books []domain.Book
	for rows.Next() {
		var row bookRow
		if err := rows.Scan(&row.BookID, &row.ISBN, &row.DtCreated, &row.DtModified, &row.Title, &row.AuthorName); err != nil {
			return nil, fmt.Errorf("scan book row: %w", err)
		}
		books = append(books, rowToDomain(row))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate book rows: %w", err)
	}
	if books == nil {
		books = []domain.Book{}
	}
	return books, nil
}

func (r *ReadRepo) GetByISBN(ctx context.Context, isbn string) (*domain.Book, error) {
	const q = `
		SELECT b.book_id, b.isbn, b.dt_created, b.dt_modified, b.title, b.author_name
		FROM library.book b
		WHERE b.isbn = $1`

	var row bookRow
	err := r.pool.QueryRow(ctx, q, isbn).Scan(
		&row.BookID, &row.ISBN, &row.DtCreated, &row.DtModified, &row.Title, &row.AuthorName,
	)
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("fetch book by isbn: %w", err)
	}
	b := rowToDomain(row)
	return &b, nil
}
