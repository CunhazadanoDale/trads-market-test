package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("API rodando em " + server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}