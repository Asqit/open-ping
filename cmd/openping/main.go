package main

import (
	"log"

	"github.com/asqit/open-ping/internal/config"
	"github.com/asqit/open-ping/internal/monitor"
	"github.com/asqit/open-ping/internal/server"
	"github.com/asqit/open-ping/internal/storage"
)

func main() {
	store, err := storage.NewSQLite("openping.db")
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	srv := server.New(store.DB())
	go srv.Start()

	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatal(err)
	}

	mon := monitor.New(store)
	mon.Run(cfg)
}
