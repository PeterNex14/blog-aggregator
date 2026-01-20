package main

import (
	"fmt"
	"os"

	"github.com/PeterNex14/blog_aggregator/internal/config"
)

func main() {
	data, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	s := &state{
		cfg: data,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	input := os.Args

	if len(input) < 2 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}

	cmd := command{
		Name: input[1],
		Args: input[2:],
	}

	if err := cmds.run(s, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
}