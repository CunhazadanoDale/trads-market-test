package main

import (
	"log"
	"net/http"

	appRouter "github.com/CunhazadanoDale/trads-market-test/internal/adapter/http"
)

func main() {
	router := appRouter.NewRouter()

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("API rodando em " + server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
