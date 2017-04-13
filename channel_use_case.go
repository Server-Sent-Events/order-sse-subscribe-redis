package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
)

func postMsg(w http.ResponseWriter, r *http.Request) {

	_ = gRedisClient.Publish("0012a2ed-c500-4b8f-83e7-c7da351d839c", "hello").Err()

}

func createChannel(w http.ResponseWriter, r *http.Request) {

	channelUUID := r.Header.Get("channel_uuid")

	if _, ok := gChannels[channelUUID]; !ok {

		ch := &Channel{
			UUID:      channelUUID,
			Terminals: make(map[string]*Terminal),
		}

		gChannels[channelUUID] = ch

		w.WriteHeader(http.StatusCreated)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func subscribeChannel(w http.ResponseWriter, r *http.Request) {

	// channel := gChannels[r.Header.Get("channel_uuid")]
	// terminalUUID := r.Header.Get("terminal_uuid")

	channel := gChannels[r.FormValue("channel_uuid")]
	terminalUUID := r.FormValue("terminal_uuid")

	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		delete(channel.Terminals, terminalUUID)

		log.Println("HTTP connection just closed.")
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	t := channel.Terminals[terminalUUID]

	var receive error
	var msg *redis.Message
	for receive == nil {

		msg, receive = t.Sub.ReceiveMessage()

		fmt.Fprintf(w, "data: Message: %s\n\n", msg.Payload)

		f.Flush()
	}
}
