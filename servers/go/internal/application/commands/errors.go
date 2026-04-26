package commands

import (
	"errors"
	"fmt"

	"github.com/christophercaldwell/model-architecture/go/internal/domain/book"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/loan"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/member"
)

// CommandError wraps domain and infrastructure errors from command execution.
type CommandError struct {
	err error
}

func (e *CommandError) Error() string { return e.err.Error() }
func (e *CommandError) Unwrap() error { return e.err }

func wrapCommand(err error) error {
	if err == nil {
		return nil
	}
	return &CommandError{err: err}
}

func wrapCommandf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return &CommandError{err: fmt.Errorf(format+": %w", append(args, err)...)}
}

// IsDomainError returns true when err is a known domain-level sentinel.
func IsDomainError(err error) bool {
	return errors.Is(err, book.ErrNotFound) ||
		errors.Is(err, bookcopy.ErrNotFound) ||
		errors.Is(err, bookcopy.ErrCannotBeBorrowed) ||
		errors.Is(err, bookcopy.ErrCannotBeSentToMaintenance) ||
		errors.Is(err, bookcopy.ErrCannotBeReturnedFromMaintenance) ||
		errors.Is(err, bookcopy.ErrCannotMarkBookLost) ||
		errors.Is(err, bookcopy.ErrCannotBeReturnedFromLost) ||
		errors.Is(err, member.ErrNotFound) ||
		errors.Is(err, member.ErrCannotBeSuspended) ||
		errors.Is(err, member.ErrCannotBeReactivated) ||
		errors.Is(err, member.ErrCannotBorrowWhileSuspended) ||
		errors.Is(err, member.ErrLoanLimitReached) ||
		errors.Is(err, loan.ErrNoActiveLoanForBookCopy) ||
		errors.Is(err, loan.ErrCannotBeReturned)
}
