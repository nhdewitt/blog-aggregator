// Package app contains shared application services and state management.
package app

import (
	"github.com/nhdewitt/blog-aggregator/internal/config"
	"github.com/nhdewitt/blog-aggregator/internal/database"
)


type State struct {
	Cfg	*config.Config
	Db	*database.Queries
}