package poker

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

type MessageType string

const (
	MessageTypeVote       MessageType = "vote"
	MessageTypeReveal     MessageType = "reveal"
	MessageTypeNewRound   MessageType = "new_round"
	MessageTypeSetStory   MessageType = "set_story"
	MessageTypeUserJoined MessageType = "user_joined"
	MessageTypeUserLeft   MessageType = "user_left"
	MessageTypeSessionState MessageType = "session_state"
)

type Message struct {
	Type    MessageType     `json:"type"`
	Data    json.RawMessage `json:"data,omitempty"`
	UserID  string          `json:"userId,omitempty"`
}

type User struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Vote        *string         `json:"vote"`
	IsOnline    bool            `json:"isOnline"`
	IsModerator bool            `json:"isModerator"`
	conn        *websocket.Conn `json:"-"`
}

type Session struct {
	ID            string             `json:"id"`
	Users         map[string]*User   `json:"users"`
	CurrentStory  string            `json:"currentStory"`
	VotesRevealed bool             `json:"votesRevealed"`
	ModeratorID   string            `json:"moderatorId"`
	CreatedAt     time.Time          `json:"createdAt"`
	mu            sync.RWMutex       `json:"-"`
}

func NewSession(id string) *Session {
	return &Session{
		ID:          id,
		Users:       make(map[string]*User),
		VotesRevealed: false,
		CreatedAt:   time.Now(),
	}
}

func (s *Session) AddUser(name string, conn *websocket.Conn) *User {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if this is the first user (moderator)
	isModerator := len(s.Users) == 0
	if isModerator {
		s.ModeratorID = uuid.New().String()
	}

	user := &User{
		ID:          uuid.New().String(),
		Name:        name,
		IsOnline:    true,
		IsModerator: isModerator,
		conn:        conn,
	}

	// If this is the moderator, use the pre-generated moderator ID
	if isModerator {
		user.ID = s.ModeratorID
	}

	s.Users[user.ID] = user
	
	// Notify all users about the new user
	s.broadcastMessage(Message{
		Type: MessageTypeUserJoined,
		Data: mustMarshal(user),
	})
	
	// Send current session state to the new user
	user.sendMessage(Message{
		Type: MessageTypeSessionState,
		Data: mustMarshal(s.getStateUnsafe()),
	})

	return user
}

func (s *Session) RemoveUser(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if user, exists := s.Users[userID]; exists {
		user.IsOnline = false
		delete(s.Users, userID)
		
		s.broadcastMessage(Message{
			Type: MessageTypeUserLeft,
			Data: mustMarshal(map[string]string{"userId": userID}),
		})
	}
}

func (s *Session) HandleMessage(userID string, msg Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.Users[userID]
	if !exists {
		return
	}

	switch msg.Type {
	case MessageTypeVote:
		var voteData struct {
			Vote string `json:"vote"`
		}
		if err := json.Unmarshal(msg.Data, &voteData); err != nil {
			log.Printf("Invalid vote data: %v", err)
			return
		}
		
		user.Vote = &voteData.Vote
		
		// Broadcast the vote (hidden) to all users
		s.broadcastMessage(Message{
			Type: MessageTypeSessionState,
			Data: mustMarshal(s.getStateUnsafe()),
		})

	case MessageTypeReveal:
		// Only allow moderator to reveal votes
		if !user.IsModerator {
			log.Printf("User %s attempted to reveal votes but is not moderator", user.Name)
			return
		}
		s.VotesRevealed = true
		s.broadcastMessage(Message{
			Type: MessageTypeSessionState,
			Data: mustMarshal(s.getStateUnsafe()),
		})

	case MessageTypeNewRound:
		// Only allow moderator to start new rounds
		if !user.IsModerator {
			log.Printf("User %s attempted to start new round but is not moderator", user.Name)
			return
		}
		s.startNewRound()
		s.broadcastMessage(Message{
			Type: MessageTypeSessionState,
			Data: mustMarshal(s.getStateUnsafe()),
		})

	case MessageTypeSetStory:
		// Only allow moderator to set stories
		if !user.IsModerator {
			log.Printf("User %s attempted to set story but is not moderator", user.Name)
			return
		}
		var storyData struct {
			Story string `json:"story"`
		}
		if err := json.Unmarshal(msg.Data, &storyData); err != nil {
			log.Printf("Invalid story data: %v", err)
			return
		}
		
		s.CurrentStory = storyData.Story
		s.startNewRound() // Reset votes when setting new story
		s.broadcastMessage(Message{
			Type: MessageTypeSessionState,
			Data: mustMarshal(s.getStateUnsafe()),
		})
	}
}

func (s *Session) startNewRound() {
	s.VotesRevealed = false
	for _, user := range s.Users {
		user.Vote = nil
	}
}

func (s *Session) GetState() interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getStateUnsafe()
}

func (s *Session) getStateUnsafe() interface{} {
	// Create a copy of users with vote visibility based on reveal state
	users := make(map[string]*User)
	for id, user := range s.Users {
		userCopy := *user
		userCopy.conn = nil // Don't include connection in JSON
		
		// Hide votes if not revealed
		if !s.VotesRevealed && userCopy.Vote != nil {
			hiddenVote := "?"
			userCopy.Vote = &hiddenVote
		}
		
		users[id] = &userCopy
	}

	return map[string]interface{}{
		"id":           s.ID,
		"users":        users,
		"currentStory": s.CurrentStory,
		"votesRevealed": s.VotesRevealed,
		"createdAt":    s.CreatedAt,
	}
}

func (s *Session) broadcastMessage(msg Message) {
	for _, user := range s.Users {
		if user.IsOnline {
			user.sendMessage(msg)
		}
	}
}

func (u *User) sendMessage(msg Message) {
	if u.conn == nil {
		return
	}
	
	if err := u.conn.WriteJSON(msg); err != nil {
		log.Printf("Error sending message to user %s: %v", u.Name, err)
		u.IsOnline = false
	}
}

func mustMarshal(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		return json.RawMessage("{}")
	}
	return data
}
