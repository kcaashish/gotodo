package web

import (
	"context"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/kcaashish/gotodo"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var noAuthRequired = []string{
			"/users/login", "/users/create",
		}

		path := r.URL.Path

		for _, route := range noAuthRequired {
			if route == path {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenString := r.Header.Get("Authorization")

		tk := &gotodo.Token{}

		token, er := jwt.ParseWithClaims(tokenString, tk, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_PASSWORD")), nil
		})

		if er != nil {
			http.Error(w, er.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), "user", tk.UserID)
		r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
