package member

import "errors"

var (
	ErrNotFound                   = errors.New("member not found")
	ErrCannotBeSuspended          = errors.New("member is already suspended")
	ErrCannotBeReactivated        = errors.New("member is not currently suspended")
	ErrCannotBorrowWhileSuspended = errors.New("member is suspended and cannot borrow new books")
	ErrLoanLimitReached           = errors.New("member has reached the maximum number of active loans")
)
