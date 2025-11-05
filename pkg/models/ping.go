package models

import "database/sql"

type Target struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Ping struct {
	ID        int64
	Target    sql.NullString
	Status    sql.NullInt64
	Success   sql.NullBool
	Latency   sql.NullInt64
	Timestamp sql.NullTime
}

type PingView struct {
	ID        int64
	Target    string
	Status    int64
	Success   bool
	Latency   int64
	Timestamp string
}

type DayStatus struct {
	Date   string
	Uptime float64
}

type TargetView struct {
	Name          string
	URL           string
	CurrentUptime float64
	AverageUptime float64
	AvgLatency    float64
	HistoryCount  int
	DailyHistory  []DayStatus
}
