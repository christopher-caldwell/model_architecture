package loan

import (
	"context"
	"fmt"
	"time"

	domainbc "github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
	domain "github.com/christophercaldwell/model-architecture/go/internal/domain/loan"
	domainm "github.com/christophercaldwell/model-architecture/go/internal/domain/member"
	"github.com/jackc/pgx/v5"
)

type WriteRepo struct {
	tx pgx.Tx
}

func NewWriteRepo(tx pgx.Tx) *WriteRepo {
	return &WriteRepo{tx: tx}
}

func (r *WriteRepo) Create(ctx context.Context, prepared domain.LoanPrepared) (*domain.Loan, error) {
	const q = `
		WITH next_id AS (
			SELECT nextval(pg_get_serial_sequence('library.loan', 'loan_id'))::integer AS loan_id
		), inserted AS (
			INSERT INTO library.loan (loan_id, loan_ident, book_copy_id, member_id)
			OVERRIDING SYSTEM VALUE
			SELECT next_id.loan_id, 'LN-' || lpad(next_id.loan_id::text, 6, '0'), $1, $2
			FROM next_id
			RETURNING loan_id, loan_ident
		)
		SELECT inserted.loan_id, inserted.loan_ident FROM inserted`

	var loanID int32
	var loanIdent string
	err := r.tx.QueryRow(ctx, q, int32(prepared.BookCopyID), int32(prepared.MemberID)).Scan(&loanID, &loanIdent)
	if err != nil {
		return nil, fmt.Errorf("create loan: %w", err)
	}
	now := time.Now()
	l := domain.Loan{
		ID:         domain.LoanID(loanID),
		Ident:      domain.LoanIdent(loanIdent),
		DtCreated:  now,
		DtModified: now,
		BookCopyID: prepared.BookCopyID,
		MemberID:   prepared.MemberID,
		DtDue:      nil,
		DtReturned: nil,
	}
	return &l, nil
}

func (r *WriteRepo) End(ctx context.Context, id domain.LoanID) error {
	const q = `
		UPDATE library.loan l
		SET dt_returned = CURRENT_TIMESTAMP
		WHERE l.loan_id = $1`

	_, err := r.tx.Exec(ctx, q, int32(id))
	if err != nil {
		return fmt.Errorf("end loan: %w", err)
	}
	return nil
}

func (r *WriteRepo) FindActiveByBookCopyIDForUpdate(ctx context.Context, id domainbc.BookCopyID) (*domain.Loan, error) {
	const activeLoanCols = `l.loan_id, l.loan_ident, l.dt_created, l.dt_modified, l.book_copy_id, l.member_id,
		NULLIF(l.dt_due, '9999-01-01 00:00:00+00'::TIMESTAMPTZ),
		NULLIF(l.dt_returned, '9999-01-01 00:00:00+00'::TIMESTAMPTZ)`
	q := `SELECT ` + activeLoanCols + `
		FROM library.loan l
		WHERE l.book_copy_id = $1
		  AND l.dt_returned = '9999-01-01 00:00:00+00'::TIMESTAMPTZ
		ORDER BY l.loan_id DESC
		LIMIT 1
		FOR UPDATE OF l`

	row, err := scanLoanRow(r.tx.QueryRow(ctx, q, int32(id)))
	if err != nil {
		return nil, fmt.Errorf("find active loan by book copy id for update: %w", err)
	}
	if row == nil {
		return nil, nil
	}
	l := rowToDomain(*row)
	return &l, nil
}

func (r *WriteRepo) CountActiveByMemberID(ctx context.Context, id domainm.MemberID) (int64, error) {
	const q = `
		SELECT COUNT(*)::BIGINT
		FROM library.loan l
		WHERE l.member_id = $1
		  AND l.dt_returned = '9999-01-01 00:00:00+00'::TIMESTAMPTZ`

	var count int64
	if err := r.tx.QueryRow(ctx, q, int32(id)).Scan(&count); err != nil {
		return 0, fmt.Errorf("count active loans by member id (write): %w", err)
	}
	return count, nil
}
