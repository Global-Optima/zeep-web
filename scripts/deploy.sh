#!/bin/bash

# ====================
# Configuration
# ====================
ENV_FILE="/root/.env"
COMPOSE_FILE="/root/docker-compose.yml"
LOG_FILE="/root/scripts/deployment.log"

# ====================
# Secure Script Execution
# ====================
set -euo pipefail  # Fail on errors, undefined variables, and pipe failures

# Logging function (logs to both console and file)
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE" >&3
}

# Error handling function
handle_error() {
    log "🚨 ERROR: $1"
    log "❌ Deployment failed. Exiting..."
    exit 1
}

# ====================
# Pre-Deployment Checks
# ====================
exec 3>&1  # Store original stdout for logging

log "🔍 Starting pre-deployment checks..."

# Ensure Docker is installed
if ! command -v docker &> /dev/null; then
    handle_error "Docker is not installed. Please install Docker and try again."
fi

# Ensure Docker Compose (new CLI) is installed
if ! command -v docker compose &> /dev/null; then
    handle_error "Docker Compose is not installed. Please install it using: 'apt install docker-compose-plugin'"
fi

# Ensure the script is running as root (for production safety)
if [[ $EUID -ne 0 ]]; then
    handle_error "This script must be run as root. Try using 'sudo ./deploy.sh'."
fi

# Ensure .env file exists
if [[ ! -f "$ENV_FILE" ]]; then
    handle_error "Missing .env file at $ENV_FILE. Please create a valid .env file and try again."
fi

# Ensure docker-compose.yml exists
if [[ ! -f "$COMPOSE_FILE" ]]; then
    handle_error "Missing docker-compose.yml at $COMPOSE_FILE. Please create it and try again."
fi

log "✅ Pre-deployment checks completed successfully."

# ====================
# Deployment Process
# ====================
log "🚀 Starting deployment process..."

# Step 1: Pull the latest images (optional but recommended for production)
log "⬇️ Pulling the latest images..."
docker compose --env-file "$ENV_FILE" pull --parallel || handle_error "Failed to pull the latest images."

# Step 2: Stop and remove existing containers
log "🛑 Stopping and removing existing containers..."
docker compose --env-file "$ENV_FILE" down --remove-orphans || handle_error "Failed to stop and remove existing containers."

# Step 3: Build and start new containers
log "⚙️ Building and starting new containers..."
docker compose --env-file "$ENV_FILE" up --build -d --force-recreate || handle_error "Failed to build and start new containers."

# Step 4: Verify container status
log "🔍 Verifying container status..."
if ! docker ps --format "table {{.Names}}\t{{.Status}}"; then
    handle_error "Failed to verify container status."
fi

log "✅ Deployment completed successfully."

# ====================
# Post-Deployment Cleanup (Safe & Efficient)
# ====================
log "🧹 Performing post-deployment cleanup..."

# Check if cleanup is necessary before running it
if docker system df | grep -q 'Reclaimable'; then
    log "♻️ Removing unused Docker resources..."
    docker system prune -f --volumes || handle_error "Failed to clean up unused resources."
else
    log "🟢 No cleanup needed."
fi

log "🏁 Post-deployment cleanup completed successfully."
log "🚀 Deployment finished successfully!"
