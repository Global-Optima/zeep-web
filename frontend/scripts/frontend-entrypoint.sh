#!/bin/sh
set -e

# Generate Nginx config from template
envsubst '${NGINX_SERVER_NAME}' < /etc/nginx/nginx.template.conf > /etc/nginx/nginx.conf

# Check for certificate renewal
if [ -d "/etc/letsencrypt/live/${NGINX_SERVER_NAME}" ]; then
  nginx -t && nginx -s reload
fi

# Start Nginx
exec nginx -g "daemon off;"