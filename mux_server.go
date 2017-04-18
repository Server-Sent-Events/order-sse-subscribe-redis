package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Inicializa o servidor HTTP
func startMux() {

	router := mux.NewRouter().StrictSlash(true)

	router.Handle("/api/v3/order", ensureCreateOrder(http.HandlerFunc(createOrder))).Methods("POST")
	router.Handle("/api/v3/order/{order_id}", ensureCreateOrder(http.HandlerFunc(updateOrder))).Methods("PUT")

	router.Handle("/subscribe", ensureSubscribe(http.HandlerFunc(subscribeChannel))).Methods("GET")
	router.Handle("/api/v1/channel", ensureCreateChannel(http.HandlerFunc(createChannel))).Methods("POST")
	router.Handle("/", http.HandlerFunc(MainPageHandler))
	router.Handle("/msg", http.HandlerFunc(postMsg)).Methods("POST")

	http.ListenAndServe(":3000", router)
}
