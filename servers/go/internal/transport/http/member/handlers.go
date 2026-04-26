package member

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/christophercaldwell/model-architecture/go/internal/application/commands"
	"github.com/christophercaldwell/model-architecture/go/internal/application/queries"
	domainm "github.com/christophercaldwell/model-architecture/go/internal/domain/member"
	"github.com/christophercaldwell/model-architecture/go/internal/transport/http/httputil"
	"github.com/go-chi/chi/v5"
)

func RegisterMember(c *commands.MembershipCommands) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body CreateMemberRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		m, err := c.RegisterMember(r.Context(), domainm.MemberCreationPayload{
			FullName:       body.FullName,
			MaxActiveLoans: body.MaxActiveLoans,
		})
		if err != nil {
			httputil.WriteCommandError(w, err)
			return
		}
		httputil.WriteJSON(w, http.StatusCreated, memberToResponse(*m))
	}
}

func GetMemberDetails(q *queries.MembershipQueries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ident := chi.URLParam(r, "ident")
		m, err := q.GetMemberDetails(r.Context(), domainm.MemberIdent(ident))
		if err != nil {
			httputil.WriteServiceError(w, err)
			return
		}
		if m == nil {
			httputil.WriteError(w, http.StatusNotFound, "member not found")
			return
		}
		httputil.WriteJSON(w, http.StatusOK, memberToResponse(*m))
	}
}

func SuspendMember(c *commands.MembershipCommands) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ident := chi.URLParam(r, "ident")
		m, err := c.SuspendMember(r.Context(), commands.MemberIdentInput{MemberIdent: ident})
		if err != nil {
			httputil.WriteCommandError(w, err)
			return
		}
		httputil.WriteJSON(w, http.StatusOK, memberToResponse(*m))
	}
}

func ReactivateMember(c *commands.MembershipCommands) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ident := chi.URLParam(r, "ident")
		m, err := c.ReactivateMember(r.Context(), commands.MemberIdentInput{MemberIdent: ident})
		if err != nil {
			httputil.WriteCommandError(w, err)
			return
		}
		httputil.WriteJSON(w, http.StatusOK, memberToResponse(*m))
	}
}

func GetMemberLoans(q *queries.LendingQueries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ident := chi.URLParam(r, "ident")
		loans, err := q.GetMemberLoans(r.Context(), domainm.MemberIdent(ident))
		if err != nil {
			httputil.WriteServiceError(w, err)
			return
		}
		type loanResp struct {
			Ident      string     `json:"ident"`
			DtCreated  time.Time  `json:"dt_created"`
			DtModified time.Time  `json:"dt_modified"`
			DtDue      *time.Time `json:"dt_due"`
			DtReturned *time.Time `json:"dt_returned"`
		}
		resp := make([]loanResp, len(loans))
		for i, l := range loans {
			resp[i] = loanResp{
				Ident:      string(l.Ident),
				DtCreated:  l.DtCreated,
				DtModified: l.DtModified,
				DtDue:      l.DtDue,
				DtReturned: l.DtReturned,
			}
		}
		httputil.WriteJSON(w, http.StatusOK, resp)
	}
}
