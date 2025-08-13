package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
)

type contextKey string

var userContextKey = contextKey("user")

func setUserInContext(r *http.Request, claims jwt.MapClaims) context.Context {
	return context.WithValue(r.Context(), userContextKey, claims)
}

func GetUserFromContext(r *http.Request) (jwt.MapClaims, bool) {
	claims, ok := r.Context().Value(userContextKey).(jwt.MapClaims)
	return claims, ok
}

func GetSession(r *http.Request) *sessions.Session {
	var store = sessions.NewCookieStore([]byte(os.Getenv("SECRET")))
	session, err := store.Get(r, "session")
	if err!=nil{
		fmt.Println(err)
	}
	return session
}
