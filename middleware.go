package main

import (
	"context"

	"github.com/kristiansroberts/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfgPtr.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
