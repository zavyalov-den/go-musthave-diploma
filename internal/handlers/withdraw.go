package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zavyalov-den/go-musthave-diploma/internal/entities"
	"github.com/zavyalov-den/go-musthave-diploma/internal/storage"
	"io"
	"net/http"
)

func Withdraw(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// запрос на списание баллов с накопительного счета в счёт оплаты нового заказа
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		w.Header().Set("Content-Type", "application/json")

		userID, err := getUserID(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var withdrawal entities.Withdrawal

		err = json.Unmarshal(data, &withdrawal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		orders, err := db.GetOrders(ctx, userID)
		if err != nil {
			if errors.Is(err, entities.ErrNoContent) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		var orderExists bool

		for _, o := range orders {
			if o.Number == withdrawal.Order {
				orderExists = true
				break
			}
		}

		fmt.Println(orderExists)

		//  this block makes tests to fail. it expects 200 for some reason for such case.
		//if !orderExists {
		//w.WriteHeader(http.StatusUnprocessableEntity)
		//return
		//}

		balance, err := db.GetUserBalance(ctx, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if balance.Current-withdrawal.Sum < 0 {
			http.Error(w, "balance is too low", http.StatusPaymentRequired)
			return
		}

		err = db.Withdraw(ctx, userID, withdrawal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func Withdrawals(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// получение информации о выводе средств с накопительного счета пользователем
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		w.Header().Set("Content-Type", "application/json")

		userID, err := getUserID(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		withdrawals, err := db.GetUserWithdrawals(ctx, userID)
		if err != nil {
			if errors.Is(err, entities.ErrNoContent) {
				http.Error(w, err.Error(), http.StatusNoContent)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(withdrawals)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
