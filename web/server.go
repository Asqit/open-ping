package web

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/asqit/open-ping/view"
	"github.com/asqit/open-ping/view/layout"
	"github.com/asqit/open-ping/view/share"
)

func StartDashboard(db *sql.DB) {
	fmt.Print("Starting web dashboard.........")

	// Static assets
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Main dashboard page
	dashboard := share.Dashboard()
	index := view.Index(dashboard)
	page := layout.Base(index)
	http.Handle("/", templ.Handler(page))

	// APIs
	http.Handle("/api/pings/paginate", PaginatePings(db))
	http.Handle("/api/targets", GetDistinctTargets(db))
	http.Handle("/api/targets/html", GetTargetsHTML(db)) // used by HTMX

	fmt.Print("[OK]\n")
	fmt.Println("Local Machine: http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
