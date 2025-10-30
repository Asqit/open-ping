package web

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
)

type Ping struct {
	ID        int64
	Target    sql.NullString
	Status    sql.NullInt64
	Success   sql.NullBool
	Latency   sql.NullInt64
	Timestamp sql.NullTime
}

func Start_dashboard(db *sql.DB) {
	tmpl := template.Must(template.ParseFiles("web/templates/dashboard.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, target, status, success, latency, timestamp FROM pings ORDER BY timestamp DESC LIMIT 100")
		if err != nil {
			http.Error(w, "query error", http.StatusInternalServerError)
			log.Println("query:", err)
			return
		}
		defer rows.Close()

		var pings []Ping
		for rows.Next() {
			var p Ping
			if err := rows.Scan(&p.ID, &p.Target, &p.Status, &p.Success, &p.Latency, &p.Timestamp); err != nil {
				log.Println("scan:", err)
				continue
			}
			pings = append(pings, p)
		}
		if err := rows.Err(); err != nil {
			log.Println("rows err:", err)
		}

		if err := tmpl.Execute(w, pings); err != nil {
			log.Println("tmpl execute:", err)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
