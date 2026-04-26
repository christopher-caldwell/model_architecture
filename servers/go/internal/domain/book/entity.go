package book

import "time"

type BookID int32

type Book struct {
	ID         BookID
	ISBN       string
	DtCreated  time.Time
	DtModified time.Time
	Title      string
	AuthorName string
}

type BookCreationPayload struct {
	ISBN       string
	Title      string
	AuthorName string
}

type BookPrepared struct {
	ISBN       string
	Title      string
	AuthorName string
}
