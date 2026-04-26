package member

import "time"

type MemberID int32

type MemberIdent string

type MemberStatus string

const (
	MemberStatusActive    MemberStatus = "active"
	MemberStatusSuspended MemberStatus = "suspended"
)

type Member struct {
	ID             MemberID
	Ident          MemberIdent
	DtCreated      time.Time
	DtModified     time.Time
	Status         MemberStatus
	FullName       string
	MaxActiveLoans int16
}

type MemberCreationPayload struct {
	FullName       string
	MaxActiveLoans int16
}

type MemberPrepared struct {
	Ident          MemberIdent
	FullName       string
	MaxActiveLoans int16
	Status         MemberStatus
}
