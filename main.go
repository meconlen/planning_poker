package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"planning-poker/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("./web/")))
	
	// WebSocket endpoint
	http.HandleFunc("/ws", server.HandleWebSocket)

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
