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

The application supports extensive configuration through environment variables:

### Server Configuration
- `PORT` - Server port (default: 8080)
- `HOST` - Server host (default: "")
- `READ_TIMEOUT` - Request read timeout (default: 15s)
- `WRITE_TIMEOUT` - Response write timeout (default: 15s)
- `IDLE_TIMEOUT` - Connection idle timeout (default: 60s)
- `SHUTDOWN_TIMEOUT` - Graceful shutdown timeout (default: 10s)

### Security Configuration
- `ALLOWED_ORIGINS` - Comma-separated list of allowed CORS origins (default: "*")
- `MAX_MESSAGE_SIZE` - Maximum WebSocket message size in bytes (default: 1024)

### Session Configuration
- `SESSION_TIMEOUT` - Session lifetime (default: 24h)
- `MAX_SESSIONS_PER_USER` - Maximum sessions per user (default: 10)

### Logging Configuration
- `LOG_LEVEL` - Log level: debug, info, warn, error (default: info)
- `LOG_FORMAT` - Log format: text, json (default: text)

### Development Configuration
- `DEVELOPMENT` - Enable development mode (default: false)
- `ENABLE_PPROF` - Enable pprof endpoints (default: false)

### Example Configuration
```bash
# Production configuration
export PORT=8080
export HOST=0.0.0.0
export ALLOWED_ORIGINS="https://yourapp.com,https://api.yourapp.com"
export LOG_LEVEL=warn
export DEVELOPMENT=false

# Development configuration
export PORT=3000
export DEVELOPMENT=true
export LOG_LEVEL=debug
```

### Health Check
The server provides a health check endpoint at `/health`:
```bash
curl http://localhost:8080/health
# Response: {"status":"ok","service":"planning-poker"}
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

## Monitoring CI/CD

This project includes GitHub Actions for continuous integration and automated releases. You can monitor the CI/CD pipeline status using the GitHub CLI.

### Setup GitHub CLI

1. **Install GitHub CLI**:
   ```bash
   # macOS
   brew install gh
   
   # Linux
   sudo apt install gh
   
   # Windows
   winget install GitHub.cli
   ```

2. **Authenticate with GitHub**:
   ```bash
   gh auth login
   ```

### Monitor GitHub Actions

Use the included monitoring scripts:

```bash
# Quick status check
./scripts/check-actions.sh

# Real-time monitoring
./scripts/monitor-actions.sh watch

# View latest logs
./scripts/monitor-actions.sh logs

# Custom status check (show last 3 runs)
./scripts/check-actions.sh 3
```

### Useful GitHub CLI Commands

```bash
# List recent workflow runs
gh run list --limit 10

# View specific run details
gh run view <run-id>

# Watch current runs in real-time
gh run watch

# List all workflows
gh workflow list

# Rerun a failed workflow
gh run rerun <run-id>

# List releases
gh release list

# View release details
gh release view <tag>
```

## Creating Releases

The project uses automated releases with comprehensive artifacts:

```bash
# Create and push a version tag
git tag v1.0.0
git push origin v1.0.0
```

The GitHub Actions workflow will automatically:
- ✅ Run all tests and quality checks
- ✅ Build documentation PDF from LaTeX sources
- ✅ Create binaries for multiple platforms (Linux, macOS, Windows)
- ✅ Create a GitHub release with all artifacts attached

**Release Artifacts Include**:
- `planning-poker-linux-amd64` - Linux binary
- `planning-poker-darwin-amd64` - macOS Intel binary  
- `planning-poker-darwin-arm64` - macOS Apple Silicon binary
- `planning-poker-windows-amd64.exe` - Windows binary
- `design.pdf` - Complete system design documentation

### Automated Documentation

The technical documentation is built automatically:
- **Source**: LaTeX files in `/docs` directory
- **Build**: Automatic during release process
- **Distribution**: PDF attached to GitHub releases
- **Benefits**: No binary files in git, always up-to-date documentation

## Contributing

We follow a structured development workflow to maintain code quality. Please read our [Contributing Guide](CONTRIBUTING.md) for detailed information.

### Quick Start for Contributors

1. **Clone and setup**:
   ```bash
   git clone https://github.com/meconlen/planning_poker.git
   cd planning_poker
   git checkout dev  # Always work on dev branch
   ```

2. **Make changes on dev branch**:
   ```bash
   # Make your changes
   go test ./...     # Run tests
   go build          # Verify build
   ```

3. **Commit and push**:
   ```bash
   git add .
   git commit -m "feat: Your feature description"
   git push origin dev
   ```

### Development Workflow

- **`dev`** - Main development branch (work here)
- **`main`** - Production-ready code only (merge from dev when stable)
- **Releases** - Tagged from main after thorough testing

### Key Rules
- ✅ All development happens on `dev` branch
- ✅ All tests must pass before merging to `main`
- ✅ `main` branch is protected and stable
- ✅ Releases are tagged from `main` only

For detailed workflow, testing requirements, and code standards, see [CONTRIBUTING.md](CONTRIBUTING.md).
