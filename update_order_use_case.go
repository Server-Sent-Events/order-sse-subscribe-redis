package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// funcao para validar valor pago x valor de items
func paymentValue(order *Order) error {
	return nil
}

func findOrderRedis(number string, uuid string) (*Order, error) {
	var redisOrder Order
	jsonOrder, _ := orderClient.redisClient.HGet(number, uuid).Result()
	err := json.Unmarshal([]byte(jsonOrder), &redisOrder)
	if err != nil {
		return nil, err
	}
	return &redisOrder, nil
}

func unmarshalOrder(r io.Reader) (*Order, error) {
	var order Order
	b, _ := ioutil.ReadAll(r)
	err := json.Unmarshal(b, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func validateOrder(order *Order, merchant string) error {
	if order.MerchantID != merchant {
		return errors.New("merchant_id can not be changed")
	}
	return nil
}

func storeDB(order Order) error {
	orderClient.redisClient.HDel(order.LogicNumber, order.UUID)
	return nil
}

func storeRedis(order Order) error {
	jsonOrder, err := json.Marshal(order)
	if err != nil {
		return err
	}

	r := orderClient.redisClient.HSet(order.LogicNumber, order.UUID, string(jsonOrder))
	if r.Err() != nil {
		return r.Err()
	}
	return nil
}

func updateOrder(w http.ResponseWriter, r *http.Request) {

	var orderUUID string
	number := r.FormValue("logic_number")
	merchantUUID := r.FormValue("merchant_id")

	if orderUUID = mux.Vars(r)["order_id"]; len(strings.TrimSpace(orderUUID)) == 0 {
		respondWithError(w, http.StatusBadRequest, "OrderId not found")
		return
	}

	// order em memoria no redis
	redisOrder, err := findOrderRedis(number, orderUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	//nova ordem
	order, err := unmarshalOrder(r.Body)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	err = validateOrder(redisOrder, merchantUUID)
	if err != nil {
		respondWithError(w, http.StatusConflict, err.Error())
		return
	}

	order.MerchantID = redisOrder.MerchantID
	order.LogicNumber = redisOrder.LogicNumber
	order.UUID = redisOrder.UUID
	order.UpdatedAt = time.Now()

	err = paymentValue(order)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if order.Status == CLOSED {
		err = storeDB(*order)
		if err != nil {
			respondWithError(w, http.StatusServiceUnavailable, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	err = storeRedis(*order)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, order)
}
