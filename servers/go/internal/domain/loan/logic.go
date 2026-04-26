package loan

func (l *Loan) EnsureCanBeReturned() error {
	if l.DtReturned != nil {
		return ErrCannotBeReturned
	}
	return nil
}

func (p LoanCreationPayload) Prepare() LoanPrepared {
	return LoanPrepared{
		MemberID:   p.MemberID,
		BookCopyID: p.BookCopyID,
	}
}
