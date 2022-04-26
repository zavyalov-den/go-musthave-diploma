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

	//r.Use(middlewares.GzipHandle)

	r.HandleFunc("/api/user/register", handlers.Register(db)).Methods(http.MethodPost)
	r.HandleFunc("/api/user/login", handlers.Login(db)).Methods(http.MethodPost)

	br := r.PathPrefix("/api/user/balance").Subrouter()
	or := r.PathPrefix("/api/user/orders").Subrouter()

	br.Use(middlewares.AuthMiddleware)
	or.Use(middlewares.AuthMiddleware)

	or.HandleFunc("", handlers.OrdersPost(db)).Methods(http.MethodPost)

	br.HandleFunc("/withdrawals", handlers.Withdrawals(db)).Methods(http.MethodGet)
	br.HandleFunc("/withdraw", handlers.Withdraw(db)).Methods(http.MethodPost)
	br.HandleFunc("/balance", handlers.Balance(db)).Methods(http.MethodGet)

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(config.GetConfig().RunAddress, r))

}
