package queries

import (
	"context"
	"fmt"

	"github.com/christophercaldwell/model-architecture/go/internal/domain/loan"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/member"
)

type LendingQueries struct {
	loanRepo loan.ReadRepository
}

func NewLendingQueries(loanRepo loan.ReadRepository) *LendingQueries {
	return &LendingQueries{loanRepo: loanRepo}
}

func (q *LendingQueries) GetMemberLoans(ctx context.Context, ident member.MemberIdent) ([]loan.Loan, error) {
	loans, err := q.loanRepo.GetByMemberIdent(ctx, ident)
	if err != nil {
		return nil, fmt.Errorf("get member loans: %w", err)
	}
	return loans, nil
}

func (q *LendingQueries) GetOverdueLoans(ctx context.Context) ([]loan.Loan, error) {
	loans, err := q.loanRepo.GetOverdue(ctx)
	if err != nil {
		return nil, fmt.Errorf("get overdue loans: %w", err)
	}
	return loans, nil
}
