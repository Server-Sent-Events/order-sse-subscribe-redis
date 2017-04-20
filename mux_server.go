package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Inicializa o servidor HTTP
func startMux() {

	router := mux.NewRouter().StrictSlash(true)

	router.Handle("/api/v3/order", ensureCreateOrder(http.HandlerFunc(createOrder))).Methods("POST")
	router.Handle("/api/v3/order", ensureCreateOrder(http.HandlerFunc(listOrder))).Methods("GET")
	router.Handle("/api/v3/order/{order_id}", ensureCreateOrder(http.HandlerFunc(findOrderByID))).Methods("GET")
	router.Handle("/api/v3/order/{order_id}", ensureCreateOrder(http.HandlerFunc(updateOrder))).Methods("PUT")

	// sse sockets
	router.Handle("/api/v3/order/shareOrder", ensureCreateOrder(http.HandlerFunc(shareOrder))).Methods("PUT")
	router.Handle("/api/v3/order/subscribe/{channel_id}", ensureCreateOrder(http.HandlerFunc(subscribeChannel))).Methods("GET")

	router.Handle("/", http.HandlerFunc(MainPageHandler))
	router.Handle("/msg", http.HandlerFunc(postMsg)).Methods("POST")

	http.ListenAndServe(":3000", router)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithString(w http.ResponseWriter, code int, payload string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(payload))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
