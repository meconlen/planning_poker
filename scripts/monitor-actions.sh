#!/bin/bash

# Real-time GitHub Actions Monitor for Planning Poker Project
# Usage: ./scripts/monitor-actions.sh [watch|status|logs]

set -e

ACTION=${1:-status}

case $ACTION in
    "watch")
        echo "🔄 Watching GitHub Actions in real-time..."
        echo "Press Ctrl+C to stop"
        echo ""
        gh run watch
        ;;
    "logs")
        echo "📜 Getting logs from latest run..."
        LATEST_RUN=$(gh run list --limit 1 --json databaseId --jq '.[0].databaseId')
        if [ "$LATEST_RUN" != "null" ]; then
            echo "Showing logs for run ID: $LATEST_RUN"
            gh run view $LATEST_RUN --log
        else
            echo "No runs found"
        fi
        ;;
    "status"|*)
        echo "🚀 Planning Poker - GitHub Actions Monitor"
        echo "=========================================="
        echo ""
        
        # Show current running workflows
        echo "🔄 Currently Running:"
        RUNNING=$(gh run list --limit 10 --json status,workflowName,createdAt,databaseId --jq '.[] | select(.status == "in_progress") | "\(.databaseId) - \(.workflowName) (started \(.createdAt))"')
        if [ -z "$RUNNING" ]; then
            echo "   No workflows currently running"
        else
            echo "$RUNNING"
        fi
        
        echo ""
        echo "✅ Recent Successful Runs:"
        gh run list --limit 5 --json status,conclusion,workflowName,createdAt,headBranch --jq '.[] | select(.conclusion == "success") | "   \(.workflowName) on \(.headBranch) - \(.createdAt)"' | head -3
        
        echo ""
        echo "❌ Recent Failures (last 10 runs):"
        FAILURES=$(gh run list --limit 10 --json status,conclusion,workflowName,createdAt,databaseId --jq '.[] | select(.conclusion == "failure" or .conclusion == "cancelled") | "   \(.databaseId) - \(.workflowName) (\(.conclusion)) - \(.createdAt)"')
        if [ -z "$FAILURES" ]; then
            echo "   No recent failures! 🎉"
        else
            echo "$FAILURES" | head -3
        fi
        
        echo ""
        echo "🏷️  Latest Release:"
        gh release view --json tagName,name,publishedAt --jq '"   \(.name) (\(.tagName)) - \(.publishedAt)"' 2>/dev/null || echo "   No releases found"
        
        echo ""
        echo "💡 Quick Actions:"
        echo "   ./scripts/monitor-actions.sh watch    # Watch runs in real-time"
        echo "   ./scripts/monitor-actions.sh logs     # Show latest run logs"
        echo "   gh run list                           # List all runs"
        echo "   gh run rerun <run-id>                 # Rerun a failed workflow"
        ;;
esac
