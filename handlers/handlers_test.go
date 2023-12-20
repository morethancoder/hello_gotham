package handlers_test

import (
	"bytes"
	"morethancoder/hello_gotham/components"
	"morethancoder/hello_gotham/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRenderViewHandlers(t *testing.T) {
    routePath := "/"
    view := components.HomeView("Hello")
    reqHx, err := http.NewRequest("GET", routePath, nil)
    if err != nil {t.Fatal(err)}
    reqHx.Header.Set("Hx-Request", "true")

    req, err := http.NewRequest("GET", routePath, nil)
    if err != nil {t.Fatal(err)}

    rrHx := httptest.NewRecorder()
    rr := httptest.NewRecorder()

    handlers.RenderView(rrHx, reqHx, view, routePath)
    if status := rrHx.Code; status != http.StatusOK {
        t.Errorf("Render with hx returned wrong status: got %v want %v", status, http.StatusOK)
    }

    handlers.RenderView(rr, req, view, routePath)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Render without hx returned wrong status: got %v want %v", status, http.StatusOK)
    }
}

func TestTimePostHandler(t *testing.T) {
    jsonData := `{"time": "00:00:00"}`
    req, err := http.NewRequest("Post","/", bytes.NewBufferString(jsonData))
    if err != nil {t.Fatal(err)}

    req.Header.Set("Content-Type","application/json")
    rr := httptest.NewRecorder()

    err = components.Time("00:00:00").Render(req.Context(), rr)
    if err != nil {t.Fatal(err)}
}


