package main

import (
	"log"
	"newscomuni/bot"
	"newscomuni/db"
	"newscomuni/rss"
	"sync"

	envr "github.com/Netflix/go-env"
)

type Env struct {
	BOT_TOKEN string `env:"BOT_TOKEN"`
}

func main() {
	var env Env
	// Obtain environment variables.
	_, err := envr.UnmarshalFromEnviron(&env)
	if err != nil {
		log.Fatal(err)
	}

	// Setup WaitGroup
	var wg sync.WaitGroup

	// Setup database
	database := db.New()
	defer database.Close()

	// Setup telegram bot
	bot := bot.New(env.BOT_TOKEN, database)
	bot.Start(&wg)

	// Setup rss parser
	rss := rss.New()

	// Scan feeds
	bot.Scanner(&wg, database, rss)

	// Wait for goroutines to return.
	wg.Wait()
}
