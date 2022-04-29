package handlers

import (
	"context"
	"encoding/json"
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
		//- `402` — на счету недостаточно средств;
		//- `422` — неверный номер заказа;
		orders, err := db.GetOrders(ctx, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var orderExists bool

		for _, o := range orders {
			if o.Number == withdrawal.Order {
				orderExists = true
				break
			}
		}

		if !orderExists {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		balance, err := db.GetUserBalance(ctx, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if balance.Current-withdrawal.Sum < 0 {
			http.Error(w, err.Error(), http.StatusPaymentRequired)
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
	}
}
