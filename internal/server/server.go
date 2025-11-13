package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/asqit/open-ping/web/templates/components"
	"github.com/asqit/open-ping/web/templates/layout"
	"github.com/asqit/open-ping/web/templates/pages"
)

type Server struct {
	db *sql.DB
}

func New(db *sql.DB) *Server {
	return &Server{db: db}
}

func printHeader() {
	fmt.Println(`
\==================================/
░█▀█░█▀█░█▀▀░█▀█░░░░░█▀█░▀█▀░█▀█░█▀▀
░█░█░█▀▀░█▀▀░█░█░▄▄▄░█▀▀░░█░░█░█░█░█
░▀▀▀░▀░░░▀▀▀░▀░▀░░░░░▀░░░▀▀▀░▀░▀░▀▀▀
/==================================\`)
}

func (s *Server) Start() {
	printHeader()
	fmt.Print("Starting web dashboard.........")

	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	dashboard := components.Dashboard()
	index := pages.Index(dashboard)
	page := layout.Base(index)
	http.Handle("/", templ.Handler(page))

	http.HandleFunc("/api/pings/paginate", s.PaginatePings)
	http.HandleFunc("/api/targets", s.GetDistinctTargets)
	http.HandleFunc("/api/targets/html", s.GetTargetsHTML)
	http.HandleFunc("/api/targets/{target}/export", s.PerTargetExport)
	http.HandleFunc("/api/export/all", s.ExportAll)

	fmt.Print("[OK]\n")
	fmt.Println("Local Machine: http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
