package handlers

import (
	"github.com/zavyalov-den/go-musthave-diploma/internal/storage"
	"net/http"
)

func Withdraw(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// запрос на списание баллов с накопительного счета в счёт оплаты нового заказа
	}
}

func Withdrawals(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// получение информации о выводе средств с накопительного счета пользователем
	}
}
