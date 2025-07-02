package main

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"planning-poker/internal/poker"

	"github.com/gorilla/websocket"
)

// FrontendBugTestClient simulates a browser client to test the frontend bug
type FrontendBugTestClient struct {
	name        string
	conn        *websocket.Conn
	isModerator bool
	hasVoted    bool
	voteValue   string
	messages    []poker.Message
}

func NewFrontendBugTestClient(name string, isModerator bool) *FrontendBugTestClient {
	return &FrontendBugTestClient{
		name:        name,
		isModerator: isModerator,
		messages:    make([]poker.Message, 0),
	}
}

func (c *FrontendBugTestClient) Connect(serverURL string) error {
	// Connect to WebSocket
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(serverURL, nil)
	if err != nil {
		return err
	}
	c.conn = conn

	// Start message listener
	go c.listen()

	// Join session
	return c.sendMessage("join", map[string]interface{}{
		"sessionId": "FRONTEND-BUG-TEST",
		"userName":  c.name,
	})
}

func (c *FrontendBugTestClient) listen() {
	for {
		var msg poker.Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Client %s read error: %v", c.name, err)
			return
		}
		c.messages = append(c.messages, msg)
		c.handleMessage(msg)
	}
}

func (c *FrontendBugTestClient) handleMessage(msg poker.Message) {
	switch msg.Type {
	case "session_state":
		// This simulates the updateSessionState() function in frontend
		var state map[string]interface{}
		json.Unmarshal(msg.Data, &state)
		
		// Check if this client's vote was cleared (simulating frontend logic)
		users := state["users"].(map[string]interface{})
		for _, userData := range users {
			user := userData.(map[string]interface{})
			if user["name"] == c.name {
				// Simulate the BUG: frontend doesn't clear vote UI when vote is nil
				vote := user["vote"]
				if vote == nil && c.hasVoted {
					// This is where the frontend bug occurs!
					// In real frontend, the .selected class remains and myVote isn't cleared
					log.Printf("🐛 BUG: Client %s had vote '%s' but received vote=nil, UI should be cleared but isn't!", 
						c.name, c.voteValue)
				}
			}
		}
	}
}

func (c *FrontendBugTestClient) Vote(value string) error {
	c.hasVoted = true
	c.voteValue = value
	return c.sendMessage("vote", map[string]interface{}{
		"vote": value,
	})
}

func (c *FrontendBugTestClient) StartNewRound() error {
	if !c.isModerator {
		return fmt.Errorf("only moderator can start new round")
	}
	return c.sendMessage("new_round", nil)
}

func (c *FrontendBugTestClient) sendMessage(msgType string, data interface{}) error {
	msg := map[string]interface{}{
		"type": msgType,
	}
	if data != nil {
		dataBytes, _ := json.Marshal(data)
		msg["data"] = json.RawMessage(dataBytes)
	}
	return c.conn.WriteJSON(msg)
}

func (c *FrontendBugTestClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *FrontendBugTestClient) GetMessageCount(msgType string) int {
	count := 0
	for _, msg := range c.messages {
		if string(msg.Type) == msgType {
			count++
		}
	}
	return count
}

// Test function to reproduce the frontend bug
func TestFrontendVoteBug(t *testing.T) {
	// This test documents the frontend bug by simulating the WebSocket interaction
	
	t.Log("🧪 Testing Frontend Vote UI Bug Reproduction")
	
	// Note: This test documents the issue but requires a running server to fully test
	// For now, it serves as documentation of the expected behavior
	
	t.Log("📝 Expected behavior:")
	t.Log("1. Moderator and participant both vote")
	t.Log("2. Moderator starts new round") 
	t.Log("3. Backend clears all votes (✅ works)")
	t.Log("4. Frontend receives session_state with vote: null for all users")
	t.Log("5. Frontend should clear vote UI for ALL users (❌ BUG: only moderator UI is cleared)")
	
	t.Log("\n🐛 The bug:")
	t.Log("- newRound() function only clears moderator's vote UI")
	t.Log("- updateSessionState() doesn't detect when votes are cleared")
	t.Log("- Participants keep .selected class on voting cards")
	t.Log("- Participants' myVote variable remains set")
	
	t.Log("\n🔧 The fix needed:")
	t.Log("- Modify updateSessionState() to detect vote clearing")
	t.Log("- Clear .selected class from all voting cards")
	t.Log("- Reset myVote = null for all users")
	
	// Simulate the detection logic that should be added to frontend
	// This shows how the frontend could detect when votes are cleared
	
	// Mock session state before new round (users have votes)
	stateBefore := map[string]interface{}{
		"users": map[string]interface{}{
			"alice-id": map[string]interface{}{
				"name": "Alice",
				"vote": "5",
				"isModerator": true,
			},
			"bob-id": map[string]interface{}{
				"name": "Bob", 
				"vote": "3",
				"isModerator": false,
			},
		},
	}
	
	// Mock session state after new round (votes cleared)
	stateAfter := map[string]interface{}{
		"users": map[string]interface{}{
			"alice-id": map[string]interface{}{
				"name": "Alice",
				"vote": nil,
				"isModerator": true,
			},
			"bob-id": map[string]interface{}{
				"name": "Bob",
				"vote": nil, 
				"isModerator": false,
			},
		},
	}
	
	// Test the detection logic
	usersBefore := stateBefore["users"].(map[string]interface{})
	usersAfter := stateAfter["users"].(map[string]interface{})
	
	votesClearedDetected := false
	for userID := range usersBefore {
		beforeUser := usersBefore[userID].(map[string]interface{})
		afterUser := usersAfter[userID].(map[string]interface{})
		
		hadVote := beforeUser["vote"] != nil
		hasVote := afterUser["vote"] != nil
		
		if hadVote && !hasVote {
			votesClearedDetected = true
			t.Logf("✓ Detected vote clearing for user: %s", afterUser["name"])
		}
	}
	
	if !votesClearedDetected {
		t.Error("❌ Failed to detect vote clearing - this logic should be added to frontend")
	} else {
		t.Log("✅ Vote clearing detection logic works - this should be added to updateSessionState()")
	}
}

// Run the test if this file is executed directly
// Note: This test documents the frontend bug and detection logic
