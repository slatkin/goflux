package ui

import (
	"github.com/slatkin/goflux/pkg/miniflux"
)

type ErrorMsg error

type EntriesMsg []miniflux.FeedEntry

type EntryContentMsg struct {
	EntryID int
	Content string
}

type ActionDoneMsg struct {
	Action  string
	Success bool
}
