export const sql = {
  book: {
    create: `
      INSERT INTO library.book (isbn, title, author_name)
      VALUES ($1, $2, $3)
      RETURNING book_id
    `,
    getByIsbn: `
      SELECT book_id, isbn, dt_created, dt_modified, title, author_name
      FROM library.book
      WHERE isbn = $1
    `,
    getCatalog: `
      SELECT book_id, isbn, dt_created, dt_modified, title, author_name
      FROM library.book
      ORDER BY title, book_id
    `
  },
  bookCopy: {
    create: `
      INSERT INTO library.book_copy (book_id, status_id, barcode)
      VALUES (
        $1,
        (
          SELECT struct_type_id
          FROM library.struct_type
          WHERE group_name = 'book_copy_status'
            AND att_pub_ident = $2
        ),
        $3
      )
      RETURNING book_copy_id
    `,
    getById: `
      SELECT bc.book_copy_id, bc.barcode, bc.dt_created, bc.dt_modified, bc.book_id, st.att_pub_ident AS status
      FROM library.book_copy bc
      JOIN library.struct_type st ON bc.status_id = st.struct_type_id
      WHERE bc.book_copy_id = $1
    `,
    getByBarcode: `
      SELECT bc.book_copy_id, bc.barcode, bc.dt_created, bc.dt_modified, bc.book_id, st.att_pub_ident AS status
      FROM library.book_copy bc
      JOIN library.struct_type st ON bc.status_id = st.struct_type_id
      WHERE bc.barcode = $1
    `,
    getByBarcodeForUpdate: `
      SELECT bc.book_copy_id, bc.barcode, bc.dt_created, bc.dt_modified, bc.book_id, st.att_pub_ident AS status
      FROM library.book_copy bc
      JOIN library.struct_type st ON bc.status_id = st.struct_type_id
      WHERE bc.barcode = $1
      FOR UPDATE OF bc
    `,
    updateStatus: `
      UPDATE library.book_copy
      SET status_id = (
        SELECT struct_type_id
        FROM library.struct_type
        WHERE group_name = 'book_copy_status'
          AND att_pub_ident = $2
      )
      WHERE book_copy_id = $1
    `
  },
  member: {
    create: `
      INSERT INTO library.member (member_ident, status_id, full_name, max_active_loans)
      VALUES (
        $1,
        (
          SELECT struct_type_id
          FROM library.struct_type
          WHERE group_name = 'member_status'
            AND att_pub_ident = $2
        ),
        $3,
        $4
      )
      RETURNING member_id
    `,
    getById: `
      SELECT m.member_id, m.member_ident, m.dt_created, m.dt_modified, st.att_pub_ident AS status, m.full_name, m.max_active_loans
      FROM library.member m
      JOIN library.struct_type st ON m.status_id = st.struct_type_id
      WHERE m.member_id = $1
        AND st.group_name = 'member_status'
        AND st.att_pub_ident IN ('active', 'suspended')
    `,
    getByIdent: `
      SELECT m.member_id, m.member_ident, m.dt_created, m.dt_modified, st.att_pub_ident AS status, m.full_name, m.max_active_loans
      FROM library.member m
      JOIN library.struct_type st ON m.status_id = st.struct_type_id
      WHERE m.member_ident = $1
        AND st.group_name = 'member_status'
        AND st.att_pub_ident IN ('active', 'suspended')
    `,
    getByIdentForUpdate: `
      SELECT m.member_id, m.member_ident, m.dt_created, m.dt_modified, st.att_pub_ident AS status, m.full_name, m.max_active_loans
      FROM library.member m
      JOIN library.struct_type st ON m.status_id = st.struct_type_id
      WHERE m.member_ident = $1
        AND st.group_name = 'member_status'
        AND st.att_pub_ident IN ('active', 'suspended')
      FOR UPDATE OF m
    `,
    updateStatus: `
      UPDATE library.member
      SET status_id = (
        SELECT struct_type_id
        FROM library.struct_type
        WHERE group_name = 'member_status'
          AND att_pub_ident = $2
      )
      WHERE member_id = $1
    `
  },
  loan: {
    create: `
      WITH next_id AS (
        SELECT nextval(pg_get_serial_sequence('library.loan', 'loan_id'))::integer AS loan_id
      ), inserted AS (
        INSERT INTO library.loan (loan_id, loan_ident, book_copy_id, member_id)
        OVERRIDING SYSTEM VALUE
        SELECT next_id.loan_id, 'LN-' || lpad(next_id.loan_id::text, 6, '0'), $1, $2
        FROM next_id
        RETURNING loan_id, loan_ident
      )
      SELECT loan_id, loan_ident
      FROM inserted
    `,
    end: `
      UPDATE library.loan
      SET dt_returned = CURRENT_TIMESTAMP
      WHERE loan_id = $1
    `,
    findActiveByBookCopyId: `
      SELECT loan_id, loan_ident, dt_created, dt_modified, book_copy_id, member_id,
        NULLIF(dt_due, '9999-01-01 00:00:00+00'::TIMESTAMPTZ) AS dt_due,
        NULLIF(dt_returned, '9999-01-01 00:00:00+00'::TIMESTAMPTZ) AS dt_returned
      FROM library.loan
      WHERE book_copy_id = $1
        AND dt_returned = '9999-01-01 00:00:00+00'::TIMESTAMPTZ
      ORDER BY loan_id DESC
      LIMIT 1
    `,
    findActiveByBookCopyIdForUpdate: `
      SELECT l.loan_id, l.loan_ident, l.dt_created, l.dt_modified, l.book_copy_id, l.member_id,
        NULLIF(l.dt_due, '9999-01-01 00:00:00+00'::TIMESTAMPTZ) AS dt_due,
        NULLIF(l.dt_returned, '9999-01-01 00:00:00+00'::TIMESTAMPTZ) AS dt_returned
      FROM library.loan l
      WHERE l.book_copy_id = $1
        AND l.dt_returned = '9999-01-01 00:00:00+00'::TIMESTAMPTZ
      ORDER BY l.loan_id DESC
      LIMIT 1
      FOR UPDATE OF l
    `,
    countActiveByMemberId: `
      SELECT COUNT(*)::BIGINT AS count
      FROM library.loan
      WHERE member_id = $1
        AND dt_returned = '9999-01-01 00:00:00+00'::TIMESTAMPTZ
    `,
    getByMemberIdent: `
      SELECT l.loan_id, l.loan_ident, l.dt_created, l.dt_modified, l.book_copy_id, l.member_id,
        NULLIF(l.dt_due, '9999-01-01 00:00:00+00'::TIMESTAMPTZ) AS dt_due,
        NULLIF(l.dt_returned, '9999-01-01 00:00:00+00'::TIMESTAMPTZ) AS dt_returned
      FROM library.loan l
      JOIN library.member m ON l.member_id = m.member_id
      WHERE m.member_ident = $1
      ORDER BY l.dt_created DESC, l.loan_id DESC
    `,
    getOverdue: `
      SELECT loan_id, loan_ident, dt_created, dt_modified, book_copy_id, member_id,
        NULLIF(dt_due, '9999-01-01 00:00:00+00'::TIMESTAMPTZ) AS dt_due,
        NULLIF(dt_returned, '9999-01-01 00:00:00+00'::TIMESTAMPTZ) AS dt_returned
      FROM library.loan
      WHERE dt_returned = '9999-01-01 00:00:00+00'::TIMESTAMPTZ
        AND dt_due < CURRENT_TIMESTAMP
      ORDER BY dt_due, loan_id
    `
  }
} as const;
