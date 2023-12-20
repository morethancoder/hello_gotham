package main

import (
	"morethancoder/hello_gotham/db"
	"morethancoder/hello_gotham/handlers"
	"morethancoder/hello_gotham/middleware"
	"morethancoder/hello_gotham/models"
	"log"
	"net/http"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)



func main() {
    dbClient, err := db.Init("./configs/.env")
    if err != nil {log.Fatal(err)}  
    sessionManager := scs.New()
    sessionManager.Store = mysqlstore.New(dbClient)

    var g models.GlobalValuesInstance
    err = models.InitGlobalValuesTable(dbClient)
    if err != nil {log.Fatal(err)}
    err = g.Create(dbClient)
    if err != nil {log.Fatal(err)}
    err = models.InitSessionsTable(dbClient)
    if err != nil {log.Fatal(err)}
    err = models.InitTodosTable(dbClient)
    if err != nil {log.Fatal()}

    r := chi.NewRouter()
    r.Use(middleware.DbClientMiddleware(dbClient))
    r.Use(middleware.SessionManagerMiddleware(sessionManager))

    fs := http.FileServer(http.Dir("static"))
    r.Handle("/static/*", http.StripPrefix("/static/", fs))
    r.Get("/", handlers.HomeGetHandler)
    r.Get("/about", handlers.AboutGetHandler)
    r.Get("/counter", handlers.CounterGetHandler)
    r.Get("/todos", handlers.TodosGetHandler)
    r.Post("/time", handlers.TimePostHandler)
    r.Post("/counter", handlers.CounterPostHandler)
    r.Post("/todos", handlers.TodosPostHandler)
    log.Println("Listening on Port 3000")
    err = http.ListenAndServe(":3000", sessionManager.LoadAndSave(r))
    if err != nil { log.Fatal(err) }
}
