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
