package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/nhdewitt/blog-aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.CmdName)
	}

	username := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Printf("User doesn't exist: %v\n", err)
		os.Exit(1)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Error setting username: %v\n", err)
	}

	fmt.Printf("Username set to %s\n", username)
	
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.CmdName)
	}

	username := cmd.Args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: username,
	})
	if err != nil {
		return fmt.Errorf("Couldn't create user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Couldn't set current user: %w", err)
	}

	fmt.Println("User created successfully")
	printUser(user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't reset database: %w", err)
	}

	fmt.Println("Databse reset")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	currentUser := s.cfg.CurrentUser

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't get user list: %w", err)
	}

	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf(" * %s (current)\n", user.Name)
		} else {
			fmt.Printf(" * %s\n", user.Name)
		}
	}

	return nil
}

func handlerAggregator(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.CmdName)
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
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feed name> <feed url>", cmd.CmdName)
	}

	c := context.Background()

	newFeed, err := s.db.AddFeed(c, database.AddFeedParams{
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

	_, err = s.db.CreateFeedFollow(c, database.CreateFeedFollowParams{
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

func handlerPrintAllFeeds(s *state, cmd command) error {
	c := context.Background()
	feeds, err := s.db.PrintAllFeeds(c)
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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.CmdName)
	}
	url := cmd.Args[0]

	c := context.Background()
	followedFeed, err := s.db.FindFeedsByURL(c, url)
	if err != nil {
		return fmt.Errorf("Error finding feed: %w", err)
	}

	createdFollow, err := s.db.CreateFeedFollow(c, database.CreateFeedFollowParams{
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

func handlerShowFollowedFeeds(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.CmdName)
	}
	c := context.Background()
	
	id := user.ID
	feeds, err := s.db.GetFeedFollowsForUser(c, id)
	if err != nil {
		return fmt.Errorf("Error getting user's feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("You are not subscribed to any feeds")
	} else {
		fmt.Printf("Subscribed feeds for %s:\n", s.cfg.CurrentUser)
	}
	for _, feed := range feeds {
		fmt.Println(feed.Name)
	}

	return nil
}

func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.CmdName)
	}
	c := context.Background()

	id := user.ID
	url := cmd.Args[0]

	err := s.db.UnfollowFeed(c, database.UnfollowFeedParams{
		ID: id,
		Url: url,
	})
	if err != nil {
		return fmt.Errorf("Error unfollowing feed: %w", err)
	}

	fmt.Println("You have unfollowed the feed")
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
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

	posts, err := s.db.GetPostsForUser(c, database.GetPostsForUserParams{
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

func printUser(u database.User) {
	fmt.Printf("\t* ID:\t%v\n", u.ID)
	fmt.Printf("\t* Name:\t%v\n", u.Name)
}