# Contributing to Planning Poker

Thank you for your interest in contributing to the Planning Poker project! This document outlines our development workflow and guidelines.

## 🌿 Branch Strategy

We use a **Git Flow-inspired** branching strategy to maintain code quality and stability:

### Branch Structure

- **`main`** - Production-ready code only
  - Always stable and deployable
  - Protected branch with required reviews
  - All commits must pass CI/CD
  - Tagged with semantic versioning

- **`dev`** - Development integration branch
  - All feature development happens here
  - Integration testing occurs here
  - Must pass all tests before merging to main

- **`feature/*`** - Individual feature branches (optional)
  - Branch from `dev` for larger features
  - Merge back to `dev` when complete

## 🚀 Development Workflow

### 1. Starting New Development

```bash
# Switch to dev branch
git checkout dev

# Pull latest changes
git pull origin dev

# Start working on your changes
# (or create a feature branch for larger changes)
git checkout -b feature/your-feature-name
```

### 2. Making Changes

```bash
# Make your changes
# Edit code, add tests, update documentation

# Run tests locally
go test ./...

# Build to ensure compilation
go build

# Run linting
go vet ./...
staticcheck ./...
```

### 3. Committing Changes

```bash
# Add your changes
git add .

# Commit with descriptive message
git commit -m "feat: Add new voting feature

- Implement custom voting scales
- Add validation for vote values
- Update UI to support new scales
- Add comprehensive tests"
```

### 4. Testing and Integration

```bash
# Push to dev branch (or your feature branch)
git push origin dev

# Verify CI passes
./scripts/monitor-actions.sh

# Run integration tests
go run test/client.go localhost:8080
```

### 5. Promoting to Main (Release Process)

**Only when dev branch is stable and all tests pass:**

```bash
# Switch to main
git checkout main

# Pull latest
git pull origin main

# Merge dev (use --no-ff for merge commit)
git merge --no-ff dev -m "Merge dev into main for release"

# Verify all tests pass one final time
go test ./...

# Push to main
git push origin main

# Tag the release
git tag -a v0.X.Y -m "Release v0.X.Y: Brief description"
git push origin v0.X.Y

# Switch back to dev for continued development
git checkout dev
```

## 🧪 Testing Requirements

All code must be thoroughly tested before merging to main:

### Automated Tests
- **Unit Tests**: All new functions must have unit tests
- **Integration Tests**: Use the test client for WebSocket flows
- **CI/CD**: All GitHub Actions workflows must pass

### Manual Testing
- **Local Testing**: Test the application locally
- **Docker Testing**: Verify Docker builds work
- **Browser Testing**: Test the web interface manually

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test file
go test ./internal/poker/

# Run integration test client
go run test/client.go localhost:8080

# Check test coverage
go test -cover ./...
```

## 📝 Commit Message Guidelines

Use conventional commit format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or modifying tests
- `chore`: Maintenance tasks

### Examples
```bash
feat(voting): Add custom voting scales
fix(websocket): Resolve connection timeout issue
docs(readme): Update installation instructions
test(session): Add comprehensive session tests
```

## 🔧 Development Environment

### Prerequisites
- Go 1.24+
- Docker (optional)
- GitHub CLI (for monitoring)

### Setup
```bash
# Clone the repository
git clone https://github.com/meconlen/planning_poker.git
cd planning_poker

# Switch to dev branch
git checkout dev

# Install dependencies
go mod download

# Run the application
go run main.go

# Or use Docker
docker-compose up -d
```

## 🚨 Code Quality Standards

### Go Standards
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions
- Handle errors appropriately
- Use Go modules for dependencies

### Project Structure
```
planning-poker/
├── main.go                 # Entry point
├── internal/
│   ├── server/            # HTTP/WebSocket server
│   └── poker/             # Game logic
├── web/                   # Frontend files
├── test/                  # Integration tests
├── scripts/               # Utility scripts
├── .github/               # CI/CD workflows
└── docs/                  # Documentation
```

## 🔍 Monitoring and Debugging

### GitHub Actions Monitoring
```bash
# Quick status check
./scripts/check-actions.sh

# Real-time monitoring
./scripts/monitor-actions.sh watch

# View logs
./scripts/monitor-actions.sh logs
```

### Local Debugging
```bash
# Run with verbose logging
go run main.go

# Use the test client for WebSocket debugging
go run test/client.go localhost:8080

# Check Docker logs
docker logs -f planning-poker
```

## 🛡️ Security Considerations

- Never commit secrets or API keys
- Use environment variables for configuration
- Validate all user inputs
- Follow OWASP security guidelines
- Keep dependencies updated

## 📋 Pull Request Process

1. **Create from dev branch** (not main)
2. **Ensure all tests pass** locally and in CI
3. **Update documentation** if needed
4. **Use descriptive title and description**
5. **Request review** from maintainers
6. **Address feedback** promptly

## 🏷️ Release Process

Releases follow semantic versioning (SemVer):

- **Major** (X.0.0): Breaking changes
- **Minor** (0.X.0): New features, backward compatible
- **Patch** (0.0.X): Bug fixes, backward compatible

### Release Checklist
- [ ] All tests pass on dev branch
- [ ] Documentation is updated
- [ ] CHANGELOG is updated
- [ ] Version number is appropriate
- [ ] Integration tests pass
- [ ] Docker build succeeds
- [ ] Manual testing complete

## 🤝 Getting Help

- **Issues**: Use GitHub Issues for bug reports and feature requests
- **Discussions**: Use GitHub Discussions for questions
- **Documentation**: Check README.md and this file
- **Code**: Review existing code for patterns and examples

## 📞 Contact

For questions about contributing, please:
1. Check existing documentation
2. Search existing issues
3. Create a new issue with the `question` label

Happy coding! 🎉
