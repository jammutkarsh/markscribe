package main

import (
	"time"

	"github.com/mmcdole/gofeed"
)

// RSSEntry represents a single RSS entry.
type RSSEntry struct {
	Title       string
	URL         string
	Description string
	PublishedAt time.Time
}

func rssFeed(url string, count int) ([]RSSEntry, error) {
	var r []RSSEntry

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}

	for _, v := range feed.Items {
		// fmt.Printf("%+v\n", v)

		var published time.Time
		if v.PublishedParsed != nil {
			published = *v.PublishedParsed
		}

		r = append(r, RSSEntry{
			Title:       v.Title,
			URL:         v.Link,
			Description: v.Description,
			PublishedAt: published,
		})
		if len(r) == count {
			break
		}
	}

	return r, nil
}
