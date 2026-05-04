package main

import (
	"cvbuilder/bot"
	"cvbuilder/config"
	"cvbuilder/db"
	"cvbuilder/handlers"
	"cvbuilder/repos"
	"cvbuilder/services"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		panic(err)
	}

	db, err := db.Init(cfg)
	if err != nil {
		panic(err)
	}

	repos, err := repos.Init(db)
	if err != nil {
		panic(err)
	}

	s, err := services.Init(repos)
	if err != nil {
		panic(err)
	}

	h, err := handlers.Init(s)

	bot := bot.Init(cfg)

	bot.RegisterHandlers(h)
	bot.Run()
}
