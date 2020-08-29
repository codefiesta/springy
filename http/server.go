package http

import (
	"fmt"
	"go.springy.io/app"
	"go.springy.io/util"
	"html/template"
	"log"
	"net/http"
)

var (
	tmpl *template.Template
	hub  *app.Hub
)

func init() {

	hub = app.NewHub()
	tmpl = template.Must(template.ParseFiles("web/templates/index.html"))

	initRoutes()
	// Run the hub in a new goroutine
	go hub.Run()
}

// Initialize the http routes
func initRoutes() {
	http.HandleFunc("/", indexRoute)
	http.HandleFunc("/ws", hub.Upgrade)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
}

/// Returns the index.html
func indexRoute(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func Start() {
	config := util.Config()
	port := fmt.Sprintf(":%d", config.Server.Port)
	log.Println("Starting http server [", port, "]")
	log.Fatal(http.ListenAndServe(port, nil))
}
