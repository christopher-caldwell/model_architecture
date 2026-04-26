package bookcopy

import (
	"net/http"

	"github.com/christophercaldwell/model-architecture/go/internal/application/commands"
	"github.com/christophercaldwell/model-architecture/go/internal/application/queries"
	"github.com/christophercaldwell/model-architecture/go/internal/transport/http/httputil"
	loandto "github.com/christophercaldwell/model-architecture/go/internal/transport/http/loan"
	"github.com/go-chi/chi/v5"
)

func GetDetails(q *queries.CatalogQueries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		barcode := chi.URLParam(r, "barcode")
		bc, err := q.GetBookCopyDetails(r.Context(), barcode)
		if err != nil {
			httputil.WriteServiceError(w, err)
			return
		}
		if bc == nil {
			httputil.WriteError(w, http.StatusNotFound, "Book copy not found")
			return
		}
		httputil.WriteJSON(w, http.StatusOK, BookCopyToResponse(*bc))
	}
}

func MarkLost(c *commands.CatalogCommands) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		barcode := chi.URLParam(r, "barcode")
		bc, err := c.MarkBookCopyLost(r.Context(), barcode)
		if err != nil {
			httputil.WriteCommandError(w, err)
			return
		}
		httputil.WriteJSON(w, http.StatusOK, BookCopyToResponse(*bc))
	}
}

func MarkFound(c *commands.CatalogCommands) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		barcode := chi.URLParam(r, "barcode")
		bc, err := c.MarkBookCopyFound(r.Context(), barcode)
		if err != nil {
			httputil.WriteCommandError(w, err)
			return
		}
		httputil.WriteJSON(w, http.StatusOK, BookCopyToResponse(*bc))
	}
}

func SendToMaintenance(c *commands.CatalogCommands) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		barcode := chi.URLParam(r, "barcode")
		bc, err := c.SendBookCopyToMaintenance(r.Context(), barcode)
		if err != nil {
			httputil.WriteCommandError(w, err)
			return
		}
		httputil.WriteJSON(w, http.StatusOK, BookCopyToResponse(*bc))
	}
}

func CompleteMaintenance(c *commands.CatalogCommands) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		barcode := chi.URLParam(r, "barcode")
		bc, err := c.CompleteBookCopyMaintenance(r.Context(), barcode)
		if err != nil {
			httputil.WriteCommandError(w, err)
			return
		}
		httputil.WriteJSON(w, http.StatusOK, BookCopyToResponse(*bc))
	}
}

func ReturnBookCopy(c *commands.LendingCommands) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		barcode := chi.URLParam(r, "barcode")
		l, err := c.ReturnBookCopy(r.Context(), barcode)
		if err != nil {
			httputil.WriteCommandError(w, err)
			return
		}
		httputil.WriteJSON(w, http.StatusOK, loandto.LoanToResponse(*l))
	}
}

func ReportLostLoaned(c *commands.LendingCommands) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		barcode := chi.URLParam(r, "barcode")
		bc, err := c.ReportLostLoanedBookCopy(r.Context(), barcode)
		if err != nil {
			httputil.WriteCommandError(w, err)
			return
		}
		httputil.WriteJSON(w, http.StatusOK, BookCopyToResponse(*bc))
	}
}
