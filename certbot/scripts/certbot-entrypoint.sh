#!/bin/sh

# Ensure required environment variables are set
if [ -z "$NGINX_SERVER_NAME" ]; then
  echo "❌ ERROR: NGINX_SERVER_NAME is not set in .env file."
  exit 1
fi

if [ -z "$CERTBOT_EMAIL" ]; then
  echo "❌ ERROR: CERTBOT_EMAIL is not set in .env file."
  exit 1
fi

echo "🚀 Generating SSL certificates for: $NGINX_SERVER_NAME"

mkdir -p /var/www/certbot

# Convert space-separated domains into -d flags for certbot
DOMAIN_ARGS=""
for domain in $NGINX_SERVER_NAME; do
  DOMAIN_ARGS="$DOMAIN_ARGS -d $domain"
done

# Request SSL certificates
if certbot certonly --webroot -w /var/www/certbot --email "$CERTBOT_EMAIL" --agree-tos --no-eff-email --force-renewal $DOMAIN_ARGS; then
  echo "✅ SSL Certificates successfully generated!"
else
  echo "❌ ERROR: Certbot failed to generate SSL certificates."
  echo "❗ Check logs for details: /var/log/letsencrypt/letsencrypt.log"
  exit 1
fi

# Auto-renew SSL and reload Nginx only if Certbot was successful
echo "🔄 Starting Certbot Auto-Renewal..."
exec certbot renew --quiet --deploy-hook "nginx -s reload"
