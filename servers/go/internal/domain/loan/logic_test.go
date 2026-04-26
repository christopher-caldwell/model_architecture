package loan_test

import (
	"errors"
	"testing"
	"time"

	"github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/loan"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/member"
)

func activeLoan() *loan.Loan {
	return &loan.Loan{
		ID:         1,
		Ident:      "LOAN-0001",
		DtCreated:  time.Now(),
		DtModified: time.Now(),
		BookCopyID: bookcopy.BookCopyID(1),
		MemberID:   member.MemberID(1),
		DtDue:      nil,
		DtReturned: nil,
	}
}

func returnedLoan() *loan.Loan {
	l := activeLoan()
	now := time.Now()
	l.DtReturned = &now
	return l
}

func TestEnsureCanBeReturned(t *testing.T) {
	tests := []struct {
		name    string
		l       *loan.Loan
		wantErr error
	}{
		{"active loan can be returned", activeLoan(), nil},
		{"returned loan cannot be returned again", returnedLoan(), loan.ErrCannotBeReturned},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.l.EnsureCanBeReturned()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoanCreationPayloadPrepare(t *testing.T) {
	payload := loan.LoanCreationPayload{
		MemberID:   member.MemberID(42),
		BookCopyID: bookcopy.BookCopyID(7),
	}
	prepared := payload.Prepare()
	if prepared.MemberID != 42 {
		t.Errorf("got member_id %v, want 42", prepared.MemberID)
	}
	if prepared.BookCopyID != 7 {
		t.Errorf("got book_copy_id %v, want 7", prepared.BookCopyID)
	}
}
