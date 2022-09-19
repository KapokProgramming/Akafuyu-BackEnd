package routes

import (
	"net/http"

	"server/handler"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func CreateRoute() http.Handler {
	r := mux.NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	r.HandleFunc("/", handler.HomeHandler)
	r.HandleFunc("/test", handler.TestHandler)
	r.HandleFunc("/post/{id:[0-9]+}", handler.PostHandler)
	r.HandleFunc("/posts", handler.PostsHandler)
	r.HandleFunc("/json", handler.JSONHandler)
	r.HandleFunc("/register", handler.RegisterHandler)
	r.HandleFunc("/login", handler.LoginHandler)

	r.NotFoundHandler = http.HandlerFunc(handler.EmptyJsonHandler)

	return handlers.CORS(headersOk, methodsOk)(r)
}
