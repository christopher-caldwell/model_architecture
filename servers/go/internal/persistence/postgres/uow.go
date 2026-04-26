package postgres

import (
	"context"
	"fmt"

	"github.com/christophercaldwell/model-architecture/go/internal/application"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/book"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/loan"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/member"
	pgbook "github.com/christophercaldwell/model-architecture/go/internal/persistence/postgres/book"
	pgbookcopy "github.com/christophercaldwell/model-architecture/go/internal/persistence/postgres/bookcopy"
	pgloan "github.com/christophercaldwell/model-architecture/go/internal/persistence/postgres/loan"
	pgmember "github.com/christophercaldwell/model-architecture/go/internal/persistence/postgres/member"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type sqlUnitOfWork struct {
	tx         pgx.Tx
	books      book.WriteRepository
	bookCopies bookcopy.WriteRepository
	members    member.WriteRepository
	loans      loan.WriteRepository
}

func (u *sqlUnitOfWork) Books() book.WriteRepository      { return u.books }
func (u *sqlUnitOfWork) BookCopies() bookcopy.WriteRepository { return u.bookCopies }
func (u *sqlUnitOfWork) Members() member.WriteRepository  { return u.members }
func (u *sqlUnitOfWork) Loans() loan.WriteRepository      { return u.loans }

func (u *sqlUnitOfWork) Commit(ctx context.Context) error {
	if err := u.tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func (u *sqlUnitOfWork) Rollback(ctx context.Context) error {
	return u.tx.Rollback(ctx)
}

type SqlUnitOfWorkFactory struct {
	Pool *pgxpool.Pool
}

func (f *SqlUnitOfWorkFactory) New(ctx context.Context) (application.UnitOfWork, error) {
	tx, err := f.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	return &sqlUnitOfWork{
		tx:         tx,
		books:      pgbook.NewWriteRepo(tx),
		bookCopies: pgbookcopy.NewWriteRepo(tx),
		members:    pgmember.NewWriteRepo(tx),
		loans:      pgloan.NewWriteRepo(tx),
	}, nil
}
