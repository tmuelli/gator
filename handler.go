package main

import (
	"errors"
	"fmt"
	"context"
	"time"
	"log"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tmuelli/blog-aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Invalid arguments for login command.")
	}

	if s == nil {
		return errors.New("Internal error - Invalid state")
	}

	// check if user exists
	if _, err := s.db.GetUserByName(context.Background(), sql.NullString{String: cmd.args[0], Valid: true}); err != nil {
		log.Fatal(err)
	}

	s.cfg.SetUser(cmd.args[0])
	fmt.Println("User was set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Invalid arguments for register command.")
	}

	if s == nil {
		return errors.New("Internal error - Invalid state")
	}

	// check if user exists
	_, err := s.db.GetUserByName(context.Background(), sql.NullString{String: cmd.args[0], Valid: true})
	if  err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// user not exist yet
		} else {
			log.Fatal(err)
		}
	}

	if err == nil {
		// user already exists
		log.Fatal(errors.New("User is already registered!"))
	}

	createdUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: 			uuid.New(),
		CreatedAt:		sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:		sql.NullTime{Time: time.Now(), Valid: true},
		Name:			sql.NullString{String: cmd.args[0], Valid: true},
	})
	if err != nil {
		log.Fatal(err)
	}

	// set user to config
	s.cfg.SetUser(createdUser.Name.String)
	fmt.Printf("User %s was successfully created at %v with id %v\n", createdUser.Name.String, createdUser.CreatedAt.Time, createdUser.ID)
	return nil
}

func handlerReset(s *state, cmd command) error {
	// reset users
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	// get all users
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range users {
		if u.Name.String == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", u.Name.String)
		} else {
			fmt.Println("*", u.Name.String)
		}
	}

	return nil
}

func handlerAggregate(s *state, cmd command) error {
	// call and fetch rss feed
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", feed)

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("Invalid arguments for addfeed command.")
	}

	if s == nil {
		return errors.New("Internal error - Invalid state")
	}

	// create feed
	createFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:		sql.NullString{String: cmd.args[0], Valid: true},
		Url:		sql.NullString{String: cmd.args[1], Valid: true},
		UserID:		user.ID,
	})
	if err != nil {
		log.Fatal(err)
	}

	// create new feed follow
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID:		user.ID,
		FeedID:		createFeed.ID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", createFeed)
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	// get all feeds
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range feeds {
		fmt.Printf("%s from %s with url %s\n", f.Feedname.String, f.Username.String, f.Feedurl.String)
	}

	return nil
}

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("Invalid arguments for follow command.")
	}

	if s == nil {
		return errors.New("Internal error - Invalid state")
	}

	// get feed by url
	feed, err := s.db.GetFeedByUrl(context.Background(), sql.NullString{String: cmd.args[0], Valid: true})
	if err != nil {
		log.Fatal(err)
	}

	// create new feed follow
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID:		user.ID,
		FeedID:		feed.ID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User %s follows feed %s now.\n", feedFollow.Username.String, feedFollow.Feedname.String)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	// get feeds the user is following
	feedFollows, err := s.db.GetFeedFollowsByUser(context.Background(), user.ID)
	if err != nil {
		log.Fatal(err)
	}

	for _, ff := range feedFollows {
		fmt.Println(ff.Feedname)
	}
	
	return nil
}

func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("Invalid arguments for unfollow command.")
	}

	if s == nil {
		return errors.New("Internal error - Invalid state")
	}

	// get feed by url
	feed, err := s.db.GetFeedByUrl(context.Background(), sql.NullString{String: cmd.args[0], Valid: true})
	if err != nil {
		log.Fatal(err)
	}

	err = s.db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID:		user.ID,
		FeedID:		feed.ID,
	})
	if err != nil {
		log.Fatal(err)
	}

	return nil
}