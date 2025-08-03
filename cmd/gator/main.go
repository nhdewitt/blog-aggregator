package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/nhdewitt/blog-aggregator/internal/config"
	"github.com/nhdewitt/blog-aggregator/internal/database"
)

func main() {	
	// Config: Read json from ~/.gatorconfig.json to get:
	//	DBUrl: URL for DB connection
	//	CurrentUserName: Currently logged in user
	c, err := config.Read()
	if err != nil {
		log.Fatal("Error reading config:", err)
	}


	// DB: open DB connection, dbQueries struct for DB queries
	d, err := sql.Open("postgres", c.DBUrl)
	dbQueries := database.New(d)

	
	// State: Holds the config and DB queries to be passed to the functions
	s := &state{
		cfg: &c,
		db: dbQueries,
	}
	
	cmds := commands{
		CommandMap: make(map[string]func(*state, command) error),
	}
	// Cmds handlers:
	//	login <username> changes the user, Exit(1) if user doesn't exist
	//	register <username> registers the user, updates cfg.CurrentUserName
	//	reset resets the user DB (cascades into the rest, clears out everything)
	//	users shows the list of registered users, marks the currently logged in user
	//	agg <duration> runs until interrupt - grabs from feed that was updated the earliest
	//	addfeed <feed name> <url> adds the feed to the feeds table
	//	follow <url> follows a url that is already in the DB - Exit(1) if feed is not present
	//	following shows a list of feeds followed by the user
	//	unfollow <url> unfollows a feed for the user
	//	browse <limit|2> grabs the latest <limit> posts from posts DB
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAggregator)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerPrintAllFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerShowFollowedFeeds))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollowFeed))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))
	
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(s, command{
		cmdName, cmdArgs,
	})
	if err != nil {
		log.Fatal(err)
	}
}