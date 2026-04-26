package loan

import (
	"time"

	"github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/member"
)

type LoanID int32

type LoanIdent string

type Loan struct {
	ID         LoanID
	Ident      LoanIdent
	DtCreated  time.Time
	DtModified time.Time
	BookCopyID bookcopy.BookCopyID
	MemberID   member.MemberID
	DtDue      *time.Time
	DtReturned *time.Time
}

type LoanCreationPayload struct {
	MemberID   member.MemberID
	BookCopyID bookcopy.BookCopyID
}

type LoanPrepared struct {
	MemberID   member.MemberID
	BookCopyID bookcopy.BookCopyID
}
