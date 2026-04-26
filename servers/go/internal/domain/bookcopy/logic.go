package bookcopy

func (c *BookCopy) EnsureCanBeBorrowed() error {
	if c.Status != BookCopyStatusActive {
		return ErrCannotBeBorrowed
	}
	return nil
}

func (c *BookCopy) SendToMaintenance() (BookCopyStatus, error) {
	if c.Status != BookCopyStatusActive {
		return "", ErrCannotBeSentToMaintenance
	}
	return BookCopyStatusMaintenance, nil
}

func (c *BookCopy) CompleteMaintenance() (BookCopyStatus, error) {
	if c.Status != BookCopyStatusMaintenance {
		return "", ErrCannotBeReturnedFromMaintenance
	}
	return BookCopyStatusActive, nil
}

func (c *BookCopy) MarkLost() (BookCopyStatus, error) {
	if c.Status == BookCopyStatusLost {
		return "", ErrCannotMarkBookLost
	}
	return BookCopyStatusLost, nil
}

func (c *BookCopy) MarkFound() (BookCopyStatus, error) {
	if c.Status != BookCopyStatusLost {
		return "", ErrCannotBeReturnedFromLost
	}
	return BookCopyStatusActive, nil
}

func (p BookCopyCreationPayload) Prepare() BookCopyPrepared {
	return BookCopyPrepared{
		Barcode: p.Barcode,
		BookID:  p.BookID,
		Status:  BookCopyStatusActive,
	}
}
