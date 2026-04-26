package bookcopy

import (
	"time"

	domainbc "github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
)

type BookCopyResponse struct {
	Barcode    string    `json:"barcode"`
	DtCreated  time.Time `json:"dt_created"`
	DtModified time.Time `json:"dt_modified"`
	Status     string    `json:"status"`
}

func BookCopyToResponse(bc domainbc.BookCopy) BookCopyResponse {
	return BookCopyResponse{
		Barcode:    bc.Barcode,
		DtCreated:  bc.DtCreated,
		DtModified: bc.DtModified,
		Status:     string(bc.Status),
	}
}
