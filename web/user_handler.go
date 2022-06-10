package web

import (
	"encoding/json"
	"net/http"
	"time"

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

		errpw := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(usr.Password))
		if errpw != nil && errpw == bcrypt.ErrMismatchedHashAndPassword {
			http.Error(w, errpw.Error(), http.StatusBadRequest)
			return
		}

		// generate access token
		var accessTokenDuration = time.Duration(1) * time.Minute
		accessToken, accessTokenExpiresAt, err := generateToken(u, accessTokenDuration)
		if err != nil {
			return
		}

		// generate refresh token
		var refreshTokenDuration = time.Duration(7) * 24 * time.Hour
		refreshToken, refreshTokenExpiresAt, err := generateToken(u, refreshTokenDuration)
		if err != nil {
			return
		}

		tokens := gotodo.Token{
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  time.Unix(accessTokenExpiresAt, 0),
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: time.Unix(refreshTokenExpiresAt, 0),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokens)
	}
}
