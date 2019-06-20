package main

import (
	"log"
	"os"

	"github.com/swrap/rdma-ds/src"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "54005"
	}
	server := src.CreateServer(port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error running server: %s", err)
	}
}
