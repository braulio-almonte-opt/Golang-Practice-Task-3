package models

type SearchOptions struct {
	Skip    int64
	Limit   int64
	Fields  []string
	Exclude []string
}
