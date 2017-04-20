package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

func subscribeChannel(w http.ResponseWriter, r *http.Request) {

	var channelUUID string
	number := r.Header.Get("logic_number")
	merchantUUID := r.Header.Get("merchant_id")
	pin := r.Header.Get("pin")

	if merchantUUID == "" {
		merchantUUID = r.FormValue("merchant_id")
	}
	if number == "" {
		number = r.FormValue("logic_number")
	}
	if pin == "" {
		pin = r.FormValue("pin")
	}

	if channelUUID = mux.Vars(r)["channel_id"]; len(strings.TrimSpace(channelUUID)) == 0 {
		respondWithError(w, http.StatusBadRequest, "chanel_id not found")
		return
	}

	f, ok := w.(http.Flusher)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Streaming unsupported!")
		return
	}

	ch, err := orderClient.getChannel(channelUUID, merchantUUID, pin)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	terminal := ch.openChannel(number, ch.UUID)

	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		ch.exitChannel("number")
		log.Println("HTTP connection just closed.")
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	var receive error
	var msg *redis.Message
	for receive == nil {

		msg, receive = terminal.Sub.ReceiveMessage()

		fmt.Fprintf(w, "data: Message: %s\n\n", msg.Payload)

		f.Flush()
	}
}
