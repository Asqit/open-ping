package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func ping_url(name, url string) {
	start := time.Now()
	resp, err := http.Get(url)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("[%s] %s → DOWN (%v)\n", time.Now().Format("15:04:05"), name, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("[%s] %s → %d OK, %v\n", time.Now().Format("15:04:05"), name, resp.StatusCode, elapsed)
	save_log(name, resp.StatusCode, resp.StatusCode < 300, int(elapsed))
}

func run_monitor(cfg *Config) {
	interval, _ := time.ParseDuration(cfg.Interval)
	for {
		var wg sync.WaitGroup
		for _, target := range cfg.Targets {
			wg.Add(1)
			go func(t Target) {
				defer wg.Done()
				ping_url(t.Name, t.Url)
			}(target)
		}

		wg.Wait()
		time.Sleep(interval)
	}
}
