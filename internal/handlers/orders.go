package handlers

import (
	"context"
	"fmt"
	"github.com/zavyalov-den/go-musthave-diploma/internal/service"
	"github.com/zavyalov-den/go-musthave-diploma/internal/storage"
	"io"
	"net/http"
	"strconv"
)

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

		orderNum, err := strconv.Atoi(string(data))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !service.Valid(orderNum) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		// db.CreateOrder(ctx, orderNum, userID)
		login := ctx.Value("login")

		fmt.Println(login)

		// return 409 conflict if order exist with another user_id

		// todo: order accepted. create order

		w.WriteHeader(http.StatusAccepted)

	}
}
