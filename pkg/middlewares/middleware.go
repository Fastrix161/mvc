package middlewares

import (
	"context"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
)

type contextKey string

var userContextKey = contextKey("user")
var store = sessions.NewCookieStore([]byte(os.Getenv("SECRET")))
func setUserInContext(r *http.Request, claims jwt.MapClaims) context.Context {
	return context.WithValue(r.Context(), userContextKey, claims)
}

func GetUserFromContext(r *http.Request) (jwt.MapClaims, bool) {
	claims, ok := r.Context().Value(userContextKey).(jwt.MapClaims)
	return claims, ok
}

func GetSession(r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "session-name")
	return session
}
