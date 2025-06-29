package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"planning-poker/internal/poker"
)

func TestNew(t *testing.T) {
	server := New()

	if server == nil {
		t.Fatal("Expected server to be created, got nil")
	}

	if server.sessions == nil {
		t.Error("Expected sessions map to be initialized")
	}

	if len(server.sessions) != 0 {
		t.Errorf("Expected new server to have 0 sessions, got %d", len(server.sessions))
	}
}

func TestHandleSessions_GET(t *testing.T) {
	server := New()

	// Add a test session
	server.sessions["TEST123"] = poker.NewSession("TEST123")

	req, err := http.NewRequest("GET", "/api/sessions", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleSessions)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Could not parse JSON response: %v", err)
	}

	sessions, exists := response["sessions"]
	if !exists {
		t.Error("Expected 'sessions' field in response")
	}

	sessionList := sessions.([]interface{})
	if len(sessionList) != 1 {
		t.Errorf("Expected 1 session in response, got %d", len(sessionList))
	}

	if sessionList[0] != "TEST123" {
		t.Errorf("Expected session ID 'TEST123', got %v", sessionList[0])
	}
}

func TestHandleSessions_POST(t *testing.T) {
	server := New()

	requestBody := map[string]string{"sessionId": "NEW123"}
	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/sessions", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleSessions)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Could not parse JSON response: %v", err)
	}

	if response["sessionId"] != "NEW123" {
		t.Errorf("Expected sessionId 'NEW123', got %s", response["sessionId"])
	}

	if response["status"] != "created" {
		t.Errorf("Expected status 'created', got %s", response["status"])
	}

	// Verify session was actually created
	server.mu.RLock()
	_, exists := server.sessions["NEW123"]
	server.mu.RUnlock()

	if !exists {
		t.Error("Expected session to be created in server sessions map")
	}
}

func TestHandleSessions_POST_InvalidJSON(t *testing.T) {
	server := New()

	req, err := http.NewRequest("POST", "/api/sessions", strings.NewReader("invalid json"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleSessions)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}
}

func TestHandleSessions_UnsupportedMethod(t *testing.T) {
	server := New()

	req, err := http.NewRequest("DELETE", "/api/sessions", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleSessions)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, status)
	}
}

func TestHandleSession_ExistingSession(t *testing.T) {
	server := New()

	// Create a test session
	session := poker.NewSession("EXISTING123")
	server.sessions["EXISTING123"] = session

	req, err := http.NewRequest("GET", "/api/sessions/EXISTING123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleSession)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Could not parse JSON response: %v", err)
	}

	if response["id"] != "EXISTING123" {
		t.Errorf("Expected session ID 'EXISTING123', got %v", response["id"])
	}
}

func TestHandleSession_NonExistentSession(t *testing.T) {
	server := New()

	req, err := http.NewRequest("GET", "/api/sessions/NONEXISTENT", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleSession)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, status)
	}
}

func TestHandleSessions_POST_ExistingSession(t *testing.T) {
	server := New()

	// Pre-create a session
	server.sessions["EXISTING123"] = poker.NewSession("EXISTING123")

	requestBody := map[string]string{"sessionId": "EXISTING123"}
	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "/api/sessions", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleSessions)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Could not parse JSON response: %v", err)
	}

	if response["sessionId"] != "EXISTING123" {
		t.Errorf("Expected sessionId 'EXISTING123', got %s", response["sessionId"])
	}

	if response["status"] != "created" {
		t.Errorf("Expected status 'created', got %s", response["status"])
	}
}

func TestConcurrentSessionCreation(t *testing.T) {
	server := New()

	// Test concurrent creation of the same session
	done := make(chan bool, 2)

	createSession := func() {
		requestBody := map[string]string{"sessionId": "CONCURRENT123"}
		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("POST", "/api/sessions", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.HandleSessions)

		handler.ServeHTTP(rr, req)
		done <- true
	}

	// Start two concurrent requests
	go createSession()
	go createSession()

	// Wait for both to complete
	<-done
	<-done

	// Verify only one session was created
	server.mu.RLock()
	_, exists := server.sessions["CONCURRENT123"]
	sessionCount := len(server.sessions)
	server.mu.RUnlock()

	if !exists {
		t.Error("Expected session to be created")
	}

	if sessionCount != 1 {
		t.Errorf("Expected 1 session, got %d", sessionCount)
	}
}

func TestSessionLifecycle(t *testing.T) {
	server := New()

	// Step 1: Create session via POST
	requestBody := map[string]string{"sessionId": "LIFECYCLE123"}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/sessions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleSessions)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Failed to create session: %d", rr.Code)
	}

	// Step 2: Verify it appears in sessions list
	req, _ = http.NewRequest("GET", "/api/sessions", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var listResponse map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &listResponse)
	sessions := listResponse["sessions"].([]interface{})

	found := false
	for _, sessionID := range sessions {
		if sessionID == "LIFECYCLE123" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Created session not found in sessions list")
	}

	// Step 3: Get specific session details
	req, _ = http.NewRequest("GET", "/api/sessions/LIFECYCLE123", nil)
	rr = httptest.NewRecorder()
	sessionHandler := http.HandlerFunc(server.HandleSession)
	sessionHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Failed to get session details: %d", rr.Code)
	}

	var sessionResponse map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &sessionResponse)

	if sessionResponse["id"] != "LIFECYCLE123" {
		t.Errorf("Expected session ID 'LIFECYCLE123', got %v", sessionResponse["id"])
	}

	if sessionResponse["status"] != "waiting" {
		t.Errorf("Expected session status 'waiting', got %v", sessionResponse["status"])
	}
}
