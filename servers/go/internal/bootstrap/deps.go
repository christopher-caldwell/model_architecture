package bootstrap

import (
	"context"
	"fmt"

	"github.com/christophercaldwell/model-architecture/go/internal/application/commands"
	"github.com/christophercaldwell/model-architecture/go/internal/application/queries"
	"github.com/christophercaldwell/model-architecture/go/internal/auth"
	pgbook "github.com/christophercaldwell/model-architecture/go/internal/persistence/postgres/book"
	pgbookcopy "github.com/christophercaldwell/model-architecture/go/internal/persistence/postgres/bookcopy"
	pgloan "github.com/christophercaldwell/model-architecture/go/internal/persistence/postgres/loan"
	pgmember "github.com/christophercaldwell/model-architecture/go/internal/persistence/postgres/member"
	"github.com/christophercaldwell/model-architecture/go/internal/persistence/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthDeps struct {
	Verifier auth.Verifier
}

type CatalogDeps struct {
	Commands *commands.CatalogCommands
	Queries  *queries.CatalogQueries
}

type LendingDeps struct {
	Commands *commands.LendingCommands
	Queries  *queries.LendingQueries
}

type MembershipDeps struct {
	Commands *commands.MembershipCommands
	Queries  *queries.MembershipQueries
}

type ServerDeps struct {
	Auth       AuthDeps
	Catalog    CatalogDeps
	Lending    LendingDeps
	Membership MembershipDeps
}

type memberIdentGenerator struct{}

func (g *memberIdentGenerator) Gen() string {
	return generateNanoID(10)
}

func CreateServerDeps(ctx context.Context, cfg *ServerConfig) (*ServerDeps, func(), error) {
	roPool, err := connectPool(ctx, cfg.DatabaseROURL)
	if err != nil {
		return nil, nil, fmt.Errorf("connect to read-only database: %w", err)
	}
	rwPool, err := connectPool(ctx, cfg.DatabaseRWURL)
	if err != nil {
		roPool.Close()
		return nil, nil, fmt.Errorf("connect to read-write database: %w", err)
	}

	cleanup := func() {
		roPool.Close()
		rwPool.Close()
	}

	uowFactory := &postgres.SqlUnitOfWorkFactory{Pool: rwPool}

	bookReadRepo := pgbook.NewReadRepo(roPool)
	bookCopyReadRepo := pgbookcopy.NewReadRepo(roPool)
	loanReadRepo := pgloan.NewReadRepo(roPool)
	memberReadRepo := pgmember.NewReadRepo(roPool)

	catalogCmds := commands.NewCatalogCommands(uowFactory)
	lendingCmds := commands.NewLendingCommands(uowFactory)
	membershipCmds := commands.NewMembershipCommands(uowFactory, &memberIdentGenerator{})

	catalogQs := queries.NewCatalogQueries(bookReadRepo, bookCopyReadRepo)
	lendingQs := queries.NewLendingQueries(loanReadRepo)
	membershipQs := queries.NewMembershipQueries(memberReadRepo)

	return &ServerDeps{
		Auth: AuthDeps{
			Verifier: auth.NewJwtVerifier(cfg.JWTSecret),
		},
		Catalog: CatalogDeps{
			Commands: catalogCmds,
			Queries:  catalogQs,
		},
		Lending: LendingDeps{
			Commands: lendingCmds,
			Queries:  lendingQs,
		},
		Membership: MembershipDeps{
			Commands: membershipCmds,
			Queries:  membershipQs,
		},
	}, cleanup, nil
}

func connectPool(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("parse database URL: %w", err)
	}
	cfg.MaxConns = 5
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	return pool, nil
}
