package web

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/asqit/open-ping/helpers"
	"github.com/asqit/open-ping/types"
	"github.com/asqit/open-ping/view/share"
)

// =============================
// Pagination for all pings (as-is)
// =============================

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
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

		bytes, err := json.Marshal(PaginatePingsResponse{Data: pings})
		if err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	}
}

// =============================
// Distinct Targets as JSON
// =============================

func GetDistinctTargets(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT DISTINCT target FROM pings")
		if err != nil {
			http.Error(w, "DB query failed", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var targets []string
		for rows.Next() {
			var t string
			if err := rows.Scan(&t); err != nil {
				log.Println("scan:", err)
				continue
			}
			targets = append(targets, t)
		}

		type response struct {
			Targets []string `json:"targets"`
		}

		bytes, err := json.Marshal(response{Targets: targets})
		if err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	}
}

// =============================
// HTMX Endpoint for Dashboard Grid
// =============================

func GetTargetsHTML(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT target,
			       COUNT(*) as total,
			       SUM(success) as success_count,
			       AVG(latency) as avg_latency
			FROM pings
			GROUP BY target
			ORDER BY target ASC;
		`)
		if err != nil {
			http.Error(w, "DB query failed", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var targets []types.TargetView
		for rows.Next() {
			var t types.TargetView
			var total, successCount int
			if err := rows.Scan(&t.Name, &total, &successCount, &t.AvgLatency); err != nil {
				log.Println("scan:", err)
				continue
			}
			t.URL = t.Name
			t.CurrentUptime = (float64(successCount) / float64(total)) * 100
			t.AverageUptime = t.CurrentUptime // later: compute real 90-day avg
			t.HistoryCount = total
			targets = append(targets, t)
		}

		// Render with Templ
		c := share.TargetList(targets)
		templ.Handler(c).ServeHTTP(w, r)
	})
}
