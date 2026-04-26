package application

import (
	"context"

	"github.com/christophercaldwell/model-architecture/go/internal/domain/book"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/loan"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/member"
)

type UnitOfWork interface {
	Books() book.WriteRepository
	BookCopies() bookcopy.WriteRepository
	Members() member.WriteRepository
	Loans() loan.WriteRepository
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type UnitOfWorkFactory interface {
	New(ctx context.Context) (UnitOfWork, error)
}
