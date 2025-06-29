#!/bin/bash

# Startup script for Planning Poker Terraform deployment
# This script runs when the instance first boots

set -e

# Log everything
exec > >(tee /var/log/planning-poker-startup.log) 2>&1

echo "[$(date)] Starting Planning Poker startup script..."

# Wait for system to be ready
sleep 30

# Ensure Docker is running
systemctl start docker
systemctl enable docker

# Start the Planning Poker service
echo "[$(date)] Starting Planning Poker service..."
systemctl start planning-poker

# Check status
if systemctl is-active --quiet planning-poker; then
    echo "[$(date)] ✅ Planning Poker service started successfully"
    
    # Wait a bit and test the application
    sleep 10
    if curl -f -s "http://localhost:8080/" > /dev/null; then
        echo "[$(date)] ✅ Planning Poker application is responding"
    else
        echo "[$(date)] ⚠️  Planning Poker application not responding yet"
    fi
else
    echo "[$(date)] ❌ Failed to start Planning Poker service"
    systemctl status planning-poker --no-pager
fi

echo "[$(date)] Startup script completed"
