package server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"planning-poker/internal/poker"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type JoinMessage struct {
	SessionID string `json:"sessionId"`
	UserName  string `json:"userName"`
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	var currentUser *poker.User
	var currentSession *poker.Session

	log.Println("Client connected")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		if messageType != websocket.TextMessage {
			continue
		}

		// If user hasn't joined a session yet, expect a join message
		if currentUser == nil {
			var joinMsg JoinMessage
			if err := json.Unmarshal(message, &joinMsg); err != nil {
				log.Printf("Error parsing join message: %v", err)
				continue
			}

			// Validate input
			if joinMsg.SessionID == "" || joinMsg.UserName == "" {
				log.Printf("Invalid join message: missing sessionId or userName")
				continue
			}

			// URL decode the user name in case it came from URL parameters
			decodedName, err := url.QueryUnescape(joinMsg.UserName)
			if err != nil {
				decodedName = joinMsg.UserName
			}

			currentSession = poker.GetSession(joinMsg.SessionID)
			currentUser = currentSession.AddUser(decodedName, conn)

			log.Printf("User %s joined session %s", currentUser.Name, currentSession.ID)
			continue
		}

		// Handle poker-related messages
		if currentSession != nil && currentUser != nil {
			currentSession.HandleMessage(currentUser.ID, message)
		}
	}

	// Clean up when user disconnects
	if currentSession != nil && currentUser != nil {
		currentSession.RemoveUser(currentUser.ID)
		log.Printf("User %s left session %s", currentUser.Name, currentSession.ID)
	}

	log.Println("Client disconnected")
}
