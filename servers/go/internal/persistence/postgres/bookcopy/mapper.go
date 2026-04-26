package bookcopy

import (
	"fmt"

	domainbook "github.com/christophercaldwell/model-architecture/go/internal/domain/book"
	domain "github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
)

func rowToDomain(r bookCopyRow) (*domain.BookCopy, error) {
	status, err := parseStatus(r.Status)
	if err != nil {
		return nil, err
	}
	return &domain.BookCopy{
		ID:         domain.BookCopyID(r.BookCopyID),
		Barcode:    r.Barcode,
		DtCreated:  r.DtCreated,
		DtModified: r.DtModified,
		BookID:     domainbook.BookID(r.BookID),
		Status:     status,
	}, nil
}

func parseStatus(s string) (domain.BookCopyStatus, error) {
	switch domain.BookCopyStatus(s) {
	case domain.BookCopyStatusActive, domain.BookCopyStatusMaintenance, domain.BookCopyStatusLost:
		return domain.BookCopyStatus(s), nil
	default:
		return "", fmt.Errorf("unknown book copy status in DB: %q", s)
	}
}
