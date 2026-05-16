package main

import (
	"database/sql"
	"context"
	"log"

	"github.com/tmuelli/blog-aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(s *state, cmd command) error {
	return func(s *state, cmd command) error {
		// get current user
		currentUser, err := s.db.GetUserByName(context.Background(), sql.NullString{String: s.cfg.CurrentUserName, Valid: true})
		if  err != nil {
			log.Fatal(err)
		}

		return handler(s, cmd, currentUser)
	}
}