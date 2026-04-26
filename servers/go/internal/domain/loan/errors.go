package loan

import "errors"

var (
	ErrNoActiveLoanForBookCopy = errors.New("book copy does not have an active loan")
	ErrCannotBeReturned        = errors.New("loan has already been returned")
)
