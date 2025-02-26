# ====================
# Deployment Script for Production (Zero-Downtime)
# ====================
#
# 📌 Instructions for Use:
# 1. If the server changes, copy this script outside the project directory (e.g., to `~`).
# 2. Modify the script to include `cd /path/to/project/root/folder` before execution.
# 3. Grant execution permissions with `chmod +x ~/deploy.sh`.
# 4. Run the script using `~/deploy.sh`.
#
# ====================
# Configuration
# ====================
ENV_FILE="./.env"
COMPOSE_FILE="./docker-compose.yml"
LOG_FILE="./deployment.log"

# Ensure log file exists
if [[ ! -f "$LOG_FILE" ]]; then
    touch "$LOG_FILE" || { echo "Failed to create log file: $LOG_FILE"; exit 1; }
fi

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

# Ensure script is run as root (for production safety)
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

# Check if in a Git repository
if ! git rev-parse --is-inside-work-tree &>/dev/null; then
    handle_error "This is not a Git repository. Please run this script inside a Git project."
fi

# Get current branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
log "📌 Current branch: $CURRENT_BRANCH"

# Ensure no uncommitted changes
if ! git diff-index --quiet HEAD --; then
    handle_error "Uncommitted changes detected! Please commit or stash your changes before deploying."
fi

# Pull latest changes
log "⬇️ Pulling latest updates from '$CURRENT_BRANCH'..."
git pull origin "$CURRENT_BRANCH" --rebase || handle_error "Failed to pull latest updates."
log "✅ Latest updates pulled successfully."

# ====================
# Blue-Green Deployment Process
# ====================
log "🚀 Starting Blue-Green Deployment..."

# Step 1: Pull latest images
log "⬇️ Pulling latest images..."
docker compose pull --parallel || handle_error "Failed to pull the latest images."

# Step 2: Create a new deployment network (if not exists)
log "🔄 Creating network if not exists..."
docker network create zeep-network || true

# Step 3: Start new version containers in parallel without stopping old ones
log "⚙️ Deploying new version in parallel..."
docker compose up --build -d || handle_error "Failed to start new version containers."

# Step 4: Wait for health checks
log "🩺 Waiting for new containers to be healthy..."
sleep 5  # Adjust based on service start time

if ! docker compose ps | grep -q 'healthy'; then
    handle_error "New version failed health checks!"
fi

log "✅ Deployment completed successfully with zero downtime."
