package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nhdewitt/blog-aggregator/internal/database"
)

const (
	timeNoLeadingZero = "Mon, 2 Jan 2006 15:04:05 MST"
	timeNoLeadingZeroZ = "Mon, 2 Jan 2006 15:04:05 -0700"
)

func scrapeFeeds(s *state) error {
	c := context.Background()
	nextFeed, err := s.db.GetNextFeedToFetch(c)
	if err != nil {
		return fmt.Errorf("Unable to grab next feed to fetch: %w", err)
	}

	feed, err := fetchFeed(c, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("Unable to fetch feed: %w", err)
	}
	feedID := nextFeed.ID

	for _, item := range feed.Channel.Item {
		description := sql.NullString{
			String: item.Description,
			Valid: item.Description != "",
		}
		validTimeStrings := []string{
			time.RFC1123,			// "Mon 02 Jan 2006 15:04:05 MST"
			time.RFC1123Z,			// "Mon 02 Jan 2006 15:04:05 -0700"
			timeNoLeadingZero,		// "Mon 2 Jan 2006 15:04:05 MST"
			timeNoLeadingZeroZ,		// "Mon 2 Jan 2006 15:04:05 -0700"
		}
		var parsedPubDate time.Time
		for _, timeString := range validTimeStrings {
			parsedPubDate, err = time.Parse(timeString, item.PubDate)
			if err == nil {
				break
			}
		}
		if err != nil {
			log.Printf("Could not parse pubDate %s from feed %s: %s\n", item.PubDate, feed.Channel.Title, err)
		}

		err = s.db.CreatePost(c, database.CreatePostParams{
			ID: uuid.New(),
			Title: item.Title,
			Url: item.Link,
			Description: description,
			PublishedAt: parsedPubDate.UTC(),
			FeedID: feedID,
		})
		if err != nil {
			log.Printf("Could not create post: %s\n", err)
		}
	}

	err = s.db.MarkFeedFetched(c, nextFeed.ID)
	if err != nil {
		return fmt.Errorf("Unable to update fetched feed: %w", err)
	}
	return nil
}