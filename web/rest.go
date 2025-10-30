package web

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/asqit/open-ping/helpers"
	"github.com/asqit/open-ping/types"
)

type PaginatePingsResponse struct {
	Data []types.PingView `json:"data"`
}

func PaginatePings(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		offsetStr := r.URL.Query().Get("offset")
		limitStr := r.URL.Query().Get("limit")

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}

		rows, err := db.Query("SELECT * FROM pings ORDER BY id DESC LIMIT ? OFFSET ?", limit, offset)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Internal Server Error"))
			return
		}

		defer rows.Close()

		var pings []types.PingView
		for rows.Next() {
			var p types.Ping
			if err := rows.Scan(&p.ID, &p.Target, &p.Status, &p.Success, &p.Latency, &p.Timestamp); err != nil {
				log.Println("scan:", err)
				continue
			}

			pings = append(pings, types.PingView{
				ID:        p.ID,
				Target:    helpers.ToString(p.Target),
				Status:    helpers.ToInt64(p.Status),
				Success:   helpers.ToBool(p.Success),
				Latency:   helpers.ToInt64(p.Latency),
				Timestamp: p.Timestamp.Time.Format("2006-01-02 15:04:05"),
			})
		}

		if err := rows.Err(); err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Internal Server Error"))
			return
		}

		bytes, err := json.Marshal(PaginatePingsResponse{Data: pings})
		if err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	}
}
