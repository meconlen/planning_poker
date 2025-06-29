package poker

import (
	"testing"
)

func TestNewSession(t *testing.T) {
	sessionID := "test-session-123"
	session := NewSession(sessionID)
	
	if session.ID != sessionID {
		t.Errorf("Expected session ID to be %s, got %s", sessionID, session.ID)
	}
	
	if session.Users == nil {
		t.Error("Expected Users map to be initialized")
	}
	
	if len(session.Users) != 0 {
		t.Errorf("Expected new session to have 0 users, got %d", len(session.Users))
	}
	
	if session.VotesRevealed {
		t.Error("Expected new session to have votes not revealed")
	}
	
	if session.CurrentStory != "" {
		t.Errorf("Expected new session to have empty story, got %s", session.CurrentStory)
	}
	
	if session.ModeratorID != "" {
		t.Errorf("Expected new session to have empty moderator ID, got %s", session.ModeratorID)
	}
}

func TestSessionAddUser(t *testing.T) {
	session := NewSession("test-session")
	
	// Add first user (should become moderator)
	user1 := session.AddUser("Alice", nil)
	
	if user1.Name != "Alice" {
		t.Errorf("Expected user name to be 'Alice', got %s", user1.Name)
	}
	
	if !user1.IsModerator {
		t.Error("Expected first user to be moderator")
	}
	
	if !user1.IsOnline {
		t.Error("Expected user to be online")
	}
	
	if session.ModeratorID != user1.ID {
		t.Error("Expected session moderator ID to match first user ID")
	}
	
	// Add second user (should not be moderator)
	user2 := session.AddUser("Bob", nil)
	
	if user2.IsModerator {
		t.Error("Expected second user to not be moderator")
	}
	
	if len(session.Users) != 2 {
		t.Errorf("Expected session to have 2 users, got %d", len(session.Users))
	}
}

func TestSessionStartNewRound(t *testing.T) {
	session := NewSession("test-session")
	user := session.AddUser("Alice", nil)
	
	// Set a vote
	vote := "5"
	user.Vote = &vote
	session.VotesRevealed = true
	
	// Start new round
	session.startNewRound()
	
	if session.VotesRevealed {
		t.Error("Expected votes to be hidden after new round")
	}
	
	if user.Vote != nil {
		t.Error("Expected user vote to be cleared after new round")
	}
}
