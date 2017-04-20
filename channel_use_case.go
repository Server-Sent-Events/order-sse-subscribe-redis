package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func postMsg(w http.ResponseWriter, r *http.Request) {

	var channelUUID string
	if channelUUID = mux.Vars(r)["channel_id"]; len(strings.TrimSpace(channelUUID)) == 0 {
		respondWithError(w, http.StatusBadRequest, "chanel_id not found")
		return
	}

	orderClient.redisClient.Publish(channelUUID, "hello")

	w.WriteHeader(http.StatusOK)

}
