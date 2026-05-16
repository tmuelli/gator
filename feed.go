package main

import (
	"context"
	"net/http"
	"encoding/xml"
	"html"
	"io"
	"fmt"
	"database/sql"

	"github.com/lib/pq"
	"github.com/tmuelli/blog-aggregator/internal/convert"
	"github.com/tmuelli/blog-aggregator/internal/database"
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

func scrapFeeds(s *state) {
	// get next feed to fetch
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Println("Log - error:", err)
		return
	}

	// mark feed as fetched
	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		fmt.Println("Log - error:", err)
		return
	}

	// fetch the feed
	rssFeed, err := fetchFeed(context.Background(), feed.Url.String)
	if err != nil {
		fmt.Println("Log - error:", err)
		return
	}

	fmt.Println("Fetched feed:", rssFeed.Channel.Title, "- items:", len(rssFeed.Channel.Item))

	// insert feed posts
	for _, item := range rssFeed.Channel.Item {
		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			Title:			sql.NullString{String: item.Title, Valid: true},
			Url:			sql.NullString{String: item.Link, Valid: true},
			Description:	sql.NullString{String: item.Description, Valid: true},
			PublishedAt:	convert.DateStringToSqlNullTime(item.PubDate),
			FeedID:			feed.ID,
		})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				continue
			}

			fmt.Println("Error creating post:", err)
		}
	}
}