package http

import (
	"fmt"
	"go.springy.io/internal/mongo"
	"go.springy.io/internal/ws"
	"go.springy.io/pkg/util"
	"html/template"
	"log"
	"net/http"
)

func init() {

	initRoutes()

	// Run the db in a new goroutine
	go mongo.Run()

	// Run the hub in a new goroutine
	go ws.Run()
}

// Initialize the http routes
func initRoutes() {
	http.HandleFunc("/", indexRoute)
	http.HandleFunc("/ws", ws.Upgrade)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
}

/// Returns the index.html
func indexRoute(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	tmpl.Execute(w, nil)
}

func Start() {
	env := util.Env()
	port := fmt.Sprintf(":%d", env.Server.Port)
	log.Printf("🌱 [Starting http server %s] 🌱", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
