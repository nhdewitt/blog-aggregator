// Package commands implements the CLI command system for the gator RSS aggregator.
package commands

import (
	"context"

	"github.com/nhdewitt/blog-aggregator/internal/app"
	"github.com/nhdewitt/blog-aggregator/internal/database"
)

// middlewareLoggedIn wraps command handlers that require user authentication.
// Automatically fetches the current user and passes it to the handler.
func middlewareLoggedIn(handler func(s *app.State, cmd Command, user database.User) error) func(*app.State, Command) error {
	return func(s *app.State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUser)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}