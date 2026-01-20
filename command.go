package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PeterNex14/blog_aggregator/internal/config"
)

type state struct {
	cfg 	*config.Config
}

type command struct {
	Name	string
	Args	[]string
}

type commands struct {
	registeredCommands 		map[string]func(*state, command) error
}


func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("username is required")
	}

	err := s.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Println("User has been set")
	return nil
}


func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]

	if !ok {
		return errors.New("There is no such argument")
	}

	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error)  {
	c.registeredCommands[name] = f
}


func cleanInput(text string) []string {
	texts := strings.Fields(strings.ToLower(text))
	return texts
}