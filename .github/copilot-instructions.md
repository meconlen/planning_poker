<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# Planning Poker Project - Copilot Instructions

This is a real-time planning poker application built with Go, WebSockets, and Docker.

## Project Structure
- `main.go` - Main server entry point with HTTP handlers
- `internal/server/` - HTTP and WebSocket server logic
- `internal/poker/` - Planning poker game logic and session management
- `web/` - Frontend HTML/CSS/JavaScript files
- Docker files for containerization

## Key Technologies
- **Backend**: Go 1.24+ with Gorilla WebSocket library
- **Frontend**: Vanilla HTML/CSS/JavaScript with WebSocket client
- **Containerization**: Docker and Docker Compose
- **Architecture**: Clean architecture with separate packages

## Development Guidelines
- Use Go modules for dependency management
- Follow Go naming conventions and idioms
- Use JSON for WebSocket message formatting
- Implement proper error handling and logging
- Use mutex locks for concurrent access to shared data
- Follow REST API conventions for HTTP endpoints

## WebSocket Message Types
- `vote` - User submits a vote
- `reveal` - Reveal all votes in the session
- `new_round` - Start a new voting round
- `set_story` - Set the current story being estimated
- `session_state` - Current session state broadcast
- `user_joined`/`user_left` - User presence notifications

## Docker Deployment
- Use multi-stage builds for optimized images
- Support environment variables for configuration
- Include health checks in containers
- Use Docker Compose for easy orchestration

## Code Style
- Use `gofmt` for code formatting
- Include comprehensive error handling
- Add meaningful log messages
- Use struct tags for JSON serialization
- Implement graceful shutdowns for long-running processes
