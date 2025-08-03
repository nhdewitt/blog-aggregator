// Package app contains shared application services and state management.
package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nhdewitt/blog-aggregator/internal/database"
)

// Time format constants for parsing various RSS date formats
const (
	timeNoLeadingZero = "Mon, 2 Jan 2006 15:04:05 MST"		// RFC1123 without leading zero
	timeNoLeadingZeroZ = "Mon, 2 Jan 2006 15:04:05 -0700"	// RFC1123Z without leading zero
)

// ScrapeFeeds fetches the next feed from the database and processes all its posts.
// It retrieves the RSS feed content, parses each post item, and stores new posts
// in the database. The feed is marked as fetched after successful processing.
//
// This function handles various RSS date formates and gracefully handles parsing errors
// by logging them and continuing with the next post. HTML entities in titles and
// descriptions are automatically unescaped.
//
// Returns an error if the feed cannot be fetched, parsed, or if db operations fail.
func ScrapeFeeds(s *State) error {
	c := context.Background()

	// Get the next feed
	nextFeed, err := s.Db.GetNextFeedToFetch(c)
	if err != nil {
		return fmt.Errorf("Unable to grab next feed to fetch: %w", err)
	}

	// Fetch and parse the RSS feed
	feed, err := fetchFeed(c, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("Unable to fetch feed: %w", err)
	}
	feedID := nextFeed.ID

	// Process each item in the feed
	for _, item := range feed.Channel.Item {
		// Handle optional description field
		description := sql.NullString{
			String: item.Description,
			Valid: item.Description != "",
		}
		// Try multiple date formats
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

		// Create the post in the database
		err = s.Db.CreatePost(c, database.CreatePostParams{
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

	// Mark the feed as successfully fetched
	err = s.Db.MarkFeedFetched(c, nextFeed.ID)
	if err != nil {
		return fmt.Errorf("Unable to update fetched feed: %w", err)
	}
	return nil
}