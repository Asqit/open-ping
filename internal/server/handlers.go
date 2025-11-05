package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/asqit/open-ping/pkg/models"
	"github.com/asqit/open-ping/web/templates/components"
)

func toString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func toInt64(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return 0
}

func toBool(nb sql.NullBool) bool {
	if nb.Valid {
		return nb.Bool
	}
	return false
}

type PaginatePingsResponse struct {
	Data []models.PingView `json:"data"`
}

func (s *Server) PaginatePings(w http.ResponseWriter, r *http.Request) {
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

	rows, err := s.db.Query("SELECT * FROM pings ORDER BY id DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pings []models.PingView
	for rows.Next() {
		var p models.Ping
		if err := rows.Scan(&p.ID, &p.Target, &p.Status, &p.Success, &p.Latency, &p.Timestamp); err != nil {
			log.Println("scan:", err)
			continue
		}

		pings = append(pings, models.PingView{
			ID:        p.ID,
			Target:    toString(p.Target),
			Status:    toInt64(p.Status),
			Success:   toBool(p.Success),
			Latency:   toInt64(p.Latency),
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

func (s *Server) GetDistinctTargets(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query("SELECT DISTINCT target FROM pings")
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

func (s *Server) GetTargetsHTML(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(`
		SELECT target,
		       target_website,
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

	var targets []models.TargetView
	for rows.Next() {
		var t models.TargetView
		var total, successCount int
		if err := rows.Scan(&t.Name, &t.URL, &total, &successCount, &t.AvgLatency); err != nil {
			log.Println("scan:", err)
			continue
		}
		t.CurrentUptime = (float64(successCount) / float64(total)) * 100
		t.AverageUptime = t.CurrentUptime
		t.HistoryCount = total
		t.DailyHistory = s.getDailyHistory(t.Name)
		targets = append(targets, t)
	}

	c := components.TargetList(targets)
	templ.Handler(c).ServeHTTP(w, r)
}

func (s *Server) getDailyHistory(target string) []models.DayStatus {
	rows, err := s.db.Query(`
		SELECT DATE(timestamp) as day,
		       COUNT(*) as total,
		       SUM(success) as success_count
		FROM pings
		WHERE target = ?
		  AND timestamp >= DATE('now', '-90 days')
		GROUP BY day
		ORDER BY day ASC;
	`, target)
	if err != nil {
		return []models.DayStatus{}
	}
	defer rows.Close()

	var history []models.DayStatus
	for rows.Next() {
		var day string
		var total, successCount int
		if err := rows.Scan(&day, &total, &successCount); err != nil {
			continue
		}
		uptime := (float64(successCount) / float64(total)) * 100
		history = append(history, models.DayStatus{Date: day, Uptime: uptime})
	}

	return history
}
