package main

import (
	"github.com/gorilla/mux"
	"github.com/zavyalov-den/go-musthave-diploma/internal/handlers"
	"github.com/zavyalov-den/go-musthave-diploma/internal/middlewares"
	"github.com/zavyalov-den/go-musthave-diploma/internal/storage"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	// todo
	db := &storage.Storage{}

	r.Use(middlewares.GzipHandle)

	r.HandleFunc("/api/user/register", handlers.Register(db))
	r.HandleFunc("/api/user/login", handlers.Login(db))
	r.HandleFunc("/api/user/orders", nil)
	r.HandleFunc("/api/user/balance/withdrawals", nil)
	r.HandleFunc("/api/user/balance/withdraw", nil)
	r.HandleFunc("/api/user/balance", nil)

	http.Handle("/", r)

}
