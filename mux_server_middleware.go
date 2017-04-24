package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func ensureBaseOrder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		if logicNumber := r.FormValue("logic_number"); len(strings.TrimSpace(logicNumber)) == 0 {
			respondWithError(w, http.StatusBadRequest, "logic_number not found")
			log.Printf("<< %s %s %v", r.Method, r.URL.Path, time.Since(start))
			return
		}

		if merchantID := r.FormValue("merchant_id"); len(strings.TrimSpace(merchantID)) == 0 {
			respondWithError(w, http.StatusBadRequest, "merchant_id not found")
			log.Printf("<< %s %s %v", r.Method, r.URL.Path, time.Since(start))
			return
		}
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)

		log.Printf("<< %s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

//r.Header.Get("logic_number")

// Middleware usado para validar o path subscribe
// se possui o UUID do canal compartilhado por outro cliente
func ensureSubscribe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var channelUUID string
		var terminal string

		if channelUUID = r.FormValue("channel_uuid"); len(strings.TrimSpace(channelUUID)) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("<< %s %s %v", r.Method, r.URL.Path, time.Since(start))
			return
		}

		if terminal = r.FormValue("terminal_uuid"); len(strings.TrimSpace(terminal)) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("<< %s %s %v", r.Method, r.URL.Path, time.Since(start))
			return
		}

		if _, ok := orderClient.channels[channelUUID]; !ok {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("<< %s %s %v", r.Method, r.URL.Path, time.Since(start))
			return
		}

		//fazer outras verificacoes para garantir EC, etc
		channel := orderClient.channels[channelUUID]
		channel.Terminals[terminal] = &Terminal{
			UUID: terminal,
			Sub:  orderClient.redisClient.Subscribe(channelUUID),
		}

		next.ServeHTTP(w, r)
		log.Printf("<< %s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// Middleware para validar a criacao de um canal
// 1) nesse momento nao iremos validar o profile manager
// 2) podemos gravar quem é o dono do canal
// 3) podemos criar com regras de ec
func ensureCreateChannel(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var channelUUID string

		if channelUUID = r.Header.Get("channel_uuid"); len(strings.TrimSpace(channelUUID)) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			next.ServeHTTP(w, r)
		}

		log.Printf("<< %s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}
