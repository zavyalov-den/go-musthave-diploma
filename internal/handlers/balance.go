package handlers

import (
	"context"
	"encoding/json"
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

		balance, err := db.GetUserBalance(ctx, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		balanceData, err := json.Marshal(balance)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(balanceData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

//func OrdersGet(db *storage.Storage) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		ctx, cancel := context.WithCancel(r.Context())
//		defer cancel()
//
//		w.Header().Set("Content-Type", "application/json")
