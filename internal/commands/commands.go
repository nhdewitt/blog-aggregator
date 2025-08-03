// Package commands implements the CLI command system for the gator RSS aggregator.
package commands

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/nhdewitt/blog-aggregator/internal/app"
	"github.com/nhdewitt/blog-aggregator/internal/database"
)

// cmdDef holds metadata for a single CLI command.
type cmdDef struct {
	name			string
	args			string
	desc			string
	handler			interface{}		// func(*State, Command) error or func(*State, Command, database.User) error
	requiresLogin	bool
}

// Command represents a parsed command with its arguments.
type Command struct {
	Name		string
	Args		[]string
}

// commandsList defines all available CLI commands.
var commandsList = []cmdDef{
	{"login",		"<username>",		"change current user",					handlerLogin,				false},
	{"register",	"<username>",		"create a new user",					handlerRegister,			false},
	{"reset",		"",					"wipe all user data",					handlerReset,				false},
	{"users",		"",					"list all users",						handlerGetUsers,			false},
	{"agg",			"<duration>",		"continuously aggregate posts",			handlerAggregator,			false},
	{"addfeed",		"<name> <url>",		"add a new feed",						handlerAddFeed,				true},
	{"feeds",		"",					"list all feeds",						handlerPrintAllFeeds,		false},
	{"follow",		"<url>",			"follow an existing feed",				handlerFollow,				true},
	{"following",	"",					"show feeds you're following",			handlerShowFollowedFeeds,	true},
	{"unfollow",	"<url>",			"stop following a feed",				handlerUnfollowFeed,		true},
	{"browse",		"[limit|2]",		"browse your latest <limit> posts",		handlerBrowse,				true},
}

// commandRegistry holds registered command handlers.
type commandRegistry struct {
	handlers		map[string]func(*app.State, Command) error
}

// newCommandRegistry creates a new command registry.
func newCommandRegistry() *commandRegistry {
	return &commandRegistry{
		handlers: make(map[string]func(*app.State, Command) error),
	}
}

// register adds a command handler to the registry.
func (r *commandRegistry) register(name string, handler func(*app.State, Command) error) {
	r.handlers[name] = handler
}

// Execute runs the specified command with the given state.
func Execute(ctx context.Context, state *app.State, cmd Command) error {
	registry := buildCommandRegistry()

	handler, exists := registry.handlers[cmd.Name]
	if !exists {
		return fmt.Errorf("unknown command: %q", cmd.Name)
	}

	return handler(state, cmd)
}

// buildCommandRegistry creates and populates the command registry.
func buildCommandRegistry() *commandRegistry {
	registry := newCommandRegistry()

	for _, cd := range commandsList {
		var handler func(*app.State, Command) error

		if cd.requiresLogin {
			// Wrap handler with authentication middleware
			userHandler := cd.handler.(func(*app.State, Command, database.User) error)
			handler = middlewareLoggedIn(userHandler)
		} else {
			// Use handler directly
			handler = cd.handler.(func(*app.State, Command) error)
		}

		registry.register(cd.name, handler)
	}

	return registry
}

// PrintHelp displays the help message and exits.
func PrintHelp() {
	fmt.Println("NAME")
	fmt.Println("    gator - RSS feed aggregator CLI")
	fmt.Println()
	fmt.Println("SYNOPSIS")
	fmt.Println("    gator <command> [arguments...]")
	fmt.Println()
	fmt.Println("DESCRIPTION")
	fmt.Println("    A command-line tool for managing RSS feeds and aggregating posts.")
	fmt.Println()
	fmt.Println("COMMANDS")

	sortedCommands := make([]cmdDef, len(commandsList))
	copy(sortedCommands, commandsList)
	slices.SortFunc(sortedCommands, func(a, b cmdDef) int {
		return strings.Compare(a.name, b.name)
	})
	
	for _, cd := range sortedCommands {
		if cd.args != "" {
			fmt.Printf("    %-12s %s\n", cd.name, cd.args)
		} else {
			fmt.Printf("    %s\n", cd.name)
		}
		fmt.Printf("        %s\n", cd.desc)
		fmt.Println()
	}
	
	fmt.Println("EXAMPLES")
	fmt.Println("    gator register alice")
	fmt.Println("        Create a new user named 'alice'")
	fmt.Println()
	fmt.Println("    gator addfeed \"Tech News\" https://example.com/feed.xml")
	fmt.Println("        Add a new RSS feed")
	fmt.Println()
	fmt.Println("    gator browse 10")
	fmt.Println("        Browse the latest 10 posts")
	
	os.Exit(1)
}