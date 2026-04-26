package queries

import (
	"context"
	"fmt"

	"github.com/christophercaldwell/model-architecture/go/internal/domain/member"
)

type MembershipQueries struct {
	memberRepo member.ReadRepository
}

func NewMembershipQueries(memberRepo member.ReadRepository) *MembershipQueries {
	return &MembershipQueries{memberRepo: memberRepo}
}

func (q *MembershipQueries) GetMemberDetails(ctx context.Context, ident member.MemberIdent) (*member.Member, error) {
	m, err := q.memberRepo.GetByIdent(ctx, ident)
	if err != nil {
		return nil, fmt.Errorf("get member details: %w", err)
	}
	return m, nil
}
