package types

import "database/sql"

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

type TargetView struct {
	Name          string  // target hostname or label
	URL           string  // same as Name for now
	CurrentUptime float64 // percentage of successful pings
	AverageUptime float64 // (placeholder) same as CurrentUptime until 90-day avg calc is added
	AvgLatency    float64 // average latency (ms)
	HistoryCount  int     // number of pings in total
}
