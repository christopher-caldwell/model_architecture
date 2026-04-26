package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/christophercaldwell/model-architecture/go/internal/auth"
	"github.com/christophercaldwell/model-architecture/go/internal/bootstrap"
	gqltransport "github.com/christophercaldwell/model-architecture/go/internal/transport/graphql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx := context.Background()

	cfg, err := bootstrap.LoadConfig()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}

	deps, cleanup, err := bootstrap.CreateServerDeps(ctx, cfg)
	if err != nil {
		slog.Error("create server deps", "error", err)
		os.Exit(1)
	}
	defer cleanup()

	srv := gqltransport.NewSchema(deps)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/graphql", playground.Handler("GraphQL Playground", "/graphql"))
	r.With(auth.Middleware(deps.Auth.Verifier)).Post("/graphql", srv.ServeHTTP)

	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	slog.Info("starting GraphQL server", "addr", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
