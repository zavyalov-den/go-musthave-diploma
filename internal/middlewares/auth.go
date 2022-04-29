package middlewares

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/zavyalov-den/go-musthave-diploma/internal/config"
	"github.com/zavyalov-den/go-musthave-diploma/internal/entities"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// todo
		jwtCookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		tokenString := jwtCookie.Value

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unxepected signing method: %v", token.Header["alg"])
			}

			return []byte(config.GetConfig().SessionKey), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["userID"]
			ctx := context.WithValue(r.Context(), entities.ContextUserID, userID)

			r = r.WithContext(ctx)
		} else {
			http.Error(w, "token is not valid", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
