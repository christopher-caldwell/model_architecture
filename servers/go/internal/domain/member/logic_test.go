package member_test

import (
	"errors"
	"testing"
	"time"

	"github.com/christophercaldwell/model-architecture/go/internal/domain/member"
)

func activeMember() *member.Member {
	return &member.Member{
		ID:             1,
		Ident:          "TEST-0001",
		DtCreated:      time.Now(),
		DtModified:     time.Now(),
		Status:         member.MemberStatusActive,
		FullName:       "Alice Smith",
		MaxActiveLoans: 3,
	}
}

func suspendedMember() *member.Member {
	m := activeMember()
	m.Status = member.MemberStatusSuspended
	return m
}

func TestSuspend(t *testing.T) {
	tests := []struct {
		name       string
		m          *member.Member
		wantStatus member.MemberStatus
		wantErr    error
	}{
		{"active -> suspended", activeMember(), member.MemberStatusSuspended, nil},
		{"suspended -> error", suspendedMember(), "", member.ErrCannotBeSuspended},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Suspend()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got err %v, want %v", err, tt.wantErr)
			}
			if err == nil && got != tt.wantStatus {
				t.Errorf("got status %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func TestReactivate(t *testing.T) {
	tests := []struct {
		name       string
		m          *member.Member
		wantStatus member.MemberStatus
		wantErr    error
	}{
		{"suspended -> active", suspendedMember(), member.MemberStatusActive, nil},
		{"active -> error", activeMember(), "", member.ErrCannotBeReactivated},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Reactivate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got err %v, want %v", err, tt.wantErr)
			}
			if err == nil && got != tt.wantStatus {
				t.Errorf("got status %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func TestEnsureCanBorrow(t *testing.T) {
	tests := []struct {
		name    string
		m       *member.Member
		wantErr error
	}{
		{"active can borrow", activeMember(), nil},
		{"suspended cannot borrow", suspendedMember(), member.ErrCannotBorrowWhileSuspended},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.m.EnsureCanBorrow()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnsureWithinLoanLimit(t *testing.T) {
	m := activeMember() // max_active_loans = 3
	tests := []struct {
		name        string
		activeLoans int16
		wantErr     error
	}{
		{"below limit (2 of 3)", 2, nil},
		{"at limit (3 of 3)", 3, member.ErrLoanLimitReached},
		{"over limit (4 of 3)", 4, member.ErrLoanLimitReached},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := m.EnsureWithinLoanLimit(tt.activeLoans)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemberCreationPayloadPrepare(t *testing.T) {
	payload := member.MemberCreationPayload{
		FullName:       "Bob Jones",
		MaxActiveLoans: 5,
	}
	prepared := payload.Prepare(member.MemberIdent("NEW-IDENT"))
	if prepared.Status != member.MemberStatusActive {
		t.Errorf("got status %v, want active", prepared.Status)
	}
	if prepared.MaxActiveLoans != 5 {
		t.Errorf("got max_active_loans %v, want 5", prepared.MaxActiveLoans)
	}
	if prepared.Ident != "NEW-IDENT" {
		t.Errorf("got ident %v, want NEW-IDENT", prepared.Ident)
	}
}
