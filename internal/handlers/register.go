package handlers

import (
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/zavyalov-den/go-musthave-diploma/internal/config"
	"github.com/zavyalov-den/go-musthave-diploma/internal/entities"
	"github.com/zavyalov-den/go-musthave-diploma/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
)

func Register(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		// register new user
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "request data is invalid", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var cred entities.Credentials

		err = json.Unmarshal(data, &cred)
		if err != nil {
			http.Error(w, "failed to parse request data", http.StatusInternalServerError)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(cred.Password), 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userID, err := db.Register(ctx, &entities.Credentials{
			Login:    cred.Login,
			Password: string(hash),
		})
		if err != nil {
			http.Error(w, "failed to create a user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		//
		resp, err := json.Marshal(cred)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"userID": userID, "login": cred.Login})
		sessionKey := config.GetConfig().SessionKey

		tokenString, err := token.SignedString([]byte(sessionKey))
		if err != nil {
			http.Error(w, "failed to sign a token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{Name: "jwt", Value: tokenString})

		// todo: check docs for the expected response
		_, err = w.Write(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
