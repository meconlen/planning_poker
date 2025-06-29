# Planning Poker

A real-time WebSocket-based planning poker application for agile teams.

## Getting Started

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

## How to Use

1. Enter a session ID and your name to join a session
2. Set the current story being estimated
3. Cast your vote using the card values
4. Wait for all team members to vote
5. Reveal votes to see everyone's estimates
6. Start a new round for the next story

## Development

This version includes basic planning poker functionality with real-time updates.
