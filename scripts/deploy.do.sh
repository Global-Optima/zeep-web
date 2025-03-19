#!/usr/bin/env bash
set -eux

# ------------------------------------------------------------------------------
# 1) Load Variables from config.env
#    Make sure config.env is in the same directory or specify an absolute path.
# ------------------------------------------------------------------------------
if [[ ! -f "config.env" ]]; then
  echo "Error: config.env not found in current directory!"
  exit 1
fi

# Load environment variables into current shell
source config.env

# ------------------------------------------------------------------------------
# 2) Remove Old Docker Packages (if any)
# ------------------------------------------------------------------------------
for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do
  sudo apt-get remove -y "$pkg" || true
done

# ------------------------------------------------------------------------------
# 3) Install Docker + Docker Compose (Official Instructions)
#    Reference: https://docs.docker.com/engine/install/ubuntu/
# ------------------------------------------------------------------------------
sudo apt-get update
sudo apt-get install -y ca-certificates curl

sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] \
  https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" \
  | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Enable and start Docker
sudo systemctl enable docker
sudo systemctl start docker

# ------------------------------------------------------------------------------
# 4) Clone Private GitHub Repository
#    We use $GITHUB_TOKEN, $GITHUB_REPO, $GITHUB_BRANCH from config.env
# ------------------------------------------------------------------------------
APP_DIR="/opt/zeep-app"
sudo mkdir -p "$APP_DIR"
sudo chown "$(whoami):$(whoami)" "$APP_DIR"
cd "$APP_DIR"

# Clone the repo
if [[ -d "app" ]]; then
  echo "Directory 'app' already exists. Skipping clone or pulling latest changes."
  # Optionally you can pull to update an existing clone:
  # cd app && git pull
else
  git clone "https://${GITHUB_TOKEN}@github.com/${GITHUB_REPO}.git" app
fi

cd app
git checkout "${GITHUB_BRANCH}"

# Clear the GitHub token from shell history
history -c || true

# ------------------------------------------------------------------------------
# 5) Generate .env from Variables
# ------------------------------------------------------------------------------
cat <<EOF > .env
# ==============================
# Backend Configuration
# ==============================
BACKEND_ENV=production
SERVER_PORT=8080

# ==============================
# PostgreSQL Database
# ==============================
DB_HOST=${POSTGRES_HOST}
DB_PORT=${POSTGRES_PORT}
DB_USER=${POSTGRES_USER}
DB_PASSWORD=${POSTGRES_PASSWORD}
DB_NAME=${POSTGRES_DB}
DB_SSL_MODE=require

# ==============================
# Redis Configuration
# ==============================
REDIS_USERNAME=${REDIS_USERNAME}
REDIS_HOST=${REDIS_HOST}
REDIS_PORT=${REDIS_PORT}
REDIS_PASSWORD=${REDIS_PASSWORD}
REDIS_DB=0
REDIS_ENABLE_TLS=true

# ==============================
# Payment Configuration
# ==============================
PAYMENT_SECRET=${PAYMENT_SECRET}

# ==============================
# Frontend Configuration
# ==============================
CLIENT_URL=https://${DOMAIN}
VITE_API_URL=https://${DOMAIN}
VITE_WS_URL=wss://${DOMAIN}/api/v1
VITE_TEST_PAYMENT=false
NGINX_SERVER_NAME=${DOMAIN}

# ==============================
# Auth Configuration
# ==============================
JWT_CUSTOMER_SECRET_KEY=${JWT_CUSTOMER_SECRET_KEY}
JWT_EMPLOYEE_SECRET_KEY=${JWT_EMPLOYEE_SECRET_KEY}
JWT_CUSTOMER_TOKEN_TTL=168h
JWT_EMPLOYEE_TOKEN_TTL=168h

# ==============================
# S3 Configuration
# ==============================
S3_ACCESS_KEY=${S3_ACCESS_KEY}
S3_SECRET_KEY=${S3_SECRET_KEY}
S3_ENDPOINT=https://${S3_ENDPOINT}
S3_BUCKET_NAME=${S3_BUCKET_NAME}

DEFAULT_PAGE=1
DEFAULT_PAGE_SIZE=10
MAX_PAGE_SIZE=100
EOF

# ------------------------------------------------------------------------------
# 6) Launch Services with Docker Compose
# ------------------------------------------------------------------------------
sudo docker compose -f docker-compose.do.yml up -d --build --force-recreate
