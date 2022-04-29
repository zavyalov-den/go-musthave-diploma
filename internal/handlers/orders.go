package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zavyalov-den/go-musthave-diploma/internal/config"
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

		userID, err := getUserID(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		err = db.CreateOrder(ctx, orderNum, userID)
		if err != nil {
			if errors.Is(err, entities.ErrUserConflict) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			} else if errors.Is(err, entities.ErrEntryExists) {
				http.Error(w, err.Error(), http.StatusOK)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = RequestAccrual(ctx, db, orderNum, userID)
		if err != nil {
			fmt.Println(err)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			//return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func OrdersGet(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		w.Header().Set("Content-Type", "application/json")

		userID, err := getUserID(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

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

func RequestAccrual(ctx context.Context, db *storage.Storage, orderNum string, userID int) error {
	url := fmt.Sprintf("%s/api/orders/%s", config.GetConfig().AccrualSystemAddress, orderNum)

	fmt.Println(url)

	r, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	//if r.StatusCode != http.StatusOK {
	fmt.Println(r.StatusCode)
	//	return err
	//}

	var order entities.AccrualOrder

	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer r.Body.Close()

	err = json.Unmarshal(data, &order)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = db.UpdateOrder(ctx, order)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if order.Accrual > 0 {
		err = db.UpdateUserBalance(ctx, userID, order.Accrual)
	}

	return nil
}
