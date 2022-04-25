package handlers

import (
	"context"
	"encoding/json"
	"github.com/zavyalov-den/go-musthave-diploma/internal/config"
	"github.com/zavyalov-den/go-musthave-diploma/internal/entities"
	"github.com/zavyalov-den/go-musthave-diploma/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func Login(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// authenticate an existing user
		data, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		defer r.Body.Close()

		var reqData entities.Credentials

		err = json.Unmarshal(data, &reqData)
		if err != nil {
			http.Error(w, "request data is invalid: "+err.Error(), http.StatusUnauthorized)
			return
		}

		user, err := db.GetUser(context.Background(), reqData.Login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqData.Password))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"username": reqData.Login})
		sessionKey := config.GetConfig().SessionKey

		tokenString, err := token.SignedString([]byte(sessionKey))
		if err != nil {
			http.Error(w, "failed to sign a token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{Name: "jwt", Value: tokenString})

		w.WriteHeader(http.StatusOK)
	}
}
