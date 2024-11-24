package main

import (
	"log"
	"matrix-compute/internal/server"
)

func main() {
	srv := server.NewServer()
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
