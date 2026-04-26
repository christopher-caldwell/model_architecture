package book

func (p BookCreationPayload) Prepare() BookPrepared {
	return BookPrepared{
		ISBN:       p.ISBN,
		Title:      p.Title,
		AuthorName: p.AuthorName,
	}
}
