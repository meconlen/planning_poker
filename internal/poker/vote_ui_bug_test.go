package poker

import (
	"encoding/json"
	"testing"
)

// TestVoteUINotClearedForNonModerators replicates the issue where vote buttons
// remain visually selected after a new round for non-moderator participants.
// This test verifies that the backend properly sends session state updates
// that would allow the frontend to detect when votes should be cleared.
func TestVoteUINotClearedForNonModerators(t *testing.T) {
	// Simulate the bug scenario:
	// 1. Create session with moderator and participant
	// 2. Both users vote
	// 3. Moderator starts new round
	// 4. Verify that participant would receive proper signals to clear their UI
	
	session := NewSession("TEST-UI-BUG")
	
	// Add moderator (Alice) and participant (Bob)
	moderator := session.AddUser("Alice", nil, true)
	participant := session.AddUser("Bob", nil, false)
	session.StartSession(moderator.ID)
	
	// Both users cast votes
	moderatorVote := "5"
	participantVote := "8"
	moderator.Vote = &moderatorVote
	participant.Vote = &participantVote
	
	// Verify initial state - both users have votes
	if moderator.Vote == nil {
		t.Error("Moderator should have a vote before new round")
	}
	if participant.Vote == nil {
		t.Error("Participant should have a vote before new round")
	}
	
	// Get initial session state to simulate what frontend would receive
	initialState := session.GetState()
	initialStateJSON, _ := json.Marshal(initialState)
	var initialStateParsed map[string]interface{}
	json.Unmarshal(initialStateJSON, &initialStateParsed)
	
	// Verify both users have votes in initial state
	users := initialStateParsed["users"].(map[string]interface{})
	moderatorData := users[moderator.ID].(map[string]interface{})
	participantData := users[participant.ID].(map[string]interface{})
	
	if moderatorData["vote"] == nil {
		t.Error("Initial state should show moderator has voted")
	}
	if participantData["vote"] == nil {
		t.Error("Initial state should show participant has voted")
	}
	
	// Moderator starts new round
	newRoundMsg := Message{
		Type: MessageTypeNewRound,
	}
	session.HandleMessage(moderator.ID, newRoundMsg)
	
	// Verify backend properly clears votes
	if moderator.Vote != nil {
		t.Error("Moderator vote should be cleared after new round")
	}
	if participant.Vote != nil {
		t.Error("Participant vote should be cleared after new round")
	}
	
	// Get new session state that would be broadcast to all clients
	newState := session.GetState()
	newStateJSON, _ := json.Marshal(newState)
	var newStateParsed map[string]interface{}
	json.Unmarshal(newStateJSON, &newStateParsed)
	
	// Verify the new state shows no votes for any user
	newUsers := newStateParsed["users"].(map[string]interface{})
	newModeratorData := newUsers[moderator.ID].(map[string]interface{})
	newParticipantData := newUsers[participant.ID].(map[string]interface{})
	
	if newModeratorData["vote"] != nil {
		t.Error("New state should show moderator has no vote after new round")
	}
	if newParticipantData["vote"] != nil {
		t.Error("New state should show participant has no vote after new round")
	}
	
	// Verify votes are not revealed in new round
	if newStateParsed["votesRevealed"].(bool) {
		t.Error("Votes should not be revealed after starting new round")
	}
	
	// Critical test: Verify the frontend would have enough information to detect
	// that a new round has started and should clear vote selections
	
	// The issue is that the frontend needs to detect when votes have been cleared
	// and reset the UI accordingly. The backend properly clears the votes,
	// but the frontend logic in updateSessionState() doesn't handle this.
	
	// This test documents the expected behavior:
	// 1. Backend correctly clears all votes ✓
	// 2. Backend sends updated session state ✓ 
	// 3. Frontend should detect cleared votes and reset UI (MISSING)
	
	t.Logf("Backend correctly clears votes for all users after new round")
	t.Logf("Frontend should detect vote clearing in updateSessionState() function")
	t.Logf("Issue: Frontend only clears vote UI for moderator, not for participants")
}

// TestDetectNewRoundInSessionState tests the ability to detect when a new round
// has started by comparing session states, which the frontend could use.
func TestDetectNewRoundInSessionState(t *testing.T) {
	session := NewSession("TEST-DETECT-NEW-ROUND")
	
	// Add users
	moderator := session.AddUser("Alice", nil, true)
	participant := session.AddUser("Bob", nil, false)
	session.StartSession(moderator.ID)
	
	// Users vote
	moderatorVote := "3"
	participantVote := "5"
	moderator.Vote = &moderatorVote
	participant.Vote = &participantVote
	
	// Capture state before new round
	stateBefore := session.GetState()
	stateBeforeJSON, _ := json.Marshal(stateBefore)
	var beforeParsed map[string]interface{}
	json.Unmarshal(stateBeforeJSON, &beforeParsed)
	
	// Start new round
	newRoundMsg := Message{Type: MessageTypeNewRound}
	session.HandleMessage(moderator.ID, newRoundMsg)
	
	// Capture state after new round
	stateAfter := session.GetState()
	stateAfterJSON, _ := json.Marshal(stateAfter)
	var afterParsed map[string]interface{}
	json.Unmarshal(stateAfterJSON, &afterParsed)
	
	// Test detection logic that frontend could use
	beforeUsers := beforeParsed["users"].(map[string]interface{})
	afterUsers := afterParsed["users"].(map[string]interface{})
	
	// Check if any user had votes before but not after (indicating new round)
	newRoundDetected := false
	for userID := range beforeUsers {
		beforeUser := beforeUsers[userID].(map[string]interface{})
		afterUser := afterUsers[userID].(map[string]interface{})
		
		hadVoteBefore := beforeUser["vote"] != nil
		hasVoteAfter := afterUser["vote"] != nil
		
		if hadVoteBefore && !hasVoteAfter {
			newRoundDetected = true
			break
		}
	}
	
	if !newRoundDetected {
		t.Error("Should be able to detect new round by comparing vote states")
	}
	
	// Also verify votes are not revealed after new round
	if afterParsed["votesRevealed"].(bool) {
		t.Error("Votes should not be revealed after new round")
	}
	
	t.Logf("Successfully detected new round by comparing session states")
	t.Logf("Frontend can use this logic to clear vote UI for all users")
}
