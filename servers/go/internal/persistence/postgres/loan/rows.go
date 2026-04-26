package loan

import "time"

type loanRow struct {
	LoanID     int32
	LoanIdent  string
	DtCreated  time.Time
	DtModified time.Time
	BookCopyID int32
	MemberID   int32
	DtDue      *time.Time
	DtReturned *time.Time
}

type countRow struct {
	Count int64
}
