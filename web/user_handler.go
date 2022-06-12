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

		userid := r.Context().Value("user").(uuid.UUID)
		if userid != id {
			http.Error(w, "Invalid request", http.StatusForbidden)
			return
		}

		u, err := s.store.User(id)
		// er if the given id is not in the database
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		resp := gotodo.UserResponse{
			ID:                u.ID,
			UserName:          u.UserName,
			FirstName:         u.FirstName,
			LastName:          u.LastName,
			Email:             u.Email,
			CreatedAt:         u.CreatedAt,
			PasswordChangedAt: u.PasswordChangedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
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
		u.ID = uuid.New()

		hashedPass, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		u.Password = string(hashedPass)

		if er := s.store.CreateUser(u); er != nil {
			http.Error(w, er.Error(), http.StatusInternalServerError)
			return
		}

		resp := gotodo.UserResponse{
			ID:                u.ID,
			UserName:          u.UserName,
			FirstName:         u.FirstName,
			LastName:          u.LastName,
			Email:             u.Email,
			CreatedAt:         u.CreatedAt,
			PasswordChangedAt: u.PasswordChangedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func (s *Server) updateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.Parse(getField(r, 0))

		// only allow logged in user to update logged in user
		userid := r.Context().Value("user").(uuid.UUID)
		if userid != id {
			http.Error(w, "Invalid request", http.StatusForbidden)
			return
		}

		checkUser, _ := s.store.User(id)

		u := &gotodo.User{}
		if err := json.NewDecoder(r.Body).Decode(u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashedPass, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		u.Password = string(hashedPass)

		if checkUser.Password != u.Password {
			u.PasswordChangedAt = time.Now()
		}

		if er := s.store.UpdateUser(id, u); er != nil {
			http.Error(w, er.Error(), http.StatusInternalServerError)
			return
		}

		resp := gotodo.UserResponse{
			ID:                u.ID,
			UserName:          u.UserName,
			FirstName:         u.FirstName,
			LastName:          u.LastName,
			Email:             u.Email,
			CreatedAt:         u.CreatedAt,
			PasswordChangedAt: u.PasswordChangedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func (s *Server) deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.Parse(getField(r, 0))

		// only allow logged in user to delete logged in user
		userid := r.Context().Value("user").(uuid.UUID)
		if userid != id {
			http.Error(w, "Invalid request", http.StatusForbidden)
			return
		}

		u, _ := s.store.User(id)

		if er := s.store.DeleteUser(id); er != nil {
			http.Error(w, er.Error(), http.StatusBadRequest)
			return
		}

		resp := gotodo.UserResponse{
			ID:                u.ID,
			UserName:          u.UserName,
			FirstName:         u.FirstName,
			LastName:          u.LastName,
			Email:             u.Email,
			CreatedAt:         u.CreatedAt,
			PasswordChangedAt: u.PasswordChangedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func (s *Server) userLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr := &gotodo.LoginUserRequest{}
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
			http.Error(w, errpw.Error(), http.StatusUnauthorized)
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

		user := gotodo.UserResponse{
			ID:                u.ID,
			UserName:          u.UserName,
			FirstName:         u.FirstName,
			LastName:          u.LastName,
			Email:             u.Email,
			CreatedAt:         u.CreatedAt,
			PasswordChangedAt: u.PasswordChangedAt,
		}

		resp := &gotodo.LoginUserResponse{
			Token: tokens,
			User:  user,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
