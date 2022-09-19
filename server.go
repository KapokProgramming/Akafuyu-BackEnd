package main

import (
	"fmt"
	"net/http"
	"time"

	"server/config"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/test", TestHandler)
	r.HandleFunc("/post/{id:[0-9]+}", PostHandler)
	r.HandleFunc("/posts", PostsHandler)
	r.HandleFunc("/json", JSONHandler)
	r.HandleFunc("/register", RegisterHandler)
	r.HandleFunc("/login", LoginHandler)

	r.NotFoundHandler = http.HandlerFunc(EmptyJsonHandler)

	cfg := config.LoadConfig()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      handlers.CORS(headersOk, methodsOk)(r),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	fmt.Println("Listening on :7700")
	err := srv.ListenAndServe()
	if err != nil {
		panic("Failed to run server")
	}

	// log.Fatal(http.ListenAndServe(":7700", )
}
