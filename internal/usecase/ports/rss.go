package ports

import (
	"context"

	domain "hufschlaeger.net/markscribe/internal/domain"
)

// RssFeedPort defines operations we use from literal.club integration.
type RssFeedPort interface {
	RecentFeedEntries(ctx context.Context, url string, count int) ([]domain.RSSEntry, error)
}
