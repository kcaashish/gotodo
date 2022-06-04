package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

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
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk.UserID)
		r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func Refresh() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			return
		}

		if time.Unix(tk.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		expiresAt := time.Now().Add(time.Minute * 5)

		tk.ExpiresAt = expiresAt.Unix()

		tokenRef := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)

		tokenString, errtk := tokenRef.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
		if errtk != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(errtk)
		}

		var resp = map[string]string{"refresh_token": string(tokenString)}
		json.NewEncoder(w).Encode(resp)
	})
}
