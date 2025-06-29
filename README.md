# Planning Poker

A real-time planning poker application built with Go and WebSockets.

## Features

- ✅ Real-time collaboration with WebSockets
- ✅ Multiple concurrent sessions
- ✅ Standard Fibonacci voting cards (0, ½, 1, 2, 3, 5, 8, 13, 21, ?, ☕)
- ✅ Vote hiding/revealing functionality
- ✅ Story management
- ✅ User presence tracking
- ✅ Clean, modern web interface

## Quick Start

### Option 1: Run with Go (Development)

1. **Install Go** (if not already installed):
   ```bash
   brew install go
   ```

2. **Clone/Navigate to the project**:
   ```bash
   cd /path/to/planning-poker
   ```

3. **Install dependencies**:
   ```bash
   go mod tidy
   ```

4. **Run the server**:
   ```bash
   go run main.go
   ```

### Option 2: Run with Docker (Recommended)

1. **Build and run with Docker**:
   ```bash
   # Using the deployment script
   ./deploy.sh
   
   # Or manually
   docker build -t planning-poker .
   docker run -d -p 8080:8080 --name planning-poker planning-poker
   ```

2. **Or use Docker Compose**:
   ```bash
   # Development environment
   docker-compose up -d
   
   # Production environment
   docker-compose -f docker-compose.prod.yml up -d
   ```

3. **Or use Make commands**:
   ```bash
   make run    # Build and run container
   make dev    # Start development environment
   make prod   # Start production environment
   make logs   # View logs
   make stop   # Stop containers
   ```

### Access the Application

Open your browser and go to:
```
http://localhost:8080
```

## Usage

1. Enter a session ID and your name to join a planning poker session
2. Set the current user story to be estimated
3. Each team member selects their estimate using the voting cards
4. Click "Reveal Votes" to show everyone's estimates
5. Click "New Round" to start estimating the next story

## Project Structure

```
planning-poker/
├── main.go                 # Main server entry point
├── internal/
│   ├── server/
│   │   └── server.go       # HTTP and WebSocket handlers
│   └── poker/
│       └── session.go      # Planning poker game logic
├── web/
│   └── index.html          # Frontend interface
├── go.mod                  # Go module file
└── README.md              # This file
```

## API Endpoints

- `GET /` - Serves the web interface
- `GET /ws` - WebSocket endpoint for real-time communication
- `GET /api/sessions` - List all active sessions
- `POST /api/sessions` - Create a new session
- `GET /api/sessions/{id}` - Get session state

## WebSocket Messages

The application uses JSON messages over WebSockets:

### Client to Server:
- `vote` - Submit a vote
- `reveal` - Reveal all votes
- `new_round` - Start a new voting round
- `set_story` - Set the current story

### Server to Client:
- `session_state` - Current session state
- `user_joined` - User joined notification
- `user_left` - User left notification

## Development

To add new features or modify the application:

1. **Backend**: Edit files in `internal/` directory
2. **Frontend**: Modify `web/index.html`
3. **Add dependencies**: Use `go get <package>`

## Configuration

The server runs on port 8080 by default. To change this, modify the port in `main.go`:

```go
log.Fatal(http.ListenAndServe(":8080", nil))
```

## Docker Commands

### Basic Docker Commands
```bash
# Build the image
docker build -t planning-poker .

# Run the container
docker run -d -p 8080:8080 --name planning-poker planning-poker

# View logs
docker logs -f planning-poker

# Stop the container
docker stop planning-poker

# Remove the container
docker rm planning-poker
```

### Docker Compose Commands
```bash
# Start development environment
docker-compose up -d

# Start production environment
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose logs -f

# Stop and remove
docker-compose down
```

### Make Commands (Recommended)
```bash
make help    # Show all available commands
make run     # Build and run container
make dev     # Start development environment
make prod    # Start production environment
make logs    # Show container logs
make stop    # Stop all containers
make clean   # Clean up containers and images
make shell   # Open shell in running container
```

## Production Deployment

For production deployment, consider:

1. **Environment Variables**: Add configuration for port, CORS origins, etc.
2. **HTTPS**: Use TLS certificates for secure WebSocket connections
3. **Database**: Add persistence for sessions and historical data
4. **Authentication**: Add user authentication and authorization
5. **Rate Limiting**: Implement rate limiting for API endpoints
6. **Monitoring**: Add logging and metrics collection

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test the application
5. Submit a pull request
