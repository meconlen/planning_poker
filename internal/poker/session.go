package poker

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

type User struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Vote     string          `json:"vote,omitempty"`
	HasVoted bool            `json:"hasVoted"`
	Conn     *websocket.Conn `json:"-"`
}

type Session struct {
	ID           string           `json:"id"`
	Story        string           `json:"story"`
	Users        map[string]*User `json:"users"`
	VotesRevealed bool            `json:"votesRevealed"`
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

var sessions = make(map[string]*Session)

func GetSession(sessionID string) *Session {
	if session, exists := sessions[sessionID]; exists {
		return session
	}
	
	// Create new session if it doesn't exist
	session := &Session{
		ID:           sessionID,
		Story:        "",
		Users:        make(map[string]*User),
		VotesRevealed: false,
	}
	sessions[sessionID] = session
	return session
}

func (s *Session) AddUser(name string, conn *websocket.Conn) *User {
	user := &User{
		ID:       uuid.New().String(),
		Name:     name,
		Vote:     "",
		HasVoted: false,
		Conn:     conn,
	}
	
	s.Users[user.ID] = user
	s.BroadcastSessionState()
	return user
}

func (s *Session) RemoveUser(userID string) {
	delete(s.Users, userID)
	s.BroadcastSessionState()
}

func (s *Session) Vote(userID, vote string) {
	if user, exists := s.Users[userID]; exists {
		user.Vote = vote
		user.HasVoted = true
		s.BroadcastSessionState()
	}
}

func (s *Session) RevealVotes() {
	s.VotesRevealed = true
	s.BroadcastSessionState()
}

func (s *Session) NewRound() {
	s.VotesRevealed = false
	for _, user := range s.Users {
		user.Vote = ""
		user.HasVoted = false
	}
	s.BroadcastSessionState()
}

func (s *Session) SetStory(story string) {
	s.Story = story
	s.BroadcastSessionState()
}

func (s *Session) BroadcastSessionState() {
	// Create a safe version of users without sensitive data
	safeUsers := make(map[string]*User)
	for id, user := range s.Users {
		safeUser := &User{
			ID:       user.ID,
			Name:     user.Name,
			HasVoted: user.HasVoted,
		}
		
		// Only include vote if votes are revealed
		if s.VotesRevealed {
			safeUser.Vote = user.Vote
		}
		
		safeUsers[id] = safeUser
	}

	sessionState := &Session{
		ID:           s.ID,
		Story:        s.Story,
		Users:        safeUsers,
		VotesRevealed: s.VotesRevealed,
	}

	message := Message{
		Type:    "session_state",
		Payload: sessionState,
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling session state: %v", err)
		return
	}

	// Broadcast to all users in the session
	for _, user := range s.Users {
		if user.Conn != nil {
			err := user.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("Error sending message to user %s: %v", user.Name, err)
			}
		}
	}
}

func (s *Session) HandleMessage(userID string, messageData []byte) {
	var msg Message
	if err := json.Unmarshal(messageData, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	switch msg.Type {
	case "vote":
		if vote, ok := msg.Payload.(string); ok {
			s.Vote(userID, vote)
		}
	case "reveal":
		s.RevealVotes()
	case "new_round":
		s.NewRound()
	case "set_story":
		if story, ok := msg.Payload.(string); ok {
			s.SetStory(story)
		}
	}
}
