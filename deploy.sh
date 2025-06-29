#!/bin/bash

# Build and run the planning poker application with Docker

echo "🐳 Building Planning Poker Docker image..."
docker build -t planning-poker:latest .

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo "🚀 Starting Planning Poker container..."
    
    # Stop any existing container
    docker stop planning-poker 2>/dev/null || true
    docker rm planning-poker 2>/dev/null || true
    
    # Run the new container
    docker run -d \
        --name planning-poker \
        -p 8080:8080 \
        --restart unless-stopped \
        planning-poker:latest
    
    if [ $? -eq 0 ]; then
        echo "✅ Planning Poker is now running!"
        echo "🌐 Open http://localhost:8080 in your browser"
        echo "📋 Container logs: docker logs -f planning-poker"
        echo "🛑 Stop container: docker stop planning-poker"
    else
        echo "❌ Failed to start container"
        exit 1
    fi
else
    echo "❌ Build failed"
    exit 1
fi
