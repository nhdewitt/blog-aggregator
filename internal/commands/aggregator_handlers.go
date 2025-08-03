// Package commands implements the CLI command system for the gator RSS aggregator.
package commands

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/nhdewitt/blog-aggregator/internal/app"
	"github.com/nhdewitt/blog-aggregator/internal/database"
)

// handlerBrowse displays the latest posts from feeds the current user follows.
// Posts are shown in reverse chronological order with title, publication date,
// description (if available), and URL.
//
// Usage: gator browse <limit>
// Default for limit is 2
func handlerBrowse(s *app.State, cmd Command, user database.User) error {
	c := context.Background()

	var limit int32 = 2
	if len(cmd.Args) > 0 {
		l, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("Please enter a number for the limit: %w", err)
		}
		limit = int32(l)
	}
	id := user.ID

	posts, err := s.Db.GetPostsForUser(c, database.GetPostsForUserParams{
		UserID : id,
		Limit: limit,
	})
	if err != nil {
		return fmt.Errorf("Error getting posts for user: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("Title: %s (published on %s at %s)\n\n", post.Title, post.PublishedAt.Format("Jan 2, 2006"), post.PublishedAt.Format("3:04 PM"))
		if post.Description.Valid {
			fmt.Printf("Description: %s\n", post.Description.String)
		}
		fmt.Printf("URL: %s\n", post.Url)
		fmt.Println()
	}

	return nil
}

// handlerAggregator starts a continuous process that fetches posts from all feeds.
// It runs indefinitely, fetching new posts at the specified interval.
//
// Usage: gator agg <duration>
// Example: gator agg "1m" (every minute), gator agg "1h" (every hour)
func handlerAggregator(s *app.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}
	timeBetweenReqs := cmd.Args[0]

	var duration time.Duration
	duration, err := time.ParseDuration(timeBetweenReqs)
	if err != nil {
		return fmt.Errorf("Please enter duration in the form (1s|1m|1h)")
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenReqs)
	ticker := time.NewTicker(duration)
	for {
		<-ticker.C
		app.ScrapeFeeds(s)
	}
}













