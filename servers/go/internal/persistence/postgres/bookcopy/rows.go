package bookcopy

import "time"

type bookCopyRow struct {
	BookCopyID int32
	Barcode    string
	DtCreated  time.Time
	DtModified time.Time
	BookID     int32
	Status     string
}
