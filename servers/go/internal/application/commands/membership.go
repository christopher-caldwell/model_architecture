package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/christophercaldwell/model-architecture/go/internal/application"
	"github.com/christophercaldwell/model-architecture/go/internal/application/ports"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/member"
)

type MemberIdentInput struct {
	MemberIdent string
}

type MembershipCommands struct {
	factory        application.UnitOfWorkFactory
	identGenerator ports.IdentGeneratorPort
}

func NewMembershipCommands(factory application.UnitOfWorkFactory, identGenerator ports.IdentGeneratorPort) *MembershipCommands {
	return &MembershipCommands{
		factory:        factory,
		identGenerator: identGenerator,
	}
}

func (c *MembershipCommands) getMemberByIdent(ctx context.Context, uow application.UnitOfWork, ident string) (*member.Member, error) {
	m, err := uow.Members().GetByIdentForUpdate(ctx, member.MemberIdent(ident))
	if err != nil {
		return nil, fmt.Errorf("load member for write: %w", err)
	}
	if m == nil {
		return nil, member.ErrNotFound
	}
	return m, nil
}

func (c *MembershipCommands) RegisterMember(ctx context.Context, payload member.MemberCreationPayload) (*member.Member, error) {
	ident := member.MemberIdent(c.identGenerator.Gen())
	prepared := payload.Prepare(ident)

	uow, err := c.factory.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("build unit of work: %w", err)
	}
	defer uow.Rollback(ctx) //nolint:errcheck

	result, err := uow.Members().Create(ctx, prepared)
	if err != nil {
		return nil, fmt.Errorf("register member: %w", err)
	}
	if err := uow.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}

func (c *MembershipCommands) SuspendMember(ctx context.Context, input MemberIdentInput) (*member.Member, error) {
	uow, err := c.factory.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("build unit of work: %w", err)
	}
	defer uow.Rollback(ctx) //nolint:errcheck

	m, err := c.getMemberByIdent(ctx, uow, input.MemberIdent)
	if err != nil {
		return nil, err
	}
	newStatus, err := m.Suspend()
	if err != nil {
		return nil, err
	}
	if err := uow.Members().UpdateStatus(ctx, m.ID, newStatus); err != nil {
		return nil, fmt.Errorf("suspend member: %w", err)
	}
	if err := uow.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	updated := *m
	updated.Status = newStatus
	updated.DtModified = time.Now()
	return &updated, nil
}

func (c *MembershipCommands) ReactivateMember(ctx context.Context, input MemberIdentInput) (*member.Member, error) {
	uow, err := c.factory.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("build unit of work: %w", err)
	}
	defer uow.Rollback(ctx) //nolint:errcheck

	m, err := c.getMemberByIdent(ctx, uow, input.MemberIdent)
	if err != nil {
		return nil, err
	}
	newStatus, err := m.Reactivate()
	if err != nil {
		return nil, err
	}
	if err := uow.Members().UpdateStatus(ctx, m.ID, newStatus); err != nil {
		return nil, fmt.Errorf("reactivate member: %w", err)
	}
	if err := uow.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	updated := *m
	updated.Status = newStatus
	updated.DtModified = time.Now()
	return &updated, nil
}
