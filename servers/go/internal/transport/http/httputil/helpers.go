package httputil

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/christophercaldwell/model-architecture/go/internal/domain/book"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/loan"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/member"
)

type ErrorBody struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, ErrorBody{Error: msg})
}

func WriteServiceError(w http.ResponseWriter, err error) {
	slog.Error("unhandled request error", "error", err)
	WriteError(w, http.StatusInternalServerError, "Something went wrong")
}

func WriteCommandError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, book.ErrNotFound):
		WriteError(w, http.StatusNotFound, "Book not found")

	case errors.Is(err, bookcopy.ErrNotFound):
		WriteError(w, http.StatusNotFound, "Book copy not found")
	case errors.Is(err, bookcopy.ErrCannotBeBorrowed),
		errors.Is(err, bookcopy.ErrCannotBeSentToMaintenance),
		errors.Is(err, bookcopy.ErrCannotBeReturnedFromMaintenance),
		errors.Is(err, bookcopy.ErrCannotMarkBookLost),
		errors.Is(err, bookcopy.ErrCannotBeReturnedFromLost):
		WriteError(w, http.StatusConflict, err.Error())

	case errors.Is(err, member.ErrNotFound):
		WriteError(w, http.StatusNotFound, "Member not found")
	case errors.Is(err, member.ErrCannotBeSuspended),
		errors.Is(err, member.ErrCannotBeReactivated),
		errors.Is(err, member.ErrCannotBorrowWhileSuspended),
		errors.Is(err, member.ErrLoanLimitReached):
		WriteError(w, http.StatusConflict, err.Error())

	case errors.Is(err, loan.ErrNoActiveLoanForBookCopy),
		errors.Is(err, loan.ErrCannotBeReturned):
		WriteError(w, http.StatusConflict, err.Error())

	default:
		slog.Error("unhandled command error", "error", err)
		WriteError(w, http.StatusInternalServerError, "Something went wrong")
	}
}
