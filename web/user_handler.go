package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/kcaashish/gotodo"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) getUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.Parse(getField(r, 0))

		u, err := s.store.User(id)
		// er if the given id is not in the database
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u)
	}
}

func (s *Server) getUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uu, err := s.store.Users()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(uu)
	}
}

func (s *Server) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &gotodo.User{}
		if err := json.NewDecoder(r.Body).Decode(u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashedPass, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		u.Password = string(hashedPass)

		if er := s.store.CreateUser(u); er != nil {
			http.Error(w, er.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u)
	}
}

func (s *Server) updateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.Parse(getField(r, 0))
		u := &gotodo.User{}
		if err := json.NewDecoder(r.Body).Decode(u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashedPass, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		u.Password = string(hashedPass)

		if er := s.store.UpdateUser(id, u); er != nil {
			http.Error(w, er.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u)
	}
}

func (s *Server) deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.Parse(getField(r, 0))

		if er := s.store.DeleteUser(id); er != nil {
			http.Error(w, er.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) userLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr := &gotodo.User{}
		if er := json.NewDecoder(r.Body).Decode(usr); er != nil {
			http.Error(w, er.Error(), http.StatusBadRequest)
			return
		}

		u, err := s.store.FindUser(usr.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		expiresAt := time.Now().Add(time.Minute * 5)

		errpw := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(usr.Password))
		if errpw != nil && errpw == bcrypt.ErrMismatchedHashAndPassword {
			http.Error(w, errpw.Error(), http.StatusBadRequest)
			return
		}

		tk := &gotodo.Token{
			UserID:   u.ID,
			UserName: u.UserName,
			Email:    u.Email,
			StandardClaims: &jwt.StandardClaims{
				ExpiresAt: expiresAt.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)

		tokenString, errtk := token.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
		if errtk != nil {
			fmt.Println(errtk)
		}

		var resp = map[string]string{"access_token": string(tokenString)}
		json.NewEncoder(w).Encode(resp)
	}
}
