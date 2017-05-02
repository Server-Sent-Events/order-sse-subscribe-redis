package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Inicializa o servidor HTTP
func startMux() {
	router := mux.NewRouter().StrictSlash(true)
	apiV3(router)
	http.ListenAndServe(":8080", router)
}

func apiV3(router *mux.Router) {
	router.Handle("/", http.HandlerFunc(MainPageHandler))

	apiV3 := router.PathPrefix("/api/v3").Subrouter()
	apiV3.Handle("/order", ensureBaseOrder(http.HandlerFunc(createOrder))).Methods("POST")
	apiV3.Handle("/order", ensureBaseOrder(http.HandlerFunc(listOrder))).Methods("GET")
	apiV3.Handle("/order/{order_id}", ensureBaseOrder(http.HandlerFunc(findOrderByID))).Methods("GET")
	apiV3.Handle("/order/{order_id}", ensureBaseOrder(http.HandlerFunc(updateOrder))).Methods("PUT")

	// sse sockets
	apiV3.Handle("/order/{order_id}/share", ensureBaseOrder(http.HandlerFunc(shareOrder))).Methods("PUT")
	apiV3.Handle("/subscribe/{channel_id}", ensureBaseOrder(http.HandlerFunc(subscribeChannel))).Methods("GET")

	apiV3.Handle("/msg/{channel_id}", http.HandlerFunc(postMsg)).Methods("POST")
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
