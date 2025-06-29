# Development Workflow Summary

This document provides a quick reference for the Planning Poker development workflow.

## 🌿 Branch Strategy

- **`dev`** - Main development branch (work here daily)
- **`main`** - Production-ready code only (merge from dev when stable)

## 🚀 Quick Commands

```bash
# Start development
./scripts/workflow.sh start

# Test your changes
./scripts/workflow.sh test

# Commit to dev branch
./scripts/workflow.sh commit "feat: Your feature description"

# Monitor CI status
./scripts/monitor-actions.sh

# When ready for release
./scripts/workflow.sh release
./scripts/workflow.sh tag v0.X.Y "Release description"
```

## 📋 Development Checklist

### Before You Start
- [ ] Switch to dev branch: `./scripts/workflow.sh start`
- [ ] Pull latest changes
- [ ] Verify CI is working

### During Development
- [ ] Make changes on `dev` branch only
- [ ] Run tests frequently: `./scripts/workflow.sh test`
- [ ] Commit regularly with good messages
- [ ] Monitor CI status after each push

### Before Release
- [ ] All tests pass: `go test ./...`
- [ ] Code builds successfully: `go build`
- [ ] Integration tests pass: `go run test/client.go localhost:8080`
- [ ] CI is green on dev branch
- [ ] Documentation is updated

### Release Process
- [ ] Merge dev to main: `./scripts/workflow.sh release`
- [ ] Verify CI passes on main
- [ ] Tag the release: `./scripts/workflow.sh tag v0.X.Y "Description"`
- [ ] Return to dev branch: `git checkout dev`

## 🔧 Tool Reference

### Workflow Tools
```bash
./scripts/workflow.sh help      # Show all commands
./scripts/workflow.sh status    # Check current status
./scripts/workflow.sh start     # Switch to dev and pull
./scripts/workflow.sh test      # Run all tests
./scripts/workflow.sh commit    # Commit and push to dev
./scripts/workflow.sh release   # Merge dev to main
./scripts/workflow.sh tag       # Tag a release
```

### Monitoring Tools
```bash
./scripts/monitor-actions.sh          # Show CI status
./scripts/monitor-actions.sh watch    # Real-time monitoring
./scripts/monitor-actions.sh logs     # View latest logs
./scripts/check-actions.sh           # Quick status check
```

### GitHub CLI
```bash
gh run list              # List workflow runs
gh run view <id>         # View run details
gh workflow list         # List workflows
gh release list          # List releases
```

## 🛡️ Protection Rules

- **main branch** is protected
- All merges to main must pass CI
- Direct pushes to main are discouraged
- All releases are tagged from main
- Development happens on dev branch

## 📖 Full Documentation

- [CONTRIBUTING.md](CONTRIBUTING.md) - Complete development guide
- [README.md](README.md) - Project overview and setup
- [.github/workflows/](/.github/workflows/) - CI/CD configuration

## 🎯 Best Practices

1. **Always work on dev branch** for day-to-day development
2. **Test before committing** using the workflow script
3. **Use semantic commit messages** (feat:, fix:, docs:, etc.)
4. **Monitor CI status** after each push
5. **Only merge to main** when dev is stable and tested
6. **Tag releases** immediately after merging to main
7. **Return to dev** after releases for continued development

This workflow ensures code quality, maintains a stable main branch, and provides clear release management.
