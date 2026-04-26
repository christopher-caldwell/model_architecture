package book

import (
	"time"

	domainbook "github.com/christophercaldwell/model-architecture/go/internal/domain/book"
)

type BookResponse struct {
	ISBN       string    `json:"isbn"`
	DtCreated  time.Time `json:"dt_created"`
	DtModified time.Time `json:"dt_modified"`
	Title      string    `json:"title"`
	AuthorName string    `json:"author_name"`
}

func bookToResponse(b domainbook.Book) BookResponse {
	return BookResponse{
		ISBN:       b.ISBN,
		DtCreated:  b.DtCreated,
		DtModified: b.DtModified,
		Title:      b.Title,
		AuthorName: b.AuthorName,
	}
}

type CreateBookRequest struct {
	ISBN       string `json:"isbn"`
	Title      string `json:"title"`
	AuthorName string `json:"author_name"`
}
