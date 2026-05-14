package main

import (
	"context"
	"net/http"
	"encoding/xml"
	"html"
	"io"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// fetch from url
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	// set additional header
	req.Header.Set("User-Agent", "gator")

	// send request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()


	// read response
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// decode Body
	var feed RSSFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, err
	}

	// unescape strings
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for _, i := range feed.Channel.Item {
		i.Title = html.UnescapeString(i.Title)
		i.Description = html.UnescapeString(i.Description)
	}
	
	return &feed, nil
}