package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/christophercaldwell/model-architecture/go/internal/application"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/bookcopy"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/loan"
	"github.com/christophercaldwell/model-architecture/go/internal/domain/member"
)

type CheckOutBookCopyInput struct {
	MemberIdent      string
	BookCopyBarcode  string
}

type LendingCommands struct {
	factory application.UnitOfWorkFactory
}

func NewLendingCommands(factory application.UnitOfWorkFactory) *LendingCommands {
	return &LendingCommands{factory: factory}
}

func (c *LendingCommands) getMemberByIdent(ctx context.Context, uow application.UnitOfWork, ident string) (*member.Member, error) {
	m, err := uow.Members().GetByIdentForUpdate(ctx, member.MemberIdent(ident))
	if err != nil {
		return nil, fmt.Errorf("load member for write: %w", err)
	}
	if m == nil {
		return nil, member.ErrNotFound
	}
	return m, nil
}

func (c *LendingCommands) getBookCopyByBarcode(ctx context.Context, uow application.UnitOfWork, barcode string) (*bookcopy.BookCopy, error) {
	bc, err := uow.BookCopies().GetByBarcodeForUpdate(ctx, barcode)
	if err != nil {
		return nil, fmt.Errorf("load book copy for write: %w", err)
	}
	if bc == nil {
		return nil, bookcopy.ErrNotFound
	}
	return bc, nil
}

func (c *LendingCommands) loadActiveLoanForCopy(ctx context.Context, uow application.UnitOfWork, id bookcopy.BookCopyID) (*loan.Loan, error) {
	l, err := uow.Loans().FindActiveByBookCopyIDForUpdate(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find active loan for book copy: %w", err)
	}
	return l, nil
}

func (c *LendingCommands) CheckOutBookCopy(ctx context.Context, input CheckOutBookCopyInput) (*loan.Loan, error) {
	uow, err := c.factory.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("build unit of work: %w", err)
	}
	defer uow.Rollback(ctx) //nolint:errcheck

	m, err := c.getMemberByIdent(ctx, uow, input.MemberIdent)
	if err != nil {
		return nil, err
	}
	bc, err := c.getBookCopyByBarcode(ctx, uow, input.BookCopyBarcode)
	if err != nil {
		return nil, err
	}

	if err := m.EnsureCanBorrow(); err != nil {
		return nil, err
	}
	if err := bc.EnsureCanBeBorrowed(); err != nil {
		return nil, err
	}

	activeCount, err := uow.Loans().CountActiveByMemberID(ctx, m.ID)
	if err != nil {
		return nil, fmt.Errorf("count active loans for member: %w", err)
	}
	if err := m.EnsureWithinLoanLimit(int16(activeCount)); err != nil {
		return nil, err
	}

	activeLoan, err := c.loadActiveLoanForCopy(ctx, uow, bc.ID)
	if err != nil {
		return nil, err
	}
	if activeLoan != nil {
		return nil, bookcopy.ErrCannotBeBorrowed
	}

	prepared := loan.LoanCreationPayload{
		MemberID:   m.ID,
		BookCopyID: bc.ID,
	}.Prepare()
	result, err := uow.Loans().Create(ctx, prepared)
	if err != nil {
		return nil, fmt.Errorf("check out book copy: %w", err)
	}
	if err := uow.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}

func (c *LendingCommands) ReturnBookCopy(ctx context.Context, barcode string) (*loan.Loan, error) {
	uow, err := c.factory.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("build unit of work: %w", err)
	}
	defer uow.Rollback(ctx) //nolint:errcheck

	bc, err := c.getBookCopyByBarcode(ctx, uow, barcode)
	if err != nil {
		return nil, err
	}
	activeLoan, err := c.loadActiveLoanForCopy(ctx, uow, bc.ID)
	if err != nil {
		return nil, err
	}
	if activeLoan == nil {
		return nil, loan.ErrNoActiveLoanForBookCopy
	}
	if err := activeLoan.EnsureCanBeReturned(); err != nil {
		return nil, err
	}
	if err := uow.Loans().End(ctx, activeLoan.ID); err != nil {
		return nil, fmt.Errorf("return book copy: %w", err)
	}
	if err := uow.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	now := time.Now()
	updated := *activeLoan
	updated.DtModified = now
	updated.DtReturned = &now
	return &updated, nil
}

func (c *LendingCommands) ReportLostLoanedBookCopy(ctx context.Context, barcode string) (*bookcopy.BookCopy, error) {
	uow, err := c.factory.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("build unit of work: %w", err)
	}
	defer uow.Rollback(ctx) //nolint:errcheck

	bc, err := c.getBookCopyByBarcode(ctx, uow, barcode)
	if err != nil {
		return nil, err
	}
	newStatus, err := bc.MarkLost()
	if err != nil {
		return nil, err
	}
	activeLoan, err := c.loadActiveLoanForCopy(ctx, uow, bc.ID)
	if err != nil {
		return nil, err
	}
	if activeLoan == nil {
		return nil, loan.ErrNoActiveLoanForBookCopy
	}
	if err := activeLoan.EnsureCanBeReturned(); err != nil {
		return nil, err
	}
	if err := uow.Loans().End(ctx, activeLoan.ID); err != nil {
		return nil, fmt.Errorf("close lost loan: %w", err)
	}
	if err := uow.BookCopies().UpdateStatus(ctx, bc.ID, newStatus); err != nil {
		return nil, fmt.Errorf("mark book copy as lost: %w", err)
	}
	if err := uow.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	updated := *bc
	updated.Status = newStatus
	updated.DtModified = time.Now()
	return &updated, nil
}
