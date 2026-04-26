package book

import (
	"encoding/json"
	"net/http"

	"github.com/christophercaldwell/model-architecture/go/internal/application/commands"
	"github.com/christophercaldwell/model-architecture/go/internal/application/queries"
	domainbook "github.com/christophercaldwell/model-architecture/go/internal/domain/book"
	bookcopydto "github.com/christophercaldwell/model-architecture/go/internal/transport/http/bookcopy"
	"github.com/christophercaldwell/model-architecture/go/internal/transport/http/httputil"
	"github.com/go-chi/chi/v5"
)

func GetCatalog(q *queries.CatalogQueries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := q.GetBookCatalog(r.Context())
		if err != nil {
			httputil.WriteServiceError(w, err)
			return
		}
		resp := make([]BookResponse, len(books))
		for i, b := range books {
			resp[i] = bookToResponse(b)
		}
		httputil.WriteJSON(w, http.StatusOK, resp)
	}
}

func AddBook(c *commands.CatalogCommands) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body CreateBookRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		b, err := c.AddBook(r.Context(), domainbook.BookCreationPayload{
			ISBN:       body.ISBN,
			Title:      body.Title,
			AuthorName: body.AuthorName,
		})
		if err != nil {
			httputil.WriteCommandError(w, err)
			return
		}
		httputil.WriteJSON(w, http.StatusCreated, bookToResponse(*b))
	}
}

func AddBookCopy(c *commands.CatalogCommands) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isbn := chi.URLParam(r, "isbn")
		var body struct {
			Barcode string `json:"barcode"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		bc, err := c.AddBookCopy(r.Context(), commands.AddBookCopyInput{
			ISBN:    isbn,
			Barcode: body.Barcode,
		})
		if err != nil {
			httputil.WriteCommandError(w, err)
			return
		}
		httputil.WriteJSON(w, http.StatusCreated, bookcopydto.BookCopyToResponse(*bc))
	}
}
