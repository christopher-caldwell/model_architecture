package member

func (m *Member) Suspend() (MemberStatus, error) {
	if m.Status == MemberStatusSuspended {
		return "", ErrCannotBeSuspended
	}
	return MemberStatusSuspended, nil
}

func (m *Member) Reactivate() (MemberStatus, error) {
	if m.Status != MemberStatusSuspended {
		return "", ErrCannotBeReactivated
	}
	return MemberStatusActive, nil
}

func (m *Member) EnsureCanBorrow() error {
	if m.Status != MemberStatusActive {
		return ErrCannotBorrowWhileSuspended
	}
	return nil
}

func (m *Member) EnsureWithinLoanLimit(activeLoanCount int16) error {
	if activeLoanCount >= m.MaxActiveLoans {
		return ErrLoanLimitReached
	}
	return nil
}

func (p MemberCreationPayload) Prepare(ident MemberIdent) MemberPrepared {
	return MemberPrepared{
		Ident:          ident,
		FullName:       p.FullName,
		MaxActiveLoans: p.MaxActiveLoans,
		Status:         MemberStatusActive,
	}
}
