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
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refTokenStr := r.Header.Get("Authorization")
		claims := &gotodo.Claims{}

		token, er := jwt.ParseWithClaims(refTokenStr, claims, func(t *jwt.Token) (interface{}, error) {
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

		expiresAt := time.Now().Add(time.Duration(1) * time.Minute)

		claims.IssuedAt = time.Now().Unix()
		claims.ExpiresAt = expiresAt.Unix()

		tokenRef := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, errtk := tokenRef.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
		if errtk != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(errtk)
		}

		var resp = map[string]string{
			"new_access_token": string(tokenString),
			"expires_at":       expiresAt.Local().String(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
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
