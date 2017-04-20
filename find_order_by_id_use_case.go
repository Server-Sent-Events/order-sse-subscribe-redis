package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func findOrderByID(w http.ResponseWriter, r *http.Request) {

	var orderUUID string
	if orderUUID = mux.Vars(r)["order_id"]; len(strings.TrimSpace(orderUUID)) == 0 {
		respondWithError(w, http.StatusBadRequest, "OrderId not found")
		return
	}

	logicNumber := r.Header.Get("logic_number")

	jsonOrder, err := orderClient.redisClient.HGet(logicNumber, orderUUID).Result()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithString(w, http.StatusOK, jsonOrder)
}
