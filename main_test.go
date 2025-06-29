package main

import (
	"testing"

	"planning-poker/internal/poker"
)

func TestSessionCreation(t *testing.T) {
	session := poker.NewSession("test-session")

	if session.ID != "test-session" {
		t.Errorf("Expected session ID to be 'test-session', got %s", session.ID)
	}

	if len(session.Users) != 0 {
		t.Errorf("Expected new session to have 0 users, got %d", len(session.Users))
	}

	if session.VotesRevealed {
		t.Error("Expected new session to have votes not revealed")
	}
}
