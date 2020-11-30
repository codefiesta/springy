package http

import (
	"fmt"
	"go.springy.io/app"
	"html/template"
	"log"
	"net/http"
)

var (
	hub *app.Hub
)

func init() {

	hub = app.NewHub()

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
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	tmpl.Execute(w, nil)
}

func Start() {
	env := app.Env()
	port := fmt.Sprintf(":%d", env.Server.Port)
	log.Println("Starting http server [", port, "]")
	log.Fatal(http.ListenAndServe(port, nil))
}
