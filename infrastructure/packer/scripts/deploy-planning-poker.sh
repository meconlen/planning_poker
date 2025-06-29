#!/bin/bash

# Planning Poker Deployment Script
# This script downloads the latest Docker image and starts the application

set -e

# Configuration
GITHUB_REPO="meconlen/planning_poker"
CONTAINER_NAME="planning-poker"
IMAGE_NAME="planning-poker"
PORT="8080"
LOG_FILE="/var/log/planning-poker-deploy.log"

# Logging function
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

log "Starting Planning Poker deployment..."

# Stop and remove existing container if it exists
if docker ps -q -f name="$CONTAINER_NAME" | grep -q .; then
    log "Stopping existing container..."
    docker stop "$CONTAINER_NAME" || true
fi

if docker ps -aq -f name="$CONTAINER_NAME" | grep -q .; then
    log "Removing existing container..."
    docker rm "$CONTAINER_NAME" || true
fi

# Remove old images to save space (keep latest)
log "Cleaning up old Docker images..."
docker image prune -f || true

# Get the latest release tag
log "Fetching latest release information..."
LATEST_TAG=$(curl -s "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | jq -r '.tag_name')

if [ "$LATEST_TAG" = "null" ] || [ -z "$LATEST_TAG" ]; then
    log "ERROR: Could not fetch latest release tag"
    exit 1
fi

log "Latest release tag: $LATEST_TAG"

# Download the Docker image from GitHub release
DOCKER_IMAGE_URL="https://github.com/$GITHUB_REPO/releases/download/$LATEST_TAG/planning-poker-docker-$LATEST_TAG.tar.gz"
TEMP_DIR="/tmp/planning-poker-deploy"
DOCKER_IMAGE_FILE="$TEMP_DIR/planning-poker-docker-$LATEST_TAG.tar.gz"

# Create temp directory
mkdir -p "$TEMP_DIR"

log "Downloading Docker image: $DOCKER_IMAGE_URL"
if ! curl -L -f -o "$DOCKER_IMAGE_FILE" "$DOCKER_IMAGE_URL"; then
    log "ERROR: Failed to download Docker image"
    exit 1
fi

# Load the Docker image
log "Loading Docker image..."
if ! docker load < "$DOCKER_IMAGE_FILE"; then
    log "ERROR: Failed to load Docker image"
    exit 1
fi

# Tag the image for easier reference
docker tag "$IMAGE_NAME:$LATEST_TAG" "$IMAGE_NAME:latest"

# Start the new container
log "Starting Planning Poker container..."
docker run -d \
    --name "$CONTAINER_NAME" \
    --restart unless-stopped \
    -p "$PORT:8080" \
    -e PORT=8080 \
    "$IMAGE_NAME:$LATEST_TAG"

# Wait a moment and check if container is running
sleep 5

if docker ps -q -f name="$CONTAINER_NAME" | grep -q .; then
    log "✅ Planning Poker is now running on port $PORT"
    log "Container ID: $(docker ps -q -f name="$CONTAINER_NAME")"
    
    # Test the application
    if curl -f -s "http://localhost:$PORT/" > /dev/null; then
        log "✅ Application health check passed"
    else
        log "⚠️  Application health check failed, but container is running"
    fi
else
    log "❌ ERROR: Container failed to start"
    log "Container logs:"
    docker logs "$CONTAINER_NAME" 2>&1 | tail -20 | tee -a "$LOG_FILE"
    exit 1
fi

# Cleanup
log "Cleaning up temporary files..."
rm -rf "$TEMP_DIR"

log "Deployment completed successfully!"
log "Access Planning Poker at: http://$(curl -s ifconfig.me):$PORT"
