package member

import (
	"context"
	"fmt"

	domain "github.com/christophercaldwell/model-architecture/go/internal/domain/member"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReadRepo struct {
	pool *pgxpool.Pool
}

func NewReadRepo(pool *pgxpool.Pool) *ReadRepo {
	return &ReadRepo{pool: pool}
}

const selectMemberCols = `
	m.member_id, m.member_ident, m.dt_created, m.dt_modified, st.att_pub_ident, m.full_name, m.max_active_loans`

const memberJoin = `
	FROM library.member m
	JOIN library.struct_type st ON m.status_id = st.struct_type_id
	WHERE st.group_name = 'member_status' AND st.att_pub_ident IN ('active', 'suspended')`

func scanMemberRow(row pgx.Row) (*domain.Member, error) {
	var r memberRow
	if err := row.Scan(&r.MemberID, &r.MemberIdent, &r.DtCreated, &r.DtModified, &r.Status, &r.FullName, &r.MaxActiveLoans); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan member row: %w", err)
	}
	return rowToDomain(r)
}

func (r *ReadRepo) GetByID(ctx context.Context, id domain.MemberID) (*domain.Member, error) {
	q := `SELECT ` + selectMemberCols + memberJoin + ` AND m.member_id = $1`
	m, err := scanMemberRow(r.pool.QueryRow(ctx, q, int32(id)))
	if err != nil {
		return nil, fmt.Errorf("fetch member by id: %w", err)
	}
	return m, nil
}

func (r *ReadRepo) GetByIdent(ctx context.Context, ident domain.MemberIdent) (*domain.Member, error) {
	q := `SELECT ` + selectMemberCols + memberJoin + ` AND m.member_ident = $1`
	m, err := scanMemberRow(r.pool.QueryRow(ctx, q, string(ident)))
	if err != nil {
		return nil, fmt.Errorf("fetch member by ident: %w", err)
	}
	return m, nil
}
