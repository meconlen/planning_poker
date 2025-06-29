package main

import (
	"log"
	"net/http"
	"os"

	"planning-poker/internal/server"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create a new server instance
	srv := server.New()

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("./web/")))

	// WebSocket endpoint
	http.HandleFunc("/ws", srv.HandleWebSocket)

	// API endpoints
	http.HandleFunc("/api/sessions", srv.HandleSessions)
	http.HandleFunc("/api/sessions/", srv.HandleSession)

	log.Printf("Planning Poker server starting on :%s", port)
	log.Printf("Open http://localhost:%s in your browser", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
