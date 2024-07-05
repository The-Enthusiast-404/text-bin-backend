package data

import "time"

// Its important in Go to keep the Fields of a struct in Capotal letter to make it public
// Any field that starts with a lowercase letter is private to the package and aren't  exported and won't be included when encoding a struct to JSON
type Text struct {
	ID        int64
	CreatedAt time.Time
	Title     string
	Content   string
	Format    string
	Version   int32
}
