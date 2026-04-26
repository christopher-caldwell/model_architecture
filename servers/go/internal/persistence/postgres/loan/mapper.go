package loan

import (
	domainbc "github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
	domain "github.com/christophercaldwell/model-architecture/go/internal/domain/loan"
	domainm "github.com/christophercaldwell/model-architecture/go/internal/domain/member"
)

func rowToDomain(r loanRow) domain.Loan {
	return domain.Loan{
		ID:         domain.LoanID(r.LoanID),
		Ident:      domain.LoanIdent(r.LoanIdent),
		DtCreated:  r.DtCreated,
		DtModified: r.DtModified,
		BookCopyID: domainbc.BookCopyID(r.BookCopyID),
		MemberID:   domainm.MemberID(r.MemberID),
		DtDue:      r.DtDue,
		DtReturned: r.DtReturned,
	}
}
