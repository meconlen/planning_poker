#!/bin/bash

# GitHub Actions Status Checker for Planning Poker Project
# Usage: ./scripts/check-actions.sh [number_of_runs]

set -e

echo "🔍 Planning Poker - GitHub Actions Status"
echo "=========================================="

# Number of runs to show (default: 5)
LIMIT=${1:-5}

echo ""
echo "📋 Recent Workflow Runs (last $LIMIT):"
echo "--------------------------------------"
gh run list --limit $LIMIT

echo ""
echo "🏷️  Recent Releases:"
echo "-------------------"
gh release list --limit 3

echo ""
echo "📊 Current Repository Status:"
echo "-----------------------------"
echo "Repository: $(gh repo view --json nameWithOwner --jq .nameWithOwner)"
echo "Default Branch: $(gh repo view --json defaultBranchRef --jq .defaultBranchRef.name)"

# Check if there are any failed runs in the last 10
echo ""
echo "⚠️  Checking for recent failures..."
FAILED_RUNS=$(gh run list --limit 10 --json status,conclusion,workflowName,createdAt --jq '.[] | select(.conclusion == "failure" or .conclusion == "cancelled") | "\(.workflowName) - \(.conclusion) (\(.createdAt))"')

if [ -z "$FAILED_RUNS" ]; then
    echo "✅ No failed runs in the last 10 workflow executions!"
else
    echo "❌ Recent failures found:"
    echo "$FAILED_RUNS"
fi

echo ""
echo "💡 Useful commands:"
echo "  gh run list                    # List recent runs"
echo "  gh run view <run-id>           # View run details"
echo "  gh run watch                   # Watch current runs"
echo "  gh release list               # List releases"
echo "  gh workflow list               # List workflows"
