package main

import "fmt"

// Basic flow:
// 1. Read config
// 2. Every interval
// 	- ping each target
// 	- measure HTTP response status & latency
// 	- log result to sqlite
// 3. Print to STDOUT stream

func main() {
	cfg, err := Read_config()
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg.Interval)
}
