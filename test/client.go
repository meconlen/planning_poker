package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type TestMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

type TestClient struct {
	name      string
	sessionID string
	isCreator bool
	conn      *websocket.Conn
	messages  []TestMessage
	done      chan struct{}
}

func NewTestClient(name, sessionID string, isCreator bool) *TestClient {
	return &TestClient{
		name:      name,
		sessionID: sessionID,
		isCreator: isCreator,
		messages:  make([]TestMessage, 0),
		done:      make(chan struct{}),
	}
}

func (c *TestClient) Connect(serverURL string) error {
	u := url.URL{Scheme: "ws", Host: serverURL, Path: "/ws"}
	q := u.Query()
	q.Set("session", c.sessionID)
	q.Set("user", c.name)
	if c.isCreator {
		q.Set("creator", "true")
	}
	u.RawQuery = q.Encode()

	log.Printf("[%s] Connecting to %s", c.name, u.String())

	var err error
	c.conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	go c.readMessages()
	return nil
}

func (c *TestClient) readMessages() {
	defer close(c.done)
	for {
		var msg TestMessage
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("[%s] Read error: %v", c.name, err)
			return
		}

		c.messages = append(c.messages, msg)
		c.logMessage(msg)
	}
}

func (c *TestClient) logMessage(msg TestMessage) {
	log.Printf("[%s] Received: %s", c.name, msg.Type)

	switch msg.Type {
	case "session_state":
		var state map[string]interface{}
		if err := json.Unmarshal(msg.Data, &state); err == nil {
			users := state["users"].(map[string]interface{})
			status := state["status"]
			log.Printf("[%s]   Session Status: %v", c.name, status)
			log.Printf("[%s]   Participants (%d):", c.name, len(users))
			for _, user := range users {
				u := user.(map[string]interface{})
				name := u["name"].(string)
				isModerator := u["isModerator"].(bool)
				moderatorTag := ""
				if isModerator {
					moderatorTag = " (MODERATOR)"
				}
				log.Printf("[%s]     - %s%s", c.name, name, moderatorTag)
			}
		}
	case "waiting_room":
		var data map[string]interface{}
		if err := json.Unmarshal(msg.Data, &data); err == nil {
			log.Printf("[%s]   Message: %v", c.name, data["message"])
		}
	case "start_session":
		var data map[string]interface{}
		if err := json.Unmarshal(msg.Data, &data); err == nil {
			log.Printf("[%s]   Message: %v", c.name, data["message"])
		}
	}
}

func (c *TestClient) SendMessage(msgType string, data interface{}) error {
	msg := map[string]interface{}{
		"type": msgType,
	}
	if data != nil {
		msg["data"] = data
	}

	log.Printf("[%s] Sending: %s", c.name, msgType)
	return c.conn.WriteJSON(msg)
}

func (c *TestClient) StartSession() error {
	return c.SendMessage("start_session", nil)
}

func (c *TestClient) Vote(vote string) error {
	return c.SendMessage("vote", map[string]string{"vote": vote})
}

func (c *TestClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *TestClient) Wait() {
	<-c.done
}

func (c *TestClient) GetMessageCount(msgType string) int {
	count := 0
	for _, msg := range c.messages {
		if msg.Type == msgType {
			count++
		}
	}
	return count
}

func (c *TestClient) GetLastMessage(msgType string) *TestMessage {
	for i := len(c.messages) - 1; i >= 0; i-- {
		if c.messages[i].Type == msgType {
			return &c.messages[i]
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run test_client.go <server_host:port>")
	}

	serverURL := os.Args[1]
	sessionID := "TEST123"

	log.Printf("=== Planning Poker Test Client ===")
	log.Printf("Server: %s", serverURL)
	log.Printf("Session: %s", sessionID)

	// Connect moderator (creator)
	log.Printf("\n=== Connecting Moderator (Alice) ===")
	moderator := NewTestClient("Alice", sessionID, true)
	if err := moderator.Connect(serverURL); err != nil {
		log.Fatal(err)
	}
	defer moderator.Close()

	// Wait a bit for connection
	time.Sleep(500 * time.Millisecond)

	// Connect participant
	log.Printf("\n=== Connecting Participant (Bob) ===")
	participant := NewTestClient("Bob", sessionID, false)
	if err := participant.Connect(serverURL); err != nil {
		log.Fatal(err)
	}
	defer participant.Close()

	// Wait for initial messages
	time.Sleep(1 * time.Second)

	// Check initial state
	log.Printf("\n=== Initial State Analysis ===")
	log.Printf("Moderator received %d session_state messages", moderator.GetMessageCount("session_state"))
	log.Printf("Participant received %d session_state messages", participant.GetMessageCount("session_state"))
	log.Printf("Participant received %d waiting_room messages", participant.GetMessageCount("waiting_room"))

	// Analyze the participant's last session state
	if lastState := participant.GetLastMessage("session_state"); lastState != nil {
		var state map[string]interface{}
		if err := json.Unmarshal(lastState.Data, &state); err == nil {
			users := state["users"].(map[string]interface{})
			log.Printf("Participant can see %d users in session state", len(users))
		}
	}

	// Start session
	log.Printf("\n=== Moderator Starting Session ===")
	if err := moderator.StartSession(); err != nil {
		log.Printf("Error starting session: %v", err)
	}

	// Wait for session start messages
	time.Sleep(1 * time.Second)

	// Check final state
	log.Printf("\n=== Post-Start State Analysis ===")
	log.Printf("Moderator received %d session_state messages total", moderator.GetMessageCount("session_state"))
	log.Printf("Participant received %d session_state messages total", participant.GetMessageCount("session_state"))
	log.Printf("Participant received %d start_session messages", participant.GetMessageCount("start_session"))

	// Check if participant transitioned from waiting room
	if lastState := participant.GetLastMessage("session_state"); lastState != nil {
		var state map[string]interface{}
		if err := json.Unmarshal(lastState.Data, &state); err == nil {
			status := state["status"]
			users := state["users"].(map[string]interface{})
			log.Printf("Participant's final session status: %v", status)
			log.Printf("Participant can see %d users in final state", len(users))
		}
	}

	// Test voting
	log.Printf("\n=== Testing Voting ===")
	moderator.Vote("5")
	participant.Vote("8")

	time.Sleep(1 * time.Second)

	log.Printf("\n=== Test Summary ===")
	log.Printf("Total messages - Moderator: %d, Participant: %d", len(moderator.messages), len(participant.messages))

	// Handle cleanup
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Printf("\n=== Shutting down ===")
		moderator.Close()
		participant.Close()
		os.Exit(0)
	}()

	log.Printf("\n=== Test complete - press Ctrl+C to exit ===")
	select {}
}
