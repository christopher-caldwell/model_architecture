package bookcopy

import "errors"

var (
	ErrNotFound                        = errors.New("book copy not found")
	ErrCannotBeBorrowed                = errors.New("book cannot currently be borrowed")
	ErrCannotBeSentToMaintenance       = errors.New("book is not active and cannot be sent to maintenance")
	ErrCannotBeReturnedFromMaintenance = errors.New("book is not currently in maintenance, and therefore cannot be returned")
	ErrCannotMarkBookLost              = errors.New("book is already marked lost")
	ErrCannotBeReturnedFromLost        = errors.New("book is not currently lost, and cannot be returned from lost")
)
