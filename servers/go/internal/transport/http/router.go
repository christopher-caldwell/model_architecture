package httptransport

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/christophercaldwell/model-architecture/go/internal/auth"
	"github.com/christophercaldwell/model-architecture/go/internal/bootstrap"
	bookhandler "github.com/christophercaldwell/model-architecture/go/internal/transport/http/book"
	bookcopyhd "github.com/christophercaldwell/model-architecture/go/internal/transport/http/bookcopy"
	healthhd "github.com/christophercaldwell/model-architecture/go/internal/transport/http/health"
	loanhd "github.com/christophercaldwell/model-architecture/go/internal/transport/http/loan"
	memberhd "github.com/christophercaldwell/model-architecture/go/internal/transport/http/member"
)

func NewRouter(deps *bootstrap.ServerDeps) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(corsMiddleware)

	r.Get("/health", healthhd.GetHealthCheck)

	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware(deps.Auth.Verifier))

		r.Get("/books", bookhandler.GetCatalog(deps.Catalog.Queries))
		r.Post("/books", bookhandler.AddBook(deps.Catalog.Commands))
		r.Post("/books/{isbn}/copies", bookhandler.AddBookCopy(deps.Catalog.Commands))

		r.Get("/book-copies/{barcode}", bookcopyhd.GetDetails(deps.Catalog.Queries))
		r.Put("/book-copies/{barcode}/lost", bookcopyhd.MarkLost(deps.Catalog.Commands))
		r.Delete("/book-copies/{barcode}/lost", bookcopyhd.MarkFound(deps.Catalog.Commands))
		r.Put("/book-copies/{barcode}/maintenance", bookcopyhd.SendToMaintenance(deps.Catalog.Commands))
		r.Delete("/book-copies/{barcode}/maintenance", bookcopyhd.CompleteMaintenance(deps.Catalog.Commands))
		r.Post("/book-copies/{barcode}/return", bookcopyhd.ReturnBookCopy(deps.Lending.Commands))
		r.Post("/book-copies/{barcode}/report-loss", bookcopyhd.ReportLostLoaned(deps.Lending.Commands))

		r.Post("/members", memberhd.RegisterMember(deps.Membership.Commands))
		r.Get("/members/{ident}", memberhd.GetMemberDetails(deps.Membership.Queries))
		r.Put("/members/{ident}/suspension", memberhd.SuspendMember(deps.Membership.Commands))
		r.Delete("/members/{ident}/suspension", memberhd.ReactivateMember(deps.Membership.Commands))
		r.Get("/members/{ident}/loans", memberhd.GetMemberLoans(deps.Lending.Queries))

		r.Post("/loans", loanhd.CheckOutBookCopy(deps.Lending.Commands))
		r.Get("/loans/overdue", loanhd.GetOverdueLoans(deps.Lending.Queries))
	})

	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
