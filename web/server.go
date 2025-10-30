package web

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/asqit/open-ping/helpers"
	"github.com/asqit/open-ping/types"
	"github.com/asqit/open-ping/view"
	"github.com/asqit/open-ping/view/layout"
	"github.com/asqit/open-ping/view/share"
)

func get_pings(db *sql.DB) ([]types.PingView, error) {
	rows, err := db.Query("SELECT id, target, status, success, latency, timestamp FROM pings ORDER BY timestamp DESC LIMIT 100")
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return pings, nil
}

func Start_dashboard(db *sql.DB) {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	pings, err := get_pings(db)
	if err != nil {
		panic(err)
	}

	c := layout.Base(view.Index(share.Dashboard(pings)))
	http.Handle("/", templ.Handler(c))
	http.Handle("/api/pings/paginate", PaginatePings(db))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
