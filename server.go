package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/test", TestHandler)
	r.HandleFunc("/posts", PostsHandler)
	r.HandleFunc("/json", JSONHandler)

	r.NotFoundHandler = http.HandlerFunc(emptyJsonHandler)
	fmt.Println("Listening on :7700")
	log.Fatal(http.ListenAndServe(":7700", handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
