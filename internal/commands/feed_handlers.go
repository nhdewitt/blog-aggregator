// Package commands implements the CLI command system for the gator RSS aggregator.
package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nhdewitt/blog-aggregator/internal/app"
	"github.com/nhdewitt/blog-aggregator/internal/database"
)

// handlerAddFeed creates a new RSS feed and automatically follows it for the current user.
// The feed URL must be valid and accessible.
//
// Usage: gator addfeed <feed_name> <feed_url>
func handlerAddFeed(s *app.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feed name> <feed url>", cmd.Name)
	}

	c := context.Background()

	newFeed, err := s.Db.AddFeed(c, database.AddFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: cmd.Args[0],
		Url: cmd.Args[1],
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("Error adding feed: %w", err)
	}

	_, err = s.Db.CreateFeedFollow(c, database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: newFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error creating feed follow: %w", err)
	}

	fmt.Printf("New Feed Added: %+v\n", newFeed)
	return nil
}

// handlerPrintAllFeeds displays all feeds in the system with their creators.
// Shows feed name, URL, and thje user who added it.
//
// Usage: gator feeds
func handlerPrintAllFeeds(s *app.State, cmd Command) error {
	c := context.Background()
	feeds, err := s.Db.PrintAllFeeds(c)
	if err != nil {
		return fmt.Errorf("Error retrieving all feeds: %w", err)
	}

	for _, feed := range feeds {
		fmt.Printf(" * Feed:\t%s\n", feed.FeedName)
		fmt.Printf(" * URL:\t\t%s\n", feed.FeedUrl)
		fmt.Printf(" * Added by:\t%s\n", feed.UserName)
		fmt.Println()
	}
	return nil
}

// handlerFollow allows the current user to follow an existing feed by URL.
// The feed must already exist in the syustem.
//
// Usage: gator follow <url>
func handlerFollow(s *app.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]

	c := context.Background()
	followedFeed, err := s.Db.FindFeedsByURL(c, url)
	if err != nil {
		return fmt.Errorf("Error finding feed: %w", err)
	}

	createdFollow, err := s.Db.CreateFeedFollow(c, database.CreateFeedFollowParams{
		ID:	uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: followedFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error following feed: %w", err)
	}

	fmt.Printf("User %s is now following feed %s.\n", createdFollow.UserName, createdFollow.FeedName)
	return nil
}

// handlerShowFollowedFeeds displays all feeds that the current user is following.
// Shows the name of each feed.
//
// Usage: gator following
func handlerShowFollowedFeeds(s *app.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	c := context.Background()
	
	id := user.ID
	feeds, err := s.Db.GetFeedFollowsForUser(c, id)
	if err != nil {
		return fmt.Errorf("Error getting user's feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("You are not subscribed to any feeds")
	} else {
		fmt.Printf("Subscribed feeds for %s:\n", s.Cfg.CurrentUser)
	}
	for _, feed := range feeds {
		fmt.Println(feed.Name)
	}

	return nil
}

// handlerUnfollowFeed removes the current user's subscription to a feed.
// The user will no longer receive posts from this feed.
//
// Usage: gator unfollow <url>
func handlerUnfollowFeed(s *app.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	c := context.Background()

	id := user.ID
	url := cmd.Args[0]

	err := s.Db.UnfollowFeed(c, database.UnfollowFeedParams{
		ID: id,
		Url: url,
	})
	if err != nil {
		return fmt.Errorf("Error unfollowing feed: %w", err)
	}

	fmt.Println("You have unfollowed the feed")
	return nil
}