package main

import (
	"github.com/gorilla/mux"
	"github.com/zavyalov-den/go-musthave-diploma/internal/config"
	"github.com/zavyalov-den/go-musthave-diploma/internal/handlers"
	"github.com/zavyalov-den/go-musthave-diploma/internal/middlewares"
	"github.com/zavyalov-den/go-musthave-diploma/internal/storage"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	db := storage.NewStorage()
	db.InitDB()

	r.Use(middlewares.GzipHandle)

	r.HandleFunc("/api/user/register", handlers.Register(db)).Methods(http.MethodPost)
	r.HandleFunc("/api/user/login", handlers.Login(db)).Methods(http.MethodPost)

	//r.Use(middlewares.Auth) todo

	r.HandleFunc("/api/user/orders", handlers.OrdersPost(db)).Methods(http.MethodPost)
	r.HandleFunc("/api/user/balance/withdrawals", handlers.Withdrawals(db)).Methods(http.MethodGet)
	r.HandleFunc("/api/user/balance/withdraw", handlers.Withdraw(db)).Methods(http.MethodPost)
	r.HandleFunc("/api/user/balance", handlers.Balance(db)).Methods(http.MethodGet)

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(config.GetConfig().RunAddress, r))

}
