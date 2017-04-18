package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// funcao para validar valor pago x valor de items
func paymentValue() {

}

func updateOrder(w http.ResponseWriter, r *http.Request) {

	var orderUUID string
	if orderUUID = mux.Vars(r)["order_id"]; len(strings.TrimSpace(orderUUID)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	//TODO: VALIDAR QUE O VALOR TOTAL DE ITEMS NAO PODE SER MAIOR
	//QUE O VALOR PAGO, MAIS TAMBEM DEVE SER TRAVADO DENTRO DO TERMINAL
	//POIS ELA FAZ O PAGAMENTO
	order.UUID = orderUUID

	if order.Status == CLOSED {
		//TODO: GRAVA BANCO DE DADOS
		//NAO VEJO A NECESSIDADE DE GERAR OS UUID DOS ITENS

		//SE SUCESSO REMOVE DO REDIS
		gRedisClient.HDel(order.LogicNumber, order.UUID)

		w.WriteHeader(http.StatusOK)
		return
	}

	jsonOrder, err := json.Marshal(order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	gRedisClient.HSet(order.LogicNumber, order.UUID, string(jsonOrder))

	w.WriteHeader(http.StatusOK)
}
