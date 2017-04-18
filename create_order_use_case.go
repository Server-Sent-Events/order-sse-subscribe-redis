package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// Primeiro passo seria identificar se uma ordem está no redis (cache) ou não
// 1- ordem sem UUID nunca foi enviada para backend
// 2- ordem sem UUID com status não CLOSED está em um ambiente transacional
// 3- ordem sem UUID com status CLOSED, fez toda sua operação offline
// 4- é necessário validar regras do objeto (itens, valores, etc)
func createOrder(w http.ResponseWriter, r *http.Request) {

	var order Order
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(b, &order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if orderUUID := order.UUID; len(strings.TrimSpace(orderUUID)) == 0 {
		order.UUID = uuid.NewV4().String()
	}

	order.LogicNumber = r.Header.Get("logic_number")
	order.MerchantID = r.Header.Get("merchant_id")

	jsonOrder, err := json.Marshal(order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//include order to redis
	//INCLUIR MUTEX PARA ORDEM COMPARTILHADA?
	if order.Status != CLOSED {
		gRedisClient.HSet(order.LogicNumber, order.UUID, string(jsonOrder))
	} else {
		//TODO: GRAVAR DIRETO NO BD
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonOrder)
}
