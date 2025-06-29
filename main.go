package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"planning-poker/internal/config"
	"planning-poker/internal/server"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create a new server instance with configuration
	srv := server.NewWithConfig(cfg)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         cfg.Address(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("./web/")))

	// WebSocket endpoint
	http.HandleFunc("/ws", srv.HandleWebSocket)

	// API endpoints
	http.HandleFunc("/api/sessions", srv.HandleSessions)
	http.HandleFunc("/api/sessions/", srv.HandleSession)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"planning-poker"}`))
	})

	log.Printf("Planning Poker server starting on %s", cfg.Address())
	log.Printf("Environment: %s", map[bool]string{true: "development", false: "production"}[cfg.IsDevelopment])
	log.Printf("Open http://localhost:%s in your browser", cfg.Port)

	// Start server in a goroutine
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}
