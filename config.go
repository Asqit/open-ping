package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

type Target struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Config struct {
	Interval string   `json:"interval"`
	Targets  []Target `json:"targets"`
}

func read_config() (*Config, error) {
	contents, error := os.ReadFile("config.json")
	if error != nil {
		return nil, error
	}

	var payload Config
	error = json.Unmarshal(contents, &payload)
	if error != nil {
		return nil, error
	}

	for i, t := range payload.Targets {
		_, parseErr := url.ParseRequestURI(t.Url)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid URL in target %d (%s): %v", i, t.Name, parseErr)
		}
	}

	return &payload, nil
}
