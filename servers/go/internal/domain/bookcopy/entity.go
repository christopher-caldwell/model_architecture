package bookcopy

import (
	"time"

	"github.com/christophercaldwell/model-architecture/go/internal/domain/book"
)

type BookCopyID int32

type BookCopyStatus string

const (
	BookCopyStatusActive      BookCopyStatus = "active"
	BookCopyStatusMaintenance BookCopyStatus = "maintenance"
	BookCopyStatusLost        BookCopyStatus = "lost"
)

type BookCopy struct {
	ID         BookCopyID
	Barcode    string
	DtCreated  time.Time
	DtModified time.Time
	BookID     book.BookID
	Status     BookCopyStatus
}

type BookCopyCreationPayload struct {
	Barcode string
	BookID  book.BookID
}

type BookCopyPrepared struct {
	Barcode string
	BookID  book.BookID
	Status  BookCopyStatus
}
