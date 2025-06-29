#!/bin/bash

# Planning Poker Development Workflow Helper
# Usage: ./scripts/workflow.sh <command>

set -e

COMMAND=${1:-help}

case $COMMAND in
    "start"|"dev")
        echo "🌿 Starting development workflow..."
        echo "=================================="
        
        # Switch to dev branch
        echo "Switching to dev branch..."
        git checkout dev
        
        # Pull latest changes
        echo "Pulling latest changes..."
        git pull origin dev
        
        echo ""
        echo "✅ Ready for development on 'dev' branch!"
        echo "Make your changes and then run: ./scripts/workflow.sh test"
        ;;
        
    "test")
        echo "🧪 Running comprehensive tests..."
        echo "================================="
        
        # Verify we're on dev branch
        CURRENT_BRANCH=$(git branch --show-current)
        if [ "$CURRENT_BRANCH" != "dev" ]; then
            echo "⚠️  Warning: You're not on the 'dev' branch (currently on '$CURRENT_BRANCH')"
            echo "Switch to dev branch first: git checkout dev"
            exit 1
        fi
        
        echo "Running Go tests..."
        go test ./...
        
        echo ""
        echo "Running build test..."
        go build
        
        echo ""
        echo "Running go vet..."
        go vet ./...
        
        echo ""
        echo "✅ All tests passed!"
        echo "Ready to commit: ./scripts/workflow.sh commit \"Your commit message\""
        ;;
        
    "commit")
        if [ -z "$2" ]; then
            echo "❌ Error: Commit message required"
            echo "Usage: ./scripts/workflow.sh commit \"Your commit message\""
            exit 1
        fi
        
        echo "📝 Committing changes to dev branch..."
        echo "====================================="
        
        # Verify we're on dev branch
        CURRENT_BRANCH=$(git branch --show-current)
        if [ "$CURRENT_BRANCH" != "dev" ]; then
            echo "❌ Error: You must be on the 'dev' branch to commit"
            exit 1
        fi
        
        # Add all changes
        git add .
        
        # Commit with message
        git commit -m "$2"
        
        # Push to dev
        git push origin dev
        
        echo ""
        echo "✅ Changes committed and pushed to dev branch!"
        echo "Monitor CI status: ./scripts/monitor-actions.sh"
        ;;
        
    "release")
        echo "🚀 Starting release process..."
        echo "============================="
        
        # Verify we're on dev branch and it's clean
        CURRENT_BRANCH=$(git branch --show-current)
        if [ "$CURRENT_BRANCH" != "dev" ]; then
            echo "❌ Error: Start release from 'dev' branch"
            exit 1
        fi
        
        # Check for uncommitted changes
        if ! git diff-index --quiet HEAD --; then
            echo "❌ Error: You have uncommitted changes. Commit them first."
            exit 1
        fi
        
        # Run final tests
        echo "Running final tests on dev branch..."
        go test ./...
        if [ $? -ne 0 ]; then
            echo "❌ Error: Tests failed. Fix issues before release."
            exit 1
        fi
        
        # Switch to main
        echo "Switching to main branch..."
        git checkout main
        git pull origin main
        
        # Merge dev into main
        echo "Merging dev into main..."
        git merge --no-ff dev -m "Merge dev into main for release"
        
        # Final test on main
        echo "Running final verification on main..."
        go test ./...
        if [ $? -ne 0 ]; then
            echo "❌ Error: Tests failed on main branch!"
            echo "Rolling back merge..."
            git reset --hard HEAD~1
            git checkout dev
            exit 1
        fi
        
        # Push main
        git push origin main
        
        echo ""
        echo "✅ Successfully merged dev into main!"
        echo "Next steps:"
        echo "1. Monitor CI: ./scripts/monitor-actions.sh"
        echo "2. Tag release: ./scripts/workflow.sh tag v0.X.Y \"Release description\""
        echo "3. Return to dev: git checkout dev"
        ;;
        
    "tag")
        if [ -z "$2" ]; then
            echo "❌ Error: Version tag required"
            echo "Usage: ./scripts/workflow.sh tag v0.X.Y \"Release description\""
            exit 1
        fi
        
        VERSION=$2
        DESCRIPTION=${3:-"Release $VERSION"}
        
        echo "🏷️  Tagging release $VERSION..."
        echo "==============================="
        
        # Verify we're on main branch
        CURRENT_BRANCH=$(git branch --show-current)
        if [ "$CURRENT_BRANCH" != "main" ]; then
            echo "❌ Error: You must be on 'main' branch to tag a release"
            exit 1
        fi
        
        # Create annotated tag
        git tag -a "$VERSION" -m "$DESCRIPTION"
        
        # Push tag
        git push origin "$VERSION"
        
        echo ""
        echo "✅ Successfully tagged release $VERSION!"
        echo "GitHub Actions will now build and publish the release."
        echo ""
        echo "Return to development:"
        echo "git checkout dev"
        ;;
        
    "status")
        echo "📊 Development Workflow Status"
        echo "==============================="
        
        CURRENT_BRANCH=$(git branch --show-current)
        echo "Current branch: $CURRENT_BRANCH"
        
        if [ "$CURRENT_BRANCH" = "dev" ]; then
            echo "✅ On development branch - ready for development"
        elif [ "$CURRENT_BRANCH" = "main" ]; then
            echo "⚠️  On production branch - switch to dev for development"
        else
            echo "❓ On feature branch - merge to dev when ready"
        fi
        
        echo ""
        echo "Recent commits:"
        git log --oneline -5
        
        echo ""
        echo "Working directory status:"
        git status --porcelain
        ;;
        
    "help"|*)
        echo "🔧 Planning Poker Development Workflow"
        echo "======================================"
        echo ""
        echo "Commands:"
        echo "  start              Switch to dev branch and pull latest"
        echo "  test               Run all tests and verification"
        echo "  commit \"message\"   Commit and push changes to dev"
        echo "  release            Merge dev to main (production)"
        echo "  tag v0.X.Y \"desc\"  Tag a release version"
        echo "  status             Show current workflow status"
        echo ""
        echo "Typical workflow:"
        echo "1. ./scripts/workflow.sh start"
        echo "2. Make your changes"
        echo "3. ./scripts/workflow.sh test"
        echo "4. ./scripts/workflow.sh commit \"feat: Your changes\""
        echo "5. ./scripts/workflow.sh release  (when ready for production)"
        echo "6. ./scripts/workflow.sh tag v0.X.Y \"Release notes\""
        echo ""
        echo "Branch strategy:"
        echo "- dev  : All development work"
        echo "- main : Production-ready code only"
        echo "- tags : Released versions"
        ;;
esac
