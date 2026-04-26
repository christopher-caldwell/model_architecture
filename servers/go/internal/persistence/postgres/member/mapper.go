package member

import (
	"fmt"

	domain "github.com/christophercaldwell/model-architecture/go/internal/domain/member"
)

func rowToDomain(r memberRow) (*domain.Member, error) {
	status, err := parseStatus(r.Status)
	if err != nil {
		return nil, err
	}
	return &domain.Member{
		ID:             domain.MemberID(r.MemberID),
		Ident:          domain.MemberIdent(r.MemberIdent),
		DtCreated:      r.DtCreated,
		DtModified:     r.DtModified,
		Status:         status,
		FullName:       r.FullName,
		MaxActiveLoans: r.MaxActiveLoans,
	}, nil
}

func parseStatus(s string) (domain.MemberStatus, error) {
	switch domain.MemberStatus(s) {
	case domain.MemberStatusActive, domain.MemberStatusSuspended:
		return domain.MemberStatus(s), nil
	default:
		return "", fmt.Errorf("unknown member status in DB: %q", s)
	}
}
