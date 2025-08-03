package main

import (
	"fmt"

	"github.com/nhdewitt/blog-aggregator/internal/config"
	"github.com/nhdewitt/blog-aggregator/internal/database"
)

type state struct {
	cfg	*config.Config
	db	*database.Queries
}

type command struct {
	CmdName		string
	Args		[]string
}

type commands struct {
	CommandMap		map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.CommandMap[cmd.CmdName]
	if !ok {
		return fmt.Errorf("Command doesn't exist: %v\n", cmd.CmdName)
	}

	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.CommandMap[name] = f
}