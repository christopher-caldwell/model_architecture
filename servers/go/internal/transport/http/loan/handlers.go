package loan

import (
	"encoding/json"
	"net/http"

	"github.com/christophercaldwell/model-architecture/go/internal/application/commands"
	"github.com/christophercaldwell/model-architecture/go/internal/application/queries"
	"github.com/christophercaldwell/model-architecture/go/internal/transport/http/httputil"
)

func CheckOutBookCopy(c *commands.LendingCommands) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body CreateLoanRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		l, err := c.CheckOutBookCopy(r.Context(), commands.CheckOutBookCopyInput{
			MemberIdent:     body.MemberIdent,
			BookCopyBarcode: body.BookCopyBarcode,
		})
		if err != nil {
			httputil.WriteCommandError(w, err)
			return
		}
		httputil.WriteJSON(w, http.StatusCreated, LoanToResponse(*l))
	}
}

func GetOverdueLoans(q *queries.LendingQueries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loans, err := q.GetOverdueLoans(r.Context())
		if err != nil {
			httputil.WriteServiceError(w, err)
			return
		}
		resp := make([]LoanResponse, len(loans))
		for i, l := range loans {
			resp[i] = LoanToResponse(l)
		}
		httputil.WriteJSON(w, http.StatusOK, resp)
	}
}
