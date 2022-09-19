package main

import (
	"fmt"
	"net/http"
	"time"

	"server/config"
	"server/routes"
)

func main() {

	cfg := config.LoadConfig()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      routes.CreateRoute(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	fmt.Println("Listening on :7700")
	err := srv.ListenAndServe()
	if err != nil {
		panic("Failed to run server")
	}

}
