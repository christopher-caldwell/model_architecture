package graphql

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

import (
	"context"
	"errors"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/christophercaldwell/model-architecture/go/internal/bootstrap"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/book"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/loan"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/member"
	"github.com/christophercaldwell/model-architecture/go/internal/transport/graphql/generated"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Resolver struct {
	deps *bootstrap.ServerDeps
}

func NewSchema(deps *bootstrap.ServerDeps) *handler.Server {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{deps: deps},
	}))
	srv.SetErrorPresenter(presentError)
	return srv
}

func presentError(ctx context.Context, err error) *gqlerror.Error {
	switch {
	case errors.Is(err, book.ErrNotFound):
		return gqlError(ctx, "Book not found", "NOT_FOUND")
	case errors.Is(err, bookcopy.ErrNotFound):
		return gqlError(ctx, "Book copy not found", "NOT_FOUND")
	case errors.Is(err, member.ErrNotFound):
		return gqlError(ctx, "Member not found", "NOT_FOUND")
	case errors.Is(err, bookcopy.ErrCannotBeBorrowed),
		errors.Is(err, bookcopy.ErrCannotBeSentToMaintenance),
		errors.Is(err, bookcopy.ErrCannotBeReturnedFromMaintenance),
		errors.Is(err, bookcopy.ErrCannotMarkBookLost),
		errors.Is(err, bookcopy.ErrCannotBeReturnedFromLost),
		errors.Is(err, member.ErrCannotBeSuspended),
		errors.Is(err, member.ErrCannotBeReactivated),
		errors.Is(err, member.ErrCannotBorrowWhileSuspended),
		errors.Is(err, member.ErrLoanLimitReached),
		errors.Is(err, loan.ErrNoActiveLoanForBookCopy),
		errors.Is(err, loan.ErrCannotBeReturned):
		return gqlError(ctx, err.Error(), "CONFLICT")
	default:
		if gqlErr, ok := err.(*gqlerror.Error); ok {
			return gqlErr
		}
		slog.Error("unhandled GraphQL error", "error", err)
		return gqlError(ctx, "Something went wrong", "INTERNAL_SERVER_ERROR")
	}
}

func gqlError(ctx context.Context, message string, code string) *gqlerror.Error {
	gqlErr := graphql.DefaultErrorPresenter(ctx, errors.New(message))
	gqlErr.Extensions = map[string]any{"code": code}
	return gqlErr
}
