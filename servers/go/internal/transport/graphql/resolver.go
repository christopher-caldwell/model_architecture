package graphql

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/christophercaldwell/model-architecture/go/internal/bootstrap"
	"github.com/christophercaldwell/model-architecture/go/internal/transport/graphql/generated"
)

type Resolver struct {
	deps *bootstrap.ServerDeps
}

func NewSchema(deps *bootstrap.ServerDeps) *handler.Server {
	return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{deps: deps},
	}))
}
