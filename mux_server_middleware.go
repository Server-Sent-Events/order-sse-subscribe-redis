package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// Middleware usado para validar o path subscribe
// se possui o UUID do canal compartilhado por outro cliente
func ensureSubscribe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var channelUUID string
		var terminal string

		if channelUUID = r.FormValue("channel_uuid"); len(strings.TrimSpace(channelUUID)) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("deu ruim 1")
			log.Printf("<< %s %s %v", r.Method, r.URL.Path, time.Since(start))
			return
		}

		if terminal = r.FormValue("terminal_uuid"); len(strings.TrimSpace(terminal)) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("deu ruim 2")
			log.Printf("<< %s %s %v", r.Method, r.URL.Path, time.Since(start))
			return
		}

		if _, ok := gChannels[channelUUID]; !ok {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("deu ruim 3")
			log.Printf("<< %s %s %v", r.Method, r.URL.Path, time.Since(start))
			return
		}

		//fazer outras verificacoes para garantir EC, etc
		channel := gChannels[channelUUID]
		channel.Terminals[terminal] = &Terminal{
			UUID: terminal,
			Sub:  gRedisClient.Subscribe(channelUUID),
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
