package bookcopy_test

import (
	"errors"
	"testing"
	"time"

	"github.com/christophercaldwell/model-architecture/go/internal/domain/book"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
)

func copyWithStatus(status bookcopy.BookCopyStatus) *bookcopy.BookCopy {
	return &bookcopy.BookCopy{
		ID:         1,
		Barcode:    "BC-001",
		DtCreated:  time.Now(),
		DtModified: time.Now(),
		BookID:     book.BookID(1),
		Status:     status,
	}
}

func TestEnsureCanBeBorrowed(t *testing.T) {
	tests := []struct {
		name    string
		status  bookcopy.BookCopyStatus
		wantErr error
	}{
		{"active can borrow", bookcopy.BookCopyStatusActive, nil},
		{"maintenance cannot borrow", bookcopy.BookCopyStatusMaintenance, bookcopy.ErrCannotBeBorrowed},
		{"lost cannot borrow", bookcopy.BookCopyStatusLost, bookcopy.ErrCannotBeBorrowed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := copyWithStatus(tt.status).EnsureCanBeBorrowed()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendToMaintenance(t *testing.T) {
	tests := []struct {
		name       string
		status     bookcopy.BookCopyStatus
		wantStatus bookcopy.BookCopyStatus
		wantErr    error
	}{
		{"active -> maintenance", bookcopy.BookCopyStatusActive, bookcopy.BookCopyStatusMaintenance, nil},
		{"maintenance -> error", bookcopy.BookCopyStatusMaintenance, "", bookcopy.ErrCannotBeSentToMaintenance},
		{"lost -> error", bookcopy.BookCopyStatusLost, "", bookcopy.ErrCannotBeSentToMaintenance},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := copyWithStatus(tt.status).SendToMaintenance()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got err %v, want %v", err, tt.wantErr)
			}
			if err == nil && got != tt.wantStatus {
				t.Errorf("got status %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func TestCompleteMaintenance(t *testing.T) {
	tests := []struct {
		name       string
		status     bookcopy.BookCopyStatus
		wantStatus bookcopy.BookCopyStatus
		wantErr    error
	}{
		{"maintenance -> active", bookcopy.BookCopyStatusMaintenance, bookcopy.BookCopyStatusActive, nil},
		{"active -> error", bookcopy.BookCopyStatusActive, "", bookcopy.ErrCannotBeReturnedFromMaintenance},
		{"lost -> error", bookcopy.BookCopyStatusLost, "", bookcopy.ErrCannotBeReturnedFromMaintenance},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := copyWithStatus(tt.status).CompleteMaintenance()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got err %v, want %v", err, tt.wantErr)
			}
			if err == nil && got != tt.wantStatus {
				t.Errorf("got status %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func TestMarkLost(t *testing.T) {
	tests := []struct {
		name       string
		status     bookcopy.BookCopyStatus
		wantStatus bookcopy.BookCopyStatus
		wantErr    error
	}{
		{"active -> lost", bookcopy.BookCopyStatusActive, bookcopy.BookCopyStatusLost, nil},
		{"maintenance -> lost", bookcopy.BookCopyStatusMaintenance, bookcopy.BookCopyStatusLost, nil},
		{"lost -> error", bookcopy.BookCopyStatusLost, "", bookcopy.ErrCannotMarkBookLost},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := copyWithStatus(tt.status).MarkLost()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got err %v, want %v", err, tt.wantErr)
			}
			if err == nil && got != tt.wantStatus {
				t.Errorf("got status %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func TestMarkFound(t *testing.T) {
	tests := []struct {
		name       string
		status     bookcopy.BookCopyStatus
		wantStatus bookcopy.BookCopyStatus
		wantErr    error
	}{
		{"lost -> active", bookcopy.BookCopyStatusLost, bookcopy.BookCopyStatusActive, nil},
		{"active -> error", bookcopy.BookCopyStatusActive, "", bookcopy.ErrCannotBeReturnedFromLost},
		{"maintenance -> error", bookcopy.BookCopyStatusMaintenance, "", bookcopy.ErrCannotBeReturnedFromLost},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := copyWithStatus(tt.status).MarkFound()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got err %v, want %v", err, tt.wantErr)
			}
			if err == nil && got != tt.wantStatus {
				t.Errorf("got status %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func TestBookCopyCreationPayloadPrepare(t *testing.T) {
	payload := bookcopy.BookCopyCreationPayload{
		Barcode: "BC-002",
		BookID:  book.BookID(2),
	}
	prepared := payload.Prepare()
	if prepared.Status != bookcopy.BookCopyStatusActive {
		t.Errorf("got status %v, want active", prepared.Status)
	}
	if prepared.Barcode != "BC-002" {
		t.Errorf("got barcode %v, want BC-002", prepared.Barcode)
	}
}
