package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/PeterNex14/blog_aggregator/internal/config"
	"github.com/PeterNex14/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
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

	dbURL := data.DBUrl
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)
	s.db = dbQueries

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUserList)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handleAddFeed))
	cmds.register("feeds", handleFeeds)
	cmds.register("follow", middlewareLoggedIn(handleFollow))
	cmds.register("following", middlewareLoggedIn(handleFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handleUnfollow))
	

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