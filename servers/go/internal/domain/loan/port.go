package loan

import (
	"context"

	"github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/member"
)

type ReadRepository interface {
	GetByMemberIdent(ctx context.Context, ident member.MemberIdent) ([]Loan, error)
	GetOverdue(ctx context.Context) ([]Loan, error)
	FindActiveByBookCopyID(ctx context.Context, id bookcopy.BookCopyID) (*Loan, error)
	CountActiveByMemberID(ctx context.Context, id member.MemberID) (int64, error)
}

type WriteRepository interface {
	Create(ctx context.Context, prepared LoanPrepared) (*Loan, error)
	End(ctx context.Context, id LoanID) error
	FindActiveByBookCopyIDForUpdate(ctx context.Context, id bookcopy.BookCopyID) (*Loan, error)
	CountActiveByMemberID(ctx context.Context, id member.MemberID) (int64, error)
}
