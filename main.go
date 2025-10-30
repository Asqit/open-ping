package main

import "github.com/asqit/open-ping/web"

// Basic flow:
// 1. Read config
// 2. Every interval
// 	- ping each target
// 	- measure HTTP response status & latency
// 	- log result to sqlite
// 3. Print to STDOUT stream

func main() {
	init_db()
	defer close_db()
	go web.Start_dashboard(db)

	cfg, err := read_config()
	if err != nil {
		panic(err)
	}

	run_monitor(cfg)
}
