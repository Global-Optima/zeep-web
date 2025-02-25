#!/bin/bash
# ====================
# Configuration
# ====================
ENV_FILE="./.env"
COMPOSE_FILE="./docker-compose.yml"
LOG_FILE="./deployment.log"

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

# ====================
# Pull Latest Updates
# ====================
log "🔄 Checking for the latest updates..."

# Check if the current directory is a Git repository
if ! git rev-parse --is-inside-work-tree &>/dev/null; then
    handle_error "This is not a Git repository. Please run this script inside a Git project."
fi

# Get the current branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
log "📌 Current branch: $CURRENT_BRANCH"

# Ensure there are no uncommitted changes
if ! git diff-index --quiet HEAD --; then
    handle_error "Uncommitted changes detected! Please commit or stash your changes before deploying."
fi

# Pull the latest changes from the current branch
log "⬇️ Pulling latest updates from '$CURRENT_BRANCH'..."
git pull origin "$CURRENT_BRANCH" --rebase || handle_error "Failed to pull latest updates."

log "✅ Latest updates pulled successfully."

# ====================
# Deployment Process
# ====================
log "🚀 Starting deployment process..."

# Step 1: Pull the latest images (optional but recommended for production)
log "⬇️ Pulling the latest images..."
docker compose pull --parallel || handle_error "Failed to pull the latest images."

# Step 2: Stop and remove existing containers
log "🛑 Stopping and removing existing containers..."
docker compose down --remove-orphans || handle_error "Failed to stop and remove existing containers."

# Step 3: Build and start new containers
log "⚙️ Building and starting new containers..."
docker compose up --build -d --force-recreate || handle_error "Failed to build and start new containers."

# Step 4: Verify container status
log "🔍 Verifying container status..."
docker compose ps --format "table {{.Service}}\t{{.State}}\t{{.Health}}" || handle_error "Failed to verify container status."

# Ensure all required services are healthy
if ! docker compose ps | grep -q 'healthy'; then
    handle_error "One or more services failed to start properly."
fi

log "✅ Deployment completed successfully."
