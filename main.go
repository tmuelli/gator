package main

import (
	"log"
	"errors"
	"os"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/tmuelli/blog-aggregator/internal/config"
	"github.com/tmuelli/blog-aggregator/internal/database"
)

func main() {
	cmds := commands{
		cmdMap: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAggregate)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollowFeed))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollowFeed))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	// read configuration
	cfg := config.Read()
	if cfg == nil {
		log.Fatal(errors.New("No configuration could be read"))
	}

	// open database
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}

	// create db queries
	dbQueries := database.New(db)

	s := state{
		db: dbQueries,
		cfg: cfg,
	}

	if len(os.Args) < 2 {
		log.Fatal(errors.New("To few arguments"))
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}