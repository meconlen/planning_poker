package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"planning-poker/internal/poker"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type Server struct {
	sessions map[string]*poker.Session
	mu       sync.RWMutex
}

func New() *Server {
	return &Server{
		sessions: make(map[string]*poker.Session),
	}
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	sessionID := r.URL.Query().Get("session")
	userName := r.URL.Query().Get("user")

	if sessionID == "" || userName == "" {
		log.Println("Missing session or user parameter")
		return
	}

	s.mu.Lock()
	session, exists := s.sessions[sessionID]
	if !exists {
		session = poker.NewSession(sessionID)
		s.sessions[sessionID] = session
	}
	s.mu.Unlock()

	// Add user to session
	user := session.AddUser(userName, conn)
	defer session.RemoveUser(user.ID)

	log.Printf("User %s joined session %s", userName, sessionID)

	// Handle messages from client
	for {
		var msg poker.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		session.HandleMessage(user.ID, msg)
	}
}

func (s *Server) HandleSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	switch r.Method {
	case "GET":
		s.mu.RLock()
		sessionList := make([]string, 0, len(s.sessions))
		for id := range s.sessions {
			sessionList = append(sessionList, id)
		}
		s.mu.RUnlock()
		
		json.NewEncoder(w).Encode(map[string]interface{}{
			"sessions": sessionList,
		})
		
	case "POST":
		var req struct {
			SessionID string `json:"sessionId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		
		s.mu.Lock()
		if _, exists := s.sessions[req.SessionID]; !exists {
			s.sessions[req.SessionID] = poker.NewSession(req.SessionID)
		}
		s.mu.Unlock()
		
		json.NewEncoder(w).Encode(map[string]string{
			"sessionId": req.SessionID,
			"status":    "created",
		})
	}
}

func (s *Server) HandleSession(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from URL path
	sessionID := strings.TrimPrefix(r.URL.Path, "/api/sessions/")
	
	s.mu.RLock()
	session, exists := s.sessions[sessionID]
	s.mu.RUnlock()
	
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session.GetState())
}
