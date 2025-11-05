package monitor

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/asqit/open-ping/internal/config"
	"github.com/asqit/open-ping/internal/storage"
	"github.com/asqit/open-ping/pkg/models"
)

type Monitor struct {
	storage storage.Storage
}

func New(s storage.Storage) *Monitor {
	return &Monitor{storage: s}
}

func (m *Monitor) Ping(name, url string) {
	start := time.Now()
	resp, err := http.Get(url)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("[%s] %s → DOWN (%v)\n", time.Now().Format("15:04:05"), name, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("[%s] %s → %d OK, %v\n", time.Now().Format("15:04:05"), name, resp.StatusCode, elapsed)
	m.storage.SavePing(name, url, resp.StatusCode, int(elapsed.Milliseconds()), resp.StatusCode < 300)
}

func (m *Monitor) Run(cfg *config.Config) {
	interval, _ := time.ParseDuration(cfg.Interval)
	fmt.Printf("Watching: %d\n", len(cfg.Targets))
	
	for {
		var wg sync.WaitGroup
		for _, target := range cfg.Targets {
			wg.Add(1)
			go func(t models.Target) {
				defer wg.Done()
				m.Ping(t.Name, t.URL)
			}(target)
		}

		wg.Wait()
		time.Sleep(interval)
	}
}
