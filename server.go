package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/test", TestHandler)
	r.HandleFunc("/posts", PostsHandler)
	r.HandleFunc("/json", JSONHandler)
	r.HandleFunc("/register", RegisterHandler)
	r.HandleFunc("/login", LoginHandler)

	r.NotFoundHandler = http.HandlerFunc(EmptyJsonHandler)
	fmt.Println("Listening on :7700")
	log.Fatal(http.ListenAndServe(":7700", handlers.CORS(headersOk, methodsOk)(r)))
}
