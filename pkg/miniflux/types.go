package miniflux

import "time"

type ReadStatus string

const (
	ReadStatusRead   ReadStatus = "read"
	ReadStatusUnread ReadStatus = "unread"
)

func (r ReadStatus) Toggle() ReadStatus {
	if r == ReadStatusRead {
		return ReadStatusUnread
	}
	return ReadStatusRead
}

type Feed struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	SiteURL string `json:"site_url"`
	FeedURL string `json:"feed_url"`
}

type FeedEntry struct {
	ID          int        `json:"id"`
	FeedID      int        `json:"feed_id"`
	Title       string     `json:"title"`
	URL         string     `json:"url"`
	Content     string     `json:"content"`
	Feed        Feed       `json:"feed"`
	Status      ReadStatus `json:"status"`
	Starred     bool       `json:"starred"`
	PublishedAt time.Time  `json:"published_at"`
	// OriginalContent is optional
	OriginalContent string `json:"original_content,omitempty"`
}

type FeedEntriesResponse struct {
	Total   int         `json:"total"`
	Entries []FeedEntry `json:"entries"`
}

type UpdateEntriesRequest struct {
	Status   string `json:"status,omitempty"`
	EntryIDs []int  `json:"entry_ids"`
}

type OriginalContentResponse struct {
	Content string `json:"content"`
}
