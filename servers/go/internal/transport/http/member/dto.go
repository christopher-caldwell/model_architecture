package member

import (
	"time"

	domainm "github.com/christophercaldwell/model-architecture/go/internal/domain/member"
)

type MemberResponse struct {
	Ident          string    `json:"ident"`
	DtCreated      time.Time `json:"dt_created"`
	DtModified     time.Time `json:"dt_modified"`
	Status         string    `json:"status"`
	FullName       string    `json:"full_name"`
	MaxActiveLoans int16     `json:"max_active_loans"`
}

func memberToResponse(m domainm.Member) MemberResponse {
	return MemberResponse{
		Ident:          string(m.Ident),
		DtCreated:      m.DtCreated,
		DtModified:     m.DtModified,
		Status:         string(m.Status),
		FullName:       m.FullName,
		MaxActiveLoans: m.MaxActiveLoans,
	}
}

type CreateMemberRequest struct {
	FullName       string `json:"full_name"`
	MaxActiveLoans int16  `json:"max_active_loans"`
}
