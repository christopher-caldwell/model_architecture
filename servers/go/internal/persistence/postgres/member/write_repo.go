package member

import (
	"context"
	"fmt"
	"time"

	domain "github.com/christophercaldwell/model-architecture/go/internal/domain/member"
	"github.com/jackc/pgx/v5"
)

type WriteRepo struct {
	tx pgx.Tx
}

func NewWriteRepo(tx pgx.Tx) *WriteRepo {
	return &WriteRepo{tx: tx}
}

func (r *WriteRepo) Create(ctx context.Context, prepared domain.MemberPrepared) (*domain.Member, error) {
	const q = `
		INSERT INTO library.member (member_ident, status_id, full_name, max_active_loans)
		VALUES (
			$1,
			(SELECT st.struct_type_id FROM library.struct_type st
			 WHERE st.group_name = 'member_status' AND st.att_pub_ident = $2),
			$3,
			$4
		)
		RETURNING member_id`

	var memberID int32
	err := r.tx.QueryRow(ctx, q,
		string(prepared.Ident),
		string(prepared.Status),
		prepared.FullName,
		prepared.MaxActiveLoans,
	).Scan(&memberID)
	if err != nil {
		return nil, fmt.Errorf("create member: %w", err)
	}
	now := time.Now()
	m := domain.Member{
		ID:             domain.MemberID(memberID),
		Ident:          prepared.Ident,
		DtCreated:      now,
		DtModified:     now,
		Status:         prepared.Status,
		FullName:       prepared.FullName,
		MaxActiveLoans: prepared.MaxActiveLoans,
	}
	return &m, nil
}

func (r *WriteRepo) GetByIdentForUpdate(ctx context.Context, ident domain.MemberIdent) (*domain.Member, error) {
	const q = `
		SELECT m.member_id, m.member_ident, m.dt_created, m.dt_modified, st.att_pub_ident, m.full_name, m.max_active_loans
		FROM library.member m
		JOIN library.struct_type st ON m.status_id = st.struct_type_id
		WHERE m.member_ident = $1
		  AND st.group_name = 'member_status'
		  AND st.att_pub_ident IN ('active', 'suspended')
		FOR UPDATE OF m`

	var row memberRow
	err := r.tx.QueryRow(ctx, q, string(ident)).Scan(
		&row.MemberID, &row.MemberIdent, &row.DtCreated, &row.DtModified,
		&row.Status, &row.FullName, &row.MaxActiveLoans,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("fetch member by ident for update: %w", err)
	}
	return rowToDomain(row)
}

func (r *WriteRepo) UpdateStatus(ctx context.Context, id domain.MemberID, status domain.MemberStatus) error {
	const q = `
		UPDATE library.member m
		SET status_id = (
			SELECT st.struct_type_id FROM library.struct_type st
			WHERE st.group_name = 'member_status' AND st.att_pub_ident = $2
		)
		WHERE m.member_id = $1`

	_, err := r.tx.Exec(ctx, q, int32(id), string(status))
	if err != nil {
		return fmt.Errorf("update member status: %w", err)
	}
	return nil
}
