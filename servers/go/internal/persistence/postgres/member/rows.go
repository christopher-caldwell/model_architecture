package member

import "time"

type memberRow struct {
	MemberID       int32
	MemberIdent    string
	DtCreated      time.Time
	DtModified     time.Time
	Status         string
	FullName       string
	MaxActiveLoans int16
}
