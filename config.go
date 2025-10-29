package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Target struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type EncodedConfig struct {
	Interval string   `json:"interval"`
	Targets  []Target `json:"targets"`
}

type Config struct {
	Interval int
	Targets  []Target
}

func Parse_duration_ms(input string) (int64, error) {
	s := strings.TrimSpace(strings.ToLower(input))
	if s == "" {
		return 0, errors.New("empty input")
	}

	if matched, _ := regexp.MatchString(`^\d{1,2}:\d{2}(:\d{2})?$`, s); matched {
		parts := strings.Split(s, ":")
		var h, m, sec int64
		if v, err := strconv.ParseInt(parts[0], 10, 64); err == nil {
			h = v
		} else {
			return 0, err
		}
		if v, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
			m = v
		} else {
			return 0, err
		}
		if len(parts) == 3 {
			if v, err := strconv.ParseInt(parts[2], 10, 64); err == nil {
				sec = v
			} else {
				return 0, err
			}
		}
		total := h*3600 + m*60 + sec
		return total * 1000, nil
	}

	re := regexp.MustCompile(`^([\d.]+)\s*(ms|s|sec|m|min|h|d)?$`)
	m := re.FindStringSubmatch(s)
	if m == nil {
		return 0, errors.New("unrecognized format")
	}

	valStr := m[1]
	unit := m[2]
	if unit == "" {
		unit = "s"
	}

	f, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return 0, err
	}

	var ms float64
	switch unit {
	case "ms":
		ms = f
	case "s", "sec":
		ms = f * 1000
	case "m", "min":
		ms = f * 60 * 1000
	case "h":
		ms = f * 3600 * 1000
	case "d":
		ms = f * 24 * 3600 * 1000
	default:
		return 0, errors.New("unknown unit, accepting: ms, s, sec, m, min, h, d")
	}

	return int64(ms + 0.5), nil // round to nearest ms
}

func Read_config() (*Config, error) {
	contents, error := os.ReadFile("config.json")
	if error != nil {
		return nil, error
	}

	var payload EncodedConfig
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

	var result Config
	result.Targets = payload.Targets
	interval, error := Parse_duration_ms(payload.Interval)
	if error != nil {
		return nil, error
	}

	result.Interval = int(interval)

	return &result, nil
}
