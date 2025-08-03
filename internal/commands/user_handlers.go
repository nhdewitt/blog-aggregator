// Package commands implements the CLI command system for the gator RSS aggregator.
package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nhdewitt/blog-aggregator/internal/app"
	"github.com/nhdewitt/blog-aggregator/internal/database"
)

// handlerLogin authenticates a user by setting them as the current user.
// The user must already exist in the database.
//
// Usage: gator login <username>
func handlerLogin(s *app.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	username := cmd.Args[0]

	_, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Printf("User doesn't exist: %v\n", err)
		os.Exit(1)
	}

	err = s.Cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Error setting username: %v\n", err)
	}

	fmt.Printf("Username set to %s\n", username)
	
	return nil
}

// handlerRegister creates a new user account and sets them as the current user.
// The username must be unique.
//
// Usage: gator register <username>
func handlerRegister(s *app.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	username := cmd.Args[0]
	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: username,
	})
	if err != nil {
		return fmt.Errorf("Couldn't create user: %w", err)
	}

	err = s.Cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Couldn't set current user: %w", err)
	}

	fmt.Println("User created successfully")
	printUser(user)
	return nil
}

// handlerReset removes all users from the database.
// This operation cannot be undone.
//
// Usage: gator reset
func handlerReset(s *app.State, cmd Command) error {
	err := s.Db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't reset database: %w", err)
	}

	fmt.Println("Databse reset")
	return nil
}

// handlerGetUsers displays a list of all registered users.
// The current user is noted with (current).
//
// Usage: gator users
func handlerGetUsers(s *app.State, cmd Command) error {
	currentUser := s.Cfg.CurrentUser

	users, err := s.Db.GetUsers(context.Background())
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

// printUser displays formatted user information.
// Runs when a user is created.
func printUser(u database.User) {
	fmt.Printf("\t* ID:\t%v\n", u.ID)
	fmt.Printf("\t* Name:\t%v\n", u.Name)
}