package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kristiansroberts/gator/internal/database"
)

func loginHandler(state *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: login <username>")
	}
	username := cmd.args[0]

	_, err := state.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("error getting user: %v", err)
	}

	err = state.cfgPtr.SetUser(username)
	if err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}
	fmt.Printf("Logged in as %s\n", username)
	return nil
}

func registerHandler(state *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: register <username>")
	}
	username := cmd.args[0]

	user, err := state.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	fmt.Printf("Registered new user: %+v\n", user)

	err = state.cfgPtr.SetUser(username)
	if err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}
	fmt.Printf("Logged in as %s\n", username)
	return nil
}

func resetHandler(state *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: reset")
	}
	err := state.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting users: %v", err)
	}
	fmt.Println("All users deleted")
	return nil
}

func usersHandler(state *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: users")
	}

	users, err := state.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users: %v", err)
	}
	for _, user := range users {
		if user.Name == state.cfgPtr.CurrentUserName {
			fmt.Printf("* %+v (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %+v\n", user.Name)
	}
	return nil
}

func aggHandler(state *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: agg")
	}
	url := "https://www.wagslane.dev/index.xml"

	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}

	fmt.Printf("%+v\n", feed)
	return nil
}

func addFeedHandler(state *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}
	name := cmd.args[0]
	url := cmd.args[1]

	feed, err := state.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed: %v", err)
	}

	_, err = state.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed follow: %v", err)
	}
	fmt.Printf("Added new feed: %+v\n", feed)
	return nil
}

func feedsHandler(state *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: feeds")
	}

	feeds, err := state.db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %v", err)
	}

	for _, feed := range feeds {
		user, _ := state.db.GetUserByID(context.Background(), feed.UserID)
		fmt.Printf("* %s\n", feed.Name)
		fmt.Printf("* %s\n", feed.Url)
		fmt.Printf("* %s\n", user.Name)
	}
	return nil
}

func followHandler(state *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: follow <url>")
	}
	url := cmd.args[0]

	feed, err := state.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error getting feed: %v", err)
	}

	feedFollow, err := state.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed follow: %v", err)
	}
	fmt.Printf("Feed followed: %+v\nCreated by: %+v\n", feedFollow, user.Name)
	return nil
}

func followingHandler(state *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: following")
	}

	follows, err := state.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting feed follows: %v", err)
	}
	for _, follow := range follows {
		fmt.Printf("* %s\n", follow.FeedName)
	}

	return nil
}

func unfollowHandler(state *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: unfollow <url>")
	}
	url := cmd.args[0]

	feed, err := state.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error getting feed: %v", err)
	}

	err = state.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error deleting feed follow: %v", err)
	}
	fmt.Printf("Feed unfollowed: %+v\n", feed.Name)
	return nil
}
