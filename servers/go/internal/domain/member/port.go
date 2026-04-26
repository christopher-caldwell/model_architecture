package member

import "context"

type ReadRepository interface {
	GetByID(ctx context.Context, id MemberID) (*Member, error)
	GetByIdent(ctx context.Context, ident MemberIdent) (*Member, error)
}

type WriteRepository interface {
	Create(ctx context.Context, prepared MemberPrepared) (*Member, error)
	GetByIdentForUpdate(ctx context.Context, ident MemberIdent) (*Member, error)
	UpdateStatus(ctx context.Context, id MemberID, status MemberStatus) error
}
