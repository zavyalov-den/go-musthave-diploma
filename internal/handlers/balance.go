package handlers

import (
	"context"
	"fmt"
	"github.com/zavyalov-den/go-musthave-diploma/internal/storage"
	"net/http"
)

// BalanceGet получение текущего количества баллов лояльности пользователя.
func BalanceGet(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		w.Header().Set("Content-Type", "application/json")

		userID, err := getUserID(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		fmt.Println(userID)

	}
}

//func OrdersGet(db *storage.Storage) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		ctx, cancel := context.WithCancel(r.Context())
//		defer cancel()
//
//		w.Header().Set("Content-Type", "application/json")
