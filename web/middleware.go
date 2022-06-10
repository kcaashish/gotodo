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

		// extract token
		tokenString := r.Header.Get("Authorization")

		// set up a claim
		claims := &gotodo.Claims{}

		// verify token
		token, er := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
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

		ctx := context.WithValue(r.Context(), "user", claims.UserID)
		r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func Refresh() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		claims := &gotodo.Claims{}

		token, er := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
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

		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		expiresAt := time.Now().Add(time.Minute * 5)

		claims.ExpiresAt = expiresAt.Unix()

		tokenRef := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, errtk := tokenRef.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
		if errtk != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(errtk)
		}

		var resp = map[string]string{"refresh_token": string(tokenString)}
		json.NewEncoder(w).Encode(resp)
	})
}

func generateToken(u gotodo.User, period time.Duration) (string, int64, error) {
	now := time.Now()
	claims := &gotodo.Claims{
		UserID:   u.ID,
		UserName: u.UserName,
		Email:    u.Email,
		StandardClaims: &jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(period).Unix(),
		},
	}

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, errtk := tokenWithClaims.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
	if errtk != nil {
		fmt.Println(errtk)
		return "", 0, errtk
	}
	return token, claims.ExpiresAt, nil
}
