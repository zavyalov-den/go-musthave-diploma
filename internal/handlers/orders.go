package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/zavyalov-den/go-musthave-diploma/internal/entities"
	"github.com/zavyalov-den/go-musthave-diploma/internal/service"
	"github.com/zavyalov-den/go-musthave-diploma/internal/storage"
	"io"
	"net/http"
	"strconv"
)

// доступные статусы обработки расчётов:
//
//- `NEW` — заказ загружен в систему, но не попал в обработку;
//- `PROCESSING` — вознаграждение за заказ рассчитывается;
//- `INVALID` — система расчёта вознаграждений отказала в расчёте;
//- `PROCESSED` — данные по заказу проверены и информация о расчёте успешно получена.
//

func OrdersPost(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// загрузка пользователем заказа для рассчета
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		orderNum := string(data)

		orderNumInt, err := strconv.Atoi(orderNum)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !service.Valid(orderNumInt) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		userID := int(ctx.Value("userID").(float64))

		err = db.CreateOrder(ctx, orderNum, userID)
		if err != nil {
			if errors.Is(err, entities.ErrUserConflict) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func OrdersGet(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		w.Header().Set("Content-Type", "application/json")

		userID := int(ctx.Value("userID").(float64))

		orders, err := db.GetOrders(ctx, userID)
		if err != nil {
			if errors.Is(err, entities.ErrNoContent) {
				http.Error(w, err.Error(), http.StatusNoContent)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		respData, err := json.Marshal(orders)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(respData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
