package loan

import (
	"time"

	domainloan "github.com/christophercaldwell/model-architecture/go/internal/domain/loan"
)

type LoanResponse struct {
	Ident      string     `json:"ident"`
	DtCreated  time.Time  `json:"dt_created"`
	DtModified time.Time  `json:"dt_modified"`
	DtDue      *time.Time `json:"dt_due"`
	DtReturned *time.Time `json:"dt_returned"`
}

func LoanToResponse(l domainloan.Loan) LoanResponse {
	return LoanResponse{
		Ident:      string(l.Ident),
		DtCreated:  l.DtCreated,
		DtModified: l.DtModified,
		DtDue:      l.DtDue,
		DtReturned: l.DtReturned,
	}
}

type CreateLoanRequest struct {
	MemberIdent     string `json:"member_ident"`
	BookCopyBarcode string `json:"book_copy_barcode"`
}
