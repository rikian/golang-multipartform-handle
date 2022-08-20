package app

import (
	"go/upload/app/middleware"
	"log"
	"net/http"
)

func ListenAndServe(address string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", middleware.Middleware)

	var server http.Server = http.Server{
		Addr:    address,
		Handler: mux,
	}

	log.Print("server listen and serve at " + address)

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err.Error())
	}
}
