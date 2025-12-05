package rss

import (
	"github.com/mmcdole/gofeed"
	"hufschlaeger.net/markscribe/internal/domain"
)

type Adapter struct {
	parser *gofeed.Parser
}

func New() *Adapter { // constructor kept simple for now
	feedParser := gofeed.NewParser()
	return &Adapter{parser: feedParser}
}

func (a *Adapter) RecentFeedEntries(url string, count int) (
	[]domain.RSSEntry, error) {

	var r []domain.RSSEntry

	feed, err := a.parser.ParseURL(url)
	if err != nil {
		panic(err)
	}

	for _, v := range feed.Items {
		r = append(r, domain.RSSEntry{
			Title:       v.Title,
			URL:         v.Link,
			PublishedAt: *v.PublishedParsed,
		})
		if len(r) == count {
			break
		}
	}

	return r, nil
}
