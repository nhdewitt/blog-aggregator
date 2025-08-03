// Gator is a command-line RSS feed aggregator that allows users to manage
// and follow RSS feeds, browse posts, and aggregate content.
//
// Usage:
//
//	gator <command> [arguments...]
//
// Available commands include login, register, addfeed, follow, browse, and more.
// Run 'gator' without arguments to see the full help.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/nhdewitt/blog-aggregator/internal/app"
	"github.com/nhdewitt/blog-aggregator/internal/commands"
	"github.com/nhdewitt/blog-aggregator/internal/config"
	"github.com/nhdewitt/blog-aggregator/internal/database"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load configuration
	c, err := config.Read()
	if err != nil {
		log.Fatal("Error reading config:", err)
	}

	// Initialize database
	d, err := sql.Open("postgres", c.DBUrl)
	dbQueries := database.New(d)
	defer d.Close()

	
	// Create application state
	s := &app.State{
		Cfg: &c,
		Db: dbQueries,
	}

	// Parse command line arguments
	if len(os.Args) < 2 {
		commands.PrintHelp()
		return nil
	}

	cmd := commands.Command{
		Name:		os.Args[1],
		Args:			os.Args[2:],
	}

	return commands.Execute(context.Background(), s, cmd)
}