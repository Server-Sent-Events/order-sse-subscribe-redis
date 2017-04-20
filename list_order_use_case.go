package main

import (
	"encoding/json"
	"net/http"
)

func listOrder(w http.ResponseWriter, r *http.Request) {

	logicNumber := r.Header.Get("logic_number")
	//gRedisClient.Del(logicNumber)

	keys, err := gRedisClient.HKeys(logicNumber).Result()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(keys) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	orders := make([]Order, 0)

	for _, v := range keys {

		j, err := gRedisClient.HGet(logicNumber, v).Result()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		var order Order
		if err := json.Unmarshal([]byte(j), &order); err == nil {
			orders = append(orders, order)
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	respondWithJSON(w, http.StatusOK, orders)
}
