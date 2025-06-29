# Planning Poker

A real-time WebSocket-based planning poker application for agile teams.

## Getting Started

### Using Docker (Recommended)

1. Start the application:
   ```bash
   make dev
   ```

2. Open your browser to `http://localhost:8080`

3. Stop the application:
   ```bash
   make down
   ```

### Using Go directly

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

3. Open your browser to `http://localhost:8080`

## Features

- Real-time planning poker sessions
- WebSocket communication for instant updates
- Vote with Fibonacci sequence cards (0, 1, 2, 3, 5, 8, 13, 21, ?)
- Session-based voting with multiple users
- Reveal votes functionality
- Start new rounds
- Set story descriptions
- Docker containerization for easy deployment

## How to Use

1. Enter a session ID and your name to join a session
2. Set the current story being estimated
3. Cast your vote using the card values
4. Wait for all team members to vote
5. Reveal votes to see everyone's estimates
6. Start a new round for the next story

## Docker Commands

- `make dev` - Start development environment
- `make down` - Stop all services
- `make logs` - View application logs
- `make clean` - Clean up Docker resources
- `make build` - Build Docker image

## Development

This version includes basic planning poker functionality with Docker containerization for easy deployment and development.
