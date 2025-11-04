package main

import (
	"fmt"

	"github.com/asqit/open-ping/helpers"
	"github.com/asqit/open-ping/web"
)

func main() {
	helpers.Print_header()
	init_db()
	defer close_db()
	go web.StartDashboard(db)

	cfg, err := read_config()
	if err != nil {
		panic(err)
	}
	run_monitor(cfg)
	fmt.Printf("Watching: %d", len(cfg.Targets))
}
