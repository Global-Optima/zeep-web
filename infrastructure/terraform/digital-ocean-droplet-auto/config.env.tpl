# File: templates/config.env.tpl

# ==============================
# PostgreSQL Database
# ==============================
POSTGRES_HOST=${POSTGRES_HOST}
POSTGRES_PORT=${POSTGRES_PORT}
POSTGRES_USER=${POSTGRES_USER}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
POSTGRES_DB=${POSTGRES_DB}

# ==============================
# Redis Configuration
# ==============================
REDIS_USERNAME=${REDIS_USERNAME}
REDIS_HOST=${REDIS_HOST}
REDIS_PORT=${REDIS_PORT}
REDIS_PASSWORD=${REDIS_PASSWORD}

# ==============================
# Domain Configuration
# ==============================
DOMAIN=${DOMAIN}

# ==============================
# S3 Configuration
# ==============================
S3_ACCESS_KEY=${S3_ACCESS_KEY}
S3_SECRET_KEY=${S3_SECRET_KEY}
S3_ENDPOINT=${S3_ENDPOINT}
S3_BUCKET_NAME=${S3_BUCKET_NAME}

# ==============================
# JWT and Payment Secret Keys
# ==============================
JWT_CUSTOMER_SECRET_KEY=${JWT_CUSTOMER_SECRET_KEY}
JWT_EMPLOYEE_SECRET_KEY=${JWT_EMPLOYEE_SECRET_KEY}
PAYMENT_SECRET=${PAYMENT_SECRET}

# ==============================
# GitHub Configuration
# ==============================
GITHUB_TOKEN=${GITHUB_TOKEN}
GITHUB_REPO=${GITHUB_REPO}
GITHUB_BRANCH=${GITHUB_BRANCH}