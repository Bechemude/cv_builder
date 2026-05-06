package main

import (
	"cvbuilder/bot"
	"cvbuilder/config"
	"cvbuilder/db"
	"cvbuilder/external"
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

	ex := external.Init(cfg)

	repos, err := repos.Init(db)
	if err != nil {
		panic(err)
	}

	s, err := services.Init(repos, ex, cfg)
	if err != nil {
		panic(err)
	}

	h, err := handlers.Init(s, ex)

	bot := bot.Init(cfg)

	bot.RegisterHandlers(h)
	bot.Run()
}
