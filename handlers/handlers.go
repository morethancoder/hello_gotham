package handlers

import (
	"database/sql"
	"errors"
	"morethancoder/hello_gotham/components"
	"morethancoder/hello_gotham/middleware"
	"morethancoder/hello_gotham/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"
)

func onError(w http.ResponseWriter, err error, msg string, code int) {
	if err != nil {
		http.Error(w, msg, code)
		log.Println(msg, err)
	}
}

func RenderView(w http.ResponseWriter, r *http.Request, view templ.Component, layoutPath string) {
	if r.Header.Get("Hx-Request") == "true" {
		err := view.Render(r.Context(), w)
		onError(w, err, "Internal server error", http.StatusInternalServerError)
		return
	}
	err := components.Layout(layoutPath).Render(r.Context(), w)
	onError(w, err, "Internal server error", http.StatusInternalServerError)
}

func HomeGetHandler(w http.ResponseWriter, r *http.Request) {
	RenderView(w, r, components.HomeView("hello, world!"), "/")
}
func AboutGetHandler(w http.ResponseWriter, r *http.Request) {
	RenderView(w, r, components.AboutView(), "/about")
}

func TimePostHandler(w http.ResponseWriter, r *http.Request) {
	clientTimeStr := r.FormValue("time")
	err := components.Time(clientTimeStr).Render(r.Context(), w)
	onError(w, err, "Internal server error", http.StatusInternalServerError)
}

func CounterGetHandler(w http.ResponseWriter, r *http.Request) {
	dbClient, ok := r.Context().Value(middleware.DbClientKey).(*sql.DB)
	if !ok {
		onError(w, errors.New("Couldnt get context value (dbClient)"),
			"Internal server error", http.StatusInternalServerError)
	}
	sessionManager, ok := r.Context().Value(middleware.SessionManagerKey).(*scs.SessionManager)
	if !ok {
		onError(w, errors.New("Couldnt get context value (sessionManager)"),
			"Internal server error", http.StatusInternalServerError)
	}

	userCount := sessionManager.GetInt(r.Context(), "count")
	g := models.GlobalValuesInstance{ID: 1}
	err := g.Read(dbClient)
	onError(w, err, "Internal server error", http.StatusInternalServerError)

	RenderView(w, r, components.CounterView(g.Count, userCount), "/counter")
}

func CounterPostHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	onError(w, err, "Internal server error", http.StatusInternalServerError)

	dbClient, ok := r.Context().Value(middleware.DbClientKey).(*sql.DB)
	if !ok {
		onError(w, errors.New("Couldnt get context value "),
			"Internal server error", http.StatusInternalServerError)
	}
	sessionManager, ok := r.Context().Value(middleware.SessionManagerKey).(*scs.SessionManager)
	if !ok {
		onError(w, errors.New("Couldnt get context value (sessionManager)"),
			"Internal server error", http.StatusInternalServerError)
	}

	userCount := sessionManager.GetInt(r.Context(), "count")
	g := models.GlobalValuesInstance{ID: 1}
	err = g.Read(dbClient)
	onError(w, err, "Internal server error", http.StatusInternalServerError)

	pressedValue := r.Form.Get("pressed")
	if pressedValue == "global" {
		err := g.Update(dbClient)
		onError(w, err, "Internal server error", http.StatusInternalServerError)
	} else if pressedValue == "session" {
		userCount = userCount + 1
		sessionManager.Put(r.Context(), "count", userCount)
	}

	RenderView(w, r, components.CounterView(g.Count, userCount), "/counter")
}

func TodosGetHandler(w http.ResponseWriter, r *http.Request) {
	dbClient, ok := r.Context().Value(middleware.DbClientKey).(*sql.DB)
	if !ok {
		onError(w, errors.New("Couldnt get context value "),
			"Internal server error", http.StatusInternalServerError)
	}
	sessionManager, ok := r.Context().Value(middleware.SessionManagerKey).(*scs.SessionManager)
	if !ok {
		onError(w, errors.New("Couldnt get context value (sessionManager)"),
			"Internal server error", http.StatusInternalServerError)
	}
    var data []models.TodoInstance

    if sessionManager.Exists(r.Context(), "todos") {
        sessionToken := sessionManager.Token(r.Context())
        sessionTodos, err := models.ReadAllSessionTodos(dbClient, sessionToken)
        onError(w, err, "Internal Server Error", http.StatusInternalServerError)
        data = sessionTodos

    } else {
        sessionManager.Put(r.Context(), "todos", nil)
    }
	RenderView(w, r, components.TodosView(data), "/todos")
}


func TodosPostHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    onError(w, err, "Internal server error", http.StatusInternalServerError)
	dbClient, ok := r.Context().Value(middleware.DbClientKey).(*sql.DB)
	if !ok {
		onError(w, errors.New("Couldnt get context value "),
			"Internal server error", http.StatusInternalServerError)
	}
	sessionManager, ok := r.Context().Value(middleware.SessionManagerKey).(*scs.SessionManager)
	if !ok {
		onError(w, errors.New("Couldnt get context value (sessionManager)"),
			"Internal server error", http.StatusInternalServerError)
	}
    sessionToken := sessionManager.Token(r.Context())
    reqType := r.FormValue("type")
    switch reqType {
        case "create":
            todo := models.TodoInstance{
                Title: r.FormValue("title"),
                Text: r.FormValue("text"),
                Done: false,
                Date: time.Now(),
                SessionToken: sessionToken,
            }
            err = todo.Create(dbClient)
            onError(w, err, "Internal server error", http.StatusInternalServerError)
        case "delete":
            todoID, err := strconv.Atoi(r.FormValue("id"))
            onError(w, err, "Internal server error", http.StatusInternalServerError)
            todo := models.TodoInstance{ID: todoID}
            err = todo.Delete(dbClient)
            onError(w, err, "Internal server error", http.StatusInternalServerError)
        case "open":
            todoID, err := strconv.Atoi(r.FormValue("id"))
            onError(w, err, "Internal server error", http.StatusInternalServerError)
            todo := models.TodoInstance{ID: todoID}
            err = todo.Read(dbClient)
            onError(w, err, "Internal server error", http.StatusInternalServerError)
            err = components.TodoItem(todo).Render(r.Context(), w)
            onError(w, err, "Internal server error", http.StatusInternalServerError)
            return
        case "update":
            todoID, err := strconv.Atoi(r.FormValue("id"))
            onError(w, err, "Internal server error", http.StatusInternalServerError)
            todo := models.TodoInstance{ID: todoID}
            err = todo.Read(dbClient)
            onError(w, err, "Internal server error", http.StatusInternalServerError)
            todo.Title = r.FormValue("title")
            todo.Text = r.FormValue("text")
            todo.Date = time.Now()
            err = todo.Update(dbClient)
            onError(w, err, "Internal server error", http.StatusInternalServerError)

    }
    sessionTodos, err := models.ReadAllSessionTodos(dbClient, sessionToken)
    onError(w, err, "Internal server error", http.StatusInternalServerError)
    RenderView(w, r, components.TodosView(sessionTodos), "/todos")

}
