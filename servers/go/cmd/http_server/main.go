package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/christophercaldwell/model-architecture/go/internal/bootstrap"
	httptransport "github.com/christophercaldwell/model-architecture/go/internal/transport/http"
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

	router := httptransport.NewRouter(deps)
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	slog.Info("starting HTTP server", "addr", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
