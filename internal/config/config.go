package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/asqit/open-ping/pkg/models"
)

type Config struct {
	Interval string          `json:"interval"`
	Targets  []models.Target `json:"targets"`
}

func Load(path string) (*Config, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(contents, &cfg); err != nil {
		return nil, err
	}

	for i, t := range cfg.Targets {
		if _, err := url.ParseRequestURI(t.URL); err != nil {
			return nil, fmt.Errorf("invalid URL in target %d (%s): %v", i, t.Name, err)
		}
	}

	return &cfg, nil
}
