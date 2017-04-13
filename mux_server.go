package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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

	router := mux.NewRouter().StrictSlash(true)

	router.Handle("/subscribe", ensureSubscribe(http.HandlerFunc(subscribeChannel))).Methods("GET")
	router.Handle("/api/v1/channel", ensureCreateChannel(http.HandlerFunc(createChannel))).Methods("POST")
	router.Handle("/", http.HandlerFunc(MainPageHandler))

	http.ListenAndServe(":3000", router)
}
