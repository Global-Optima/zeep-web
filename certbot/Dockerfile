FROM certbot/certbot:latest

# Copy the certbot entrypoint script
COPY scripts/entrypoint.sh /entrypoint.sh

# Make script executable
RUN chmod +x /entrypoint.sh

# Set entrypoint to the script
ENTRYPOINT ["/entrypoint.sh"]
