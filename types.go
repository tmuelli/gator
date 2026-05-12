package main

import (
	"github.com/tmuelli/blog-aggregator/internal/config"
	"github.com/tmuelli/blog-aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmdMap map[string]func(s *state, cmd command) error
}