package handlers

import (
	"context"
	"errors"
	"github.com/zavyalov-den/go-musthave-diploma/internal/entities"
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
