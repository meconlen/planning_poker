package poker

import (
	"testing"
)

func TestSessionCreation(t *testing.T) {
	session := NewSession("TEST123")
	
	if session.ID != "TEST123" {
		t.Errorf("Expected session ID 'TEST123', got %s", session.ID)
	}
	
	if session.Status != SessionStatusWaiting {
		t.Errorf("Expected new session status to be 'waiting', got %s", session.Status)
	}
	
	if len(session.Users) != 0 {
		t.Errorf("Expected new session to have no users, got %d", len(session.Users))
	}
}

func TestAddUserAsCreator(t *testing.T) {
	session := NewSession("TEST123")
	
	// Add user as creator
	user := session.AddUser("Alice", nil, true)
	
	if user.Name != "Alice" {
		t.Errorf("Expected user name 'Alice', got %s", user.Name)
	}
	
	if !user.IsModerator {
		t.Error("Creator should be moderator")
	}
	
	if session.CreatorID != user.ID {
		t.Error("Session creator ID should match user ID")
	}
	
	if session.ModeratorID != user.ID {
		t.Error("Session moderator ID should match user ID")
	}
}

func TestAddUserAsParticipant(t *testing.T) {
	session := NewSession("TEST123")
	
	// Add creator first
	creator := session.AddUser("Alice", nil, true)
	
	// Add participant
	participant := session.AddUser("Bob", nil, false)
	
	if participant.Name != "Bob" {
		t.Errorf("Expected participant name 'Bob', got %s", participant.Name)
	}
	
	if participant.IsModerator {
		t.Error("Participant should not be moderator")
	}
	
	if session.CreatorID != creator.ID {
		t.Error("Session creator ID should still be the first user")
	}
	
	if len(session.Users) != 2 {
		t.Errorf("Expected 2 users in session, got %d", len(session.Users))
	}
}

func TestStartSessionAsCreator(t *testing.T) {
	session := NewSession("TEST123")
	
	// Add creator
	creator := session.AddUser("Alice", nil, true)
	
	// Start session
	success := session.StartSession(creator.ID)
	
	if !success {
		t.Error("Creator should be able to start session")
	}
	
	if session.Status != SessionStatusActive {
		t.Errorf("Expected session status to be 'active', got %s", session.Status)
	}
}

func TestStartSessionAsNonCreator(t *testing.T) {
	session := NewSession("TEST123")
	
	// Add creator
	creator := session.AddUser("Alice", nil, true)
	
	// Add participant
	participant := session.AddUser("Bob", nil, false)
	
	// Try to start session as participant
	success := session.StartSession(participant.ID)
	
	if success {
		t.Error("Non-creator should not be able to start session")
	}
	
	if session.Status != SessionStatusWaiting {
		t.Errorf("Expected session status to remain 'waiting', got %s", session.Status)
	}
	
	// Creator should still be able to start
	success = session.StartSession(creator.ID)
	if !success {
		t.Error("Creator should be able to start session")
	}
}

func TestVoting(t *testing.T) {
	session := NewSession("TEST123")
	
	// Add users and start session
	creator := session.AddUser("Alice", nil, true)
	participant := session.AddUser("Bob", nil, false)
	session.StartSession(creator.ID)
	
	// Test voting
	voteMsg := Message{
		Type: MessageTypeVote,
		Data: mustMarshal(map[string]string{"vote": "5"}),
	}
	
	session.HandleMessage(creator.ID, voteMsg)
	
	if creator.Vote == nil || *creator.Vote != "5" {
		t.Errorf("Expected creator vote to be '5', got %v", creator.Vote)
	}
	
	// Test participant voting
	participantVoteMsg := Message{
		Type: MessageTypeVote,
		Data: mustMarshal(map[string]string{"vote": "8"}),
	}
	
	session.HandleMessage(participant.ID, participantVoteMsg)
	
	if participant.Vote == nil || *participant.Vote != "8" {
		t.Errorf("Expected participant vote to be '8', got %v", participant.Vote)
	}
}

func TestRevealVotes(t *testing.T) {
	session := NewSession("TEST123")
	
	// Add users and start session
	creator := session.AddUser("Alice", nil, true)
	participant := session.AddUser("Bob", nil, false)
	session.StartSession(creator.ID)
	
	// Cast votes
	vote1 := "5"
	vote2 := "8"
	creator.Vote = &vote1
	participant.Vote = &vote2
	
	// Reveal votes (only moderator can do this)
	revealMsg := Message{
		Type: MessageTypeReveal,
	}
	
	session.HandleMessage(creator.ID, revealMsg)
	
	if !session.VotesRevealed {
		t.Error("Votes should be revealed after moderator reveals them")
	}
	
	// Test that non-moderator cannot reveal votes
	session.VotesRevealed = false
	session.HandleMessage(participant.ID, revealMsg)
	
	if session.VotesRevealed {
		t.Error("Non-moderator should not be able to reveal votes")
	}
}

func TestNewRound(t *testing.T) {
	session := NewSession("TEST123")
	
	// Add users and start session
	creator := session.AddUser("Alice", nil, true)
	participant := session.AddUser("Bob", nil, false)
	session.StartSession(creator.ID)
	
	// Cast votes and reveal
	vote1 := "5"
	vote2 := "8"
	creator.Vote = &vote1
	participant.Vote = &vote2
	session.VotesRevealed = true
	
	// Start new round
	newRoundMsg := Message{
		Type: MessageTypeNewRound,
	}
	
	session.HandleMessage(creator.ID, newRoundMsg)
	
	if session.VotesRevealed {
		t.Error("Votes should be hidden after new round")
	}
	
	if creator.Vote != nil {
		t.Error("Creator vote should be cleared after new round")
	}
	
	if participant.Vote != nil {
		t.Error("Participant vote should be cleared after new round")
	}
}

func TestSetStory(t *testing.T) {
	session := NewSession("TEST123")
	
	// Add users and start session
	creator := session.AddUser("Alice", nil, true)
	participant := session.AddUser("Bob", nil, false)
	session.StartSession(creator.ID)
	
	// Set story as moderator
	storyMsg := Message{
		Type: MessageTypeSetStory,
		Data: mustMarshal(map[string]string{"story": "User can login"}),
	}
	
	session.HandleMessage(creator.ID, storyMsg)
	
	if session.CurrentStory != "User can login" {
		t.Errorf("Expected story to be 'User can login', got %s", session.CurrentStory)
	}
	
	// Test that non-moderator cannot set story
	participantStoryMsg := Message{
		Type: MessageTypeSetStory,
		Data: mustMarshal(map[string]string{"story": "Should not work"}),
	}
	
	session.HandleMessage(participant.ID, participantStoryMsg)
	
	if session.CurrentStory == "Should not work" {
		t.Error("Non-moderator should not be able to set story")
	}
}

func TestSessionStateIncluesStatus(t *testing.T) {
	session := NewSession("TEST123")
	
	// Test waiting state
	state := session.GetState().(map[string]interface{})
	if state["status"] != SessionStatusWaiting {
		t.Errorf("Expected session state to include status 'waiting', got %v", state["status"])
	}
	
	// Add creator and start session
	creator := session.AddUser("Alice", nil, true)
	session.StartSession(creator.ID)
	
	// Test active state
	activeState := session.GetState().(map[string]interface{})
	if activeState["status"] != SessionStatusActive {
		t.Errorf("Expected session state to include status 'active', got %v", activeState["status"])
	}
}
