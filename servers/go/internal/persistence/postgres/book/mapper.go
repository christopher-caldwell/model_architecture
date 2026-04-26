package book

import (
	domain "github.com/christophercaldwell/model-architecture/go/internal/domain/book"
)

func rowToDomain(r bookRow) domain.Book {
	return domain.Book{
		ID:         domain.BookID(r.BookID),
		ISBN:       r.ISBN,
		DtCreated:  r.DtCreated,
		DtModified: r.DtModified,
		Title:      r.Title,
		AuthorName: r.AuthorName,
	}
}
