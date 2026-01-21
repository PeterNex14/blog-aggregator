package main

import (
	"context"
	"fmt"

	"github.com/PeterNex14/blog_aggregator/internal/database"
)


func middlewareLoggedIn(handler func(s *state, cmd command, user database.User)error) func(*state, command) error {
	return func (s *state, cmd command) error  {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Error occured when retrieve user data: %w", err)
		}
		return handler(s, cmd, user)
	}
}