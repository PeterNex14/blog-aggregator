package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/PeterNex14/blog_aggregator/internal/config"
	"github.com/PeterNex14/blog_aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type state struct {
	db 		*database.Queries
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

	user, err := s.db.GetUser(
		context.Background(),
		cmd.Args[0],
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("user doesn't exists")
			os.Exit(1)
		}
	}

	
	if err := s.cfg.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Println("User has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {

	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name: cmd.Args[0],
		},
	)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			fmt.Println("Postgres error:", pgErr.Message)
			fmt.Println("user already exist")
			os.Exit(1)
		} else {
			return err
		}
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Println("User successfuly created")
	fmt.Printf("uuid : %v\n", user.ID)
	fmt.Printf("created_at : %v\n", user.CreatedAt)
	fmt.Printf("updated_at : %v\n", user.UpdatedAt)
	fmt.Printf("name : %s\n", user.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.RemoveUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete users: %w", err)
	} 

	fmt.Println("Data Reset Successfuly")
	return nil
}

func handlerUserList(s *state, cmd command) error {
	list, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to retrieve users data: %v", err)
	}

	if len(list) == 0 {
		fmt.Println("There is no existing data")
		os.Exit(1)
	}

	for _, value := range list{
		if value.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", value.Name)
		} else {
			fmt.Printf("* %s\n", value.Name)
		}
	}

	return nil

}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		fmt.Println("Not enough arguments!")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("go run . agg <duration>")
		os.Exit(1)
	}
	duration, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Error parsing time argument: %w", err)
	}

	fmt.Printf("Collecting feeds every %s\n", duration)

	ticker := time.NewTicker(duration)
	for ; ; <- ticker.C {
		if err := scrapeFeeds(s, cmd); err != nil {
			return err
		}
	}

}

func handleAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		fmt.Println("Not enough arguments!")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("go run . addfeed <name> <url>")
		os.Exit(1)
	}

	connect, err := s.db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name: cmd.Args[0],
			Url: cmd.Args[1],
			UserID: user.ID,
		},
	)

	

	if err != nil {
		return fmt.Errorf("Error Creating Feed, %v", err)
	}

	feed_follow, err := s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID: user.ID,
			FeedID: connect.ID,
		},
	)

	if err != nil {
		return fmt.Errorf("Error following newly created feed: %w", err)
	}

	fmt.Println("Feeds saved successfuly")
	fmt.Printf("Feed with title %v successfuly followed\n", feed_follow.FeedName)
	fmt.Println()
	fmt.Printf("id: %v\n", connect.ID)
	fmt.Printf("created_at: %v\n", connect.CreatedAt)
	fmt.Printf("updated_at: %v\n", connect.UpdatedAt)
	fmt.Printf("name: %v\n", connect.Name)
	fmt.Printf("url: %v\n", connect.Url)
	fmt.Printf("user_id: %v\n", connect.UserID)

	return nil
}

func handleFeeds(s *state, cmd command) error {
	data, err := s.db.GetFeedsUser(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to retrieve data: %v", err)
	}

	for _, item := range data {
		fmt.Printf("Title: %v\n", item.Name)
		fmt.Printf("URL: %v\n", item.Url)
		fmt.Printf("Name: %v\n", item.UserName)
		fmt.Println()
	}

	return nil
}

func handleFollow(s *state, cmd command, user database.User) error {

	feed, err := s.db.GetFeedsByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Error occured when retrieving feed: %w", err)
	}


	feed_follow, err := s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID: user.ID,
			FeedID: feed.ID,
		},
	)

	fmt.Println("Title: ", feed_follow.FeedName)
	fmt.Println("Current User: ", feed_follow.UserName)

	return nil 
}

func handleFollowing(s *state, cmd command, user database.User) error {

	data, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Error occured when retrieving following data: %w", err)
	}

	for _, item := range data {
		fmt.Println(item.FeedName)
	}

	return nil
}

func handleUnfollow(s *state, cmd command, user database.User) error {

	feed, err := s.db.GetFeedsByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Error occured when retrieve feed data:")
	}

	err = s.db.DeleteFollowByUserAndFeedId(
		context.Background(),
		database.DeleteFollowByUserAndFeedIdParams{
			UserID: user.ID,
			FeedID: feed.ID,
		},
	)

	if err != nil {
		return fmt.Errorf("Failed to Delete Data: %w", err)
	}

	fmt.Println("Feeds successfully unfollowed")

	return nil
}

func scrapeFeeds(s *state, cmd command) error {
	var ctx = context.Background()
	next_feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("Failed to fetch next feed: %w", err)
	}

	err = s.db.MarkFeedFetched(ctx, next_feed.ID)
	if err != nil {
		return fmt.Errorf("Failed to mark feed: %w", err)
	}

	fetch_feed, err := fetchFeed(ctx, next_feed.Url)
	if err != nil {
		return fmt.Errorf("Failed to fetch current feed: %w", err)
	}

	for _, value := range fetch_feed.Channel.Item {
		fmt.Println(value.Title)
	}

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

