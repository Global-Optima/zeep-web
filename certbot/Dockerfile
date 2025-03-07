FROM certbot/certbot:latest

# Copy the certbot entrypoint script
COPY scripts/certbot-entrypoint.sh /certbot-entrypoint.sh

# Ensure script is executable
RUN chmod +x /certbot-entrypoint.sh

# Set entrypoint to the script
ENTRYPOINT ["/certbot-entrypoint.sh"]