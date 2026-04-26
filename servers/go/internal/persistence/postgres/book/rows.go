package book

import "time"

type bookRow struct {
	BookID     int32
	ISBN       string
	DtCreated  time.Time
	DtModified time.Time
	Title      string
	AuthorName string
}
