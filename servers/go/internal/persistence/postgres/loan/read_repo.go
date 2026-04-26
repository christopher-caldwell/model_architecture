package loan

import (
	"context"
	"fmt"

	domainbc "github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
	domain "github.com/christophercaldwell/model-architecture/go/internal/domain/loan"
	domainm "github.com/christophercaldwell/model-architecture/go/internal/domain/member"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReadRepo struct {
	pool *pgxpool.Pool
}

func NewReadRepo(pool *pgxpool.Pool) *ReadRepo {
	return &ReadRepo{pool: pool}
}

const loanCols = `
	l.loan_id, l.loan_ident, l.dt_created, l.dt_modified, l.book_copy_id, l.member_id,
	NULLIF(l.dt_due, '9999-01-01 00:00:00+00'::TIMESTAMPTZ),
	NULLIF(l.dt_returned, '9999-01-01 00:00:00+00'::TIMESTAMPTZ)`

func scanLoanRow(row pgx.Row) (*loanRow, error) {
	var r loanRow
	if err := row.Scan(&r.LoanID, &r.LoanIdent, &r.DtCreated, &r.DtModified, &r.BookCopyID, &r.MemberID, &r.DtDue, &r.DtReturned); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan loan row: %w", err)
	}
	return &r, nil
}

func (r *ReadRepo) GetByMemberIdent(ctx context.Context, ident domainm.MemberIdent) ([]domain.Loan, error) {
	q := `SELECT ` + loanCols + `
		FROM library.loan l
		JOIN library.member m ON l.member_id = m.member_id
		WHERE m.member_ident = $1
		ORDER BY l.dt_created DESC, l.loan_id DESC`

	rows, err := r.pool.Query(ctx, q, string(ident))
	if err != nil {
		return nil, fmt.Errorf("fetch loans by member ident: %w", err)
	}
	defer rows.Close()

	var loans []domain.Loan
	for rows.Next() {
		var row loanRow
		if err := rows.Scan(&row.LoanID, &row.LoanIdent, &row.DtCreated, &row.DtModified, &row.BookCopyID, &row.MemberID, &row.DtDue, &row.DtReturned); err != nil {
			return nil, fmt.Errorf("scan loan row: %w", err)
		}
		loans = append(loans, rowToDomain(row))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate loan rows: %w", err)
	}
	if loans == nil {
		loans = []domain.Loan{}
	}
	return loans, nil
}

func (r *ReadRepo) GetOverdue(ctx context.Context) ([]domain.Loan, error) {
	q := `SELECT ` + loanCols + `
		FROM library.loan l
		WHERE l.dt_returned = '9999-01-01 00:00:00+00'::TIMESTAMPTZ
		  AND l.dt_due < CURRENT_TIMESTAMP
		ORDER BY l.dt_due, l.loan_id`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("fetch overdue loans: %w", err)
	}
	defer rows.Close()

	var loans []domain.Loan
	for rows.Next() {
		var row loanRow
		if err := rows.Scan(&row.LoanID, &row.LoanIdent, &row.DtCreated, &row.DtModified, &row.BookCopyID, &row.MemberID, &row.DtDue, &row.DtReturned); err != nil {
			return nil, fmt.Errorf("scan loan row: %w", err)
		}
		loans = append(loans, rowToDomain(row))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate loan rows: %w", err)
	}
	if loans == nil {
		loans = []domain.Loan{}
	}
	return loans, nil
}

func (r *ReadRepo) FindActiveByBookCopyID(ctx context.Context, id domainbc.BookCopyID) (*domain.Loan, error) {
	q := `SELECT ` + loanCols + `
		FROM library.loan l
		WHERE l.book_copy_id = $1
		  AND l.dt_returned = '9999-01-01 00:00:00+00'::TIMESTAMPTZ
		ORDER BY l.loan_id DESC
		LIMIT 1`

	row, err := scanLoanRow(r.pool.QueryRow(ctx, q, int32(id)))
	if err != nil {
		return nil, fmt.Errorf("find active loan by book copy id: %w", err)
	}
	if row == nil {
		return nil, nil
	}
	l := rowToDomain(*row)
	return &l, nil
}

func (r *ReadRepo) CountActiveByMemberID(ctx context.Context, id domainm.MemberID) (int64, error) {
	const q = `
		SELECT COUNT(*)::BIGINT
		FROM library.loan l
		WHERE l.member_id = $1
		  AND l.dt_returned = '9999-01-01 00:00:00+00'::TIMESTAMPTZ`

	var count int64
	if err := r.pool.QueryRow(ctx, q, int32(id)).Scan(&count); err != nil {
		return 0, fmt.Errorf("count active loans by member id: %w", err)
	}
	return count, nil
}
