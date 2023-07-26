package rss

import (
	"time"

	"github.com/mmcdole/gofeed"
)

type RSS struct {
	parser *gofeed.Parser
}

func New() *RSS {
	return &RSS{
		parser: gofeed.NewParser(),
	}
}

func (r *RSS) Parse(
	feedURL string,
	storedLastPublished int64,
) ([]gofeed.Item, int64, error) {
	feed, err := r.parser.ParseURL(feedURL)
	if err != nil {
		return nil, 0, err
	}

	newItems := []gofeed.Item{}

	var lastPublished = storedLastPublished

	for _, item := range feed.Items {
		// Parse the published date of the current item
		pubTime := item.PublishedParsed
		if pubTime != nil &&
			pubTime.After(time.UnixMilli(storedLastPublished)) {
			newItems = append(newItems, *item)
			// Update the last published date to the current item's date
			lastPublished = pubTime.UnixMilli()
		}
	}

	return newItems, lastPublished, nil
}
