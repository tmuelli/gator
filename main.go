package main

import (
	"log"
	"errors"
	"fmt"

	"github.com/tmuelli/blog-aggregator/internal/config"
)

func main() {
	cfg := config.Read()
	if cfg == nil {
		log.Fatal(errors.New("No configuration could be read"))
	}

	cfg.SetUser("tmmue")


	cfg = config.Read()
	if cfg == nil {
		log.Fatal(errors.New("No configuration could be read"))
	}

	fmt.Println("User name:", cfg.CurrentUserName)
	fmt.Println("Database url:", cfg.DbUrl)
}