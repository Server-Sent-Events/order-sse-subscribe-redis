package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ServerMux struct {
	clients        map[chan string]bool
	newClients     chan chan string
	defunctClients chan chan string
	messages       chan string
}

// Middleware usado para validar o path subscribe
// se possui o UUID do canal compartilhado por outro cliente
func ensureSubscribe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if uuid := r.Header.Get("uuid"); uuid != "" {

			if _, ok := gChannels[uuid]; !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			//fazer outras verificacoes para garantir EC, etc

			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
	})
}

// This ServerMux method handles and HTTP request at the "/events/" URL.
func (s *ServerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		log.Println("HTTP connection just closed.")
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	pubsub := gRedisClient.Subscribe("canal")

	for {

		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "data: Message: %s\n\n", msg.Payload)

		f.Flush()
	}
}

// Middleware http
func ensureParamsHeaders(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		log.Printf(">> %s %s", r.Method, r.URL.Path)

		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)

		log.Printf("<< %s %s %v", r.Method, r.URL.Path, time.Since(start))
	})

}

func getTeste(w http.ResponseWriter, r *http.Request) {
	if msg := r.FormValue("msg"); msg != "" {

		err := gRedisClient.Publish("canal", msg).Err()
		if err != nil {
			panic(err)
		}
	}
}

func MainPageHandler(w http.ResponseWriter, r *http.Request) {

	// Did you know Golang's ServeMux matches only the
	// prefix of the request URL?  It's true.  Here we
	// insist the path is just "/".
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Read in the template with our SSE JavaScript code.
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatal("WTF dude, error parsing your template.")

	}

	// Render the template, writing to `w`.
	t.Execute(w, "Duder")

	// Done.
	log.Println("Finished HTTP request at ", r.URL.Path)
}

// Inicializa o servidor HTTP
func startMux() {

	// Make a new ServerMux instance
	b := &ServerMux{
		make(map[chan string]bool),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
	}

	router := mux.NewRouter().StrictSlash(true)

	router.Handle("/subscribe", ensureSubscribe(b)).Methods("GET")
	router.Handle("/nova", ensureParamsHeaders(http.HandlerFunc(getTeste))).Methods("POST")
	router.Handle("/", http.HandlerFunc(MainPageHandler))

	http.ListenAndServe(":3000", router)
}
