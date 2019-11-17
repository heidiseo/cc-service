package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/heidiseo/cc-service"
)

func main() {
	port := os.Getenv("PORT")
	fmt.Println(port)

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := mux.NewRouter()
	r.HandleFunc("/v1/creditcard", handler).Methods(http.MethodPost)
	err := http.ListenAndServe(":"+port, r)

	if err != nil {
		log.Fatal("error occurred")
	}
}
