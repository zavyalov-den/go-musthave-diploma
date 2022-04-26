package handlers

import (
	"context"
	"github.com/zavyalov-den/go-musthave-diploma/internal/service"
	"github.com/zavyalov-den/go-musthave-diploma/internal/storage"
	"io"
	"net/http"
	"strconv"
)

func OrdersPost(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// загрузка пользователем заказа для рассчета
		_, cancel := context.WithCancel(r.Context())
		defer cancel()

		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		orderNum, err := strconv.Atoi(string(data))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if !service.IsValid(orderNum) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}

		//db.CreateOrder(ctx, orderNum, userID)

		// todo: order accepted. create order

		w.WriteHeader(http.StatusAccepted)

	}
}
