package handlers

import (
	"github.com/zavyalov-den/go-musthave-diploma/internal/storage"
	"net/http"
)

func Balance(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// получение текущего баланса счета баллов лояльности пользователя.
	}
}
