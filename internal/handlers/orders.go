package handlers

import (
	"github.com/zavyalov-den/go-musthave-diploma/internal/storage"
	"net/http"
)

func Orders(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// загрузка пользователем заказа для рассчета
	}
}
