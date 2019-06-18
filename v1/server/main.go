package main

import (
	"log"
	"os"

	"github.com/swrap/rdma-ds/v1"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "54005"
	}
	server := v1.CreateServer(port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error running server: %s", err)
	}
}
