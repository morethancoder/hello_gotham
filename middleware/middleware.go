package middleware

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type key int

const (
    DbClientKey key = iota
    SessionManagerKey key = iota
)

func DbClientMiddleware(dbClient *sql.DB) func(http.Handler) http.Handler {
    return func(h http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx := context.WithValue(r.Context(), DbClientKey, dbClient)
            h.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func SessionManagerMiddleware(sessionManager *scs.SessionManager) func(http.Handler) http.Handler {
    return func(h http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx := context.WithValue(r.Context(), SessionManagerKey , sessionManager)
            h.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
