package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel	struct {
		Title		string		`xml:"title"`
		Link		string		`xml:"link"`
		Description	string		`xml:"description"`
		Item		[]RSSItem	`xml:"item"`
	}	`xml:"channel"`
}

type RSSItem struct {
	Title 		string		`xml:"title"`
	Link		string		`xml:"link"`
	Description	string		`xml:"description"`
	PubDate		string		`xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	var rss_feed *RSSFeed

	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return rss_feed, fmt.Errorf("error creating request : %w", err)
	}
	
	req.Header.Set("User-Agent", "gator")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return rss_feed, fmt.Errorf("error sending request: %w", err)
	}

	defer res.Body.Close()


	body, err := io.ReadAll(res.Body)
	if err != nil {
		return rss_feed, fmt.Errorf("error reading response body: %w", err)
	}

	if err := xml.Unmarshal(body, &rss_feed); err != nil {
		return rss_feed, err
	}

	rss_feed.Channel.Title = html.UnescapeString(rss_feed.Channel.Title)
	rss_feed.Channel.Description = html.UnescapeString(rss_feed.Channel.Description)
	for i := 0; i < len(rss_feed.Channel.Item); i++ {
		rss_feed.Channel.Item[i].Title = html.UnescapeString(rss_feed.Channel.Item[i].Title)
		rss_feed.Channel.Item[i].Description = html.UnescapeString(rss_feed.Channel.Item[i].Description)
	}

	return rss_feed, nil
}