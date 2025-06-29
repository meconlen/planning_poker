package poker

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageType string

const (
	MessageTypeVote         MessageType = "vote"
	MessageTypeReveal       MessageType = "reveal"
	MessageTypeNewRound     MessageType = "new_round"
	MessageTypeSetStory     MessageType = "set_story"
	MessageTypeUserJoined   MessageType = "user_joined"
	MessageTypeUserLeft     MessageType = "user_left"
	MessageTypeSessionState MessageType = "session_state"
	MessageTypeStartSession MessageType = "start_session"
	MessageTypeWaitingRoom  MessageType = "waiting_room"
)

type SessionStatus string

const (
	SessionStatusWaiting SessionStatus = "waiting" // Waiting for moderator to start
	SessionStatusActive  SessionStatus = "active"  // Session is active
	SessionStatusEnded   SessionStatus = "ended"   // Session has ended
)

type Message struct {
	Type   MessageType     `json:"type"`
	Data   json.RawMessage `json:"data,omitempty"`
	UserID string          `json:"userId,omitempty"`
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
	ID            string           `json:"id"`
	Users         map[string]*User `json:"users"`
	CurrentStory  string           `json:"currentStory"`
	VotesRevealed bool             `json:"votesRevealed"`
	ModeratorID   string           `json:"moderatorId"`
	CreatorID     string           `json:"creatorId"`     // Who created the session
	Status        SessionStatus    `json:"status"`       // Session status
	CreatedAt     time.Time        `json:"createdAt"`
	mu            sync.RWMutex     `json:"-"`
}

func NewSession(id string) *Session {
	return &Session{
		ID:            id,
		Users:         make(map[string]*User),
		VotesRevealed: false,
		Status:        SessionStatusWaiting,
		CreatedAt:     time.Now(),
	}
}

func (s *Session) AddUser(name string, conn *websocket.Conn, isCreator bool) *User {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := &User{
		ID:          uuid.New().String(),
		Name:        name,
		IsOnline:    true,
		IsModerator: isCreator, // Set moderator status if creator
		conn:        conn,
	}

	s.Users[user.ID] = user

	// Set creator information if this is the creator
	if isCreator {
		s.CreatorID = user.ID
		s.ModeratorID = user.ID
	}

	// Notify all users about the new user
	s.broadcastMessage(Message{
		Type: MessageTypeUserJoined,
		Data: mustMarshal(user),
	})

	// Send appropriate state based on session status and creator status
	if s.Status == SessionStatusWaiting && !isCreator {
		// Send waiting room message to non-creators
		user.sendMessage(Message{
			Type: MessageTypeWaitingRoom,
			Data: mustMarshal(map[string]interface{}{
				"sessionId": s.ID,
				"message":   "Waiting for the session creator to start the session...",
			}),
		})
	}
	
	// Always send session state so users know about other participants
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
		log.Printf("HandleMessage: User %s not found in session", userID)
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

	case MessageTypeStartSession:
		// Only allow creator to start session
		if s.CreatorID != userID {
			log.Printf("User %s attempted to start session but is not creator", user.Name)
			return
		}
		
		if !s.startSessionUnsafe(userID) {
			log.Printf("Failed to start session for user %s", user.Name)
		}
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
		"id":            s.ID,
		"users":         users,
		"currentStory":  s.CurrentStory,
		"votesRevealed": s.VotesRevealed,
		"status":        s.Status,
		"createdAt":     s.CreatedAt,
	}
}

func (s *Session) broadcastMessage(msg Message) {
	for userID, user := range s.Users {
		if user.IsOnline {
			user.sendMessage(msg)
		} else {
			log.Printf("Skipping offline user %s (%s)", user.Name, userID)
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

// SetCreator sets the session creator and makes them the moderator
func (s *Session) SetCreator(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.CreatorID = userID
	s.ModeratorID = userID
	
	// Make the creator the moderator
	if user, exists := s.Users[userID]; exists {
		user.IsModerator = true
	}
}

// StartSession starts the session (only creator can do this)
func (s *Session) StartSession(userID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Only the creator can start the session
	if s.CreatorID != userID {
		return false
	}

	if s.Status != SessionStatusWaiting {
		return false
	}

	s.Status = SessionStatusActive
	log.Printf("Session %s started by %s", s.ID, s.Users[userID].Name)

	// Notify all users that session has started
	s.broadcastMessage(Message{
		Type: MessageTypeStartSession,
		Data: mustMarshal(map[string]interface{}{
			"message": "Session has been started by the moderator!",
		}),
	})

	// Send current session state to all users
	s.broadcastSessionState()

	return true
}

// startSessionUnsafe starts the session without acquiring locks (for internal use)
func (s *Session) startSessionUnsafe(userID string) bool {
	// Only the creator can start the session
	if s.CreatorID != userID {
		return false
	}

	if s.Status != SessionStatusWaiting {
		return false
	}

	s.Status = SessionStatusActive
	log.Printf("Session %s started by %s", s.ID, s.Users[userID].Name)

	// Notify all users that session has started
	s.broadcastMessage(Message{
		Type: MessageTypeStartSession,
		Data: mustMarshal(map[string]interface{}{
			"message": "Session has been started by the moderator!",
		}),
	})

	// Send current session state to all users
	s.broadcastSessionState()

	return true
}

// broadcastSessionState sends the current session state to all users
func (s *Session) broadcastSessionState() {
	state := s.getStateUnsafe()
	s.broadcastMessage(Message{
		Type: MessageTypeSessionState,
		Data: mustMarshal(state),
	})
}
