# ====================
# Build Stage
# ====================
FROM oven/bun:latest AS builder

# Set the working directory
WORKDIR /app

# Build arguments for passing env variables at build time
ARG VITE_API_URL
ARG VITE_WS_URL
ARG VITE_PAYMENT_SECRET
ARG VITE_TEST_PAYMENT
ARG VITE_SAVE_ON_PRINT

# Copy only package files first for better caching
COPY package*.json ./

# Install dependencies using bun
RUN bun install

# Copy the rest of the project files
COPY . .

# Inject environment variables for Vite build
ENV VITE_API_URL=$VITE_API_URL
ENV VITE_WS_URL=$VITE_WS_URL
ENV VITE_PAYMENT_SECRET=$VITE_PAYMENT_SECRET
ENV VITE_TEST_PAYMENT=$VITE_TEST_PAYMENT
ENV VITE_SAVE_ON_PRINT=$VITE_SAVE_ON_PRINT


# Build the production files
RUN bun run build

# ====================
# Final Stage (Nginx)
# ====================
FROM nginx:alpine AS final

# Remove default config and copy template
RUN rm -f /etc/nginx/conf.d/default.conf
COPY nginx.template.conf /etc/nginx/nginx.template.conf

# Copy static files
COPY --from=builder /app/dist /usr/share/nginx/html

# Add custom entrypoint script
COPY scripts/envsubst_nginx.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose port and set entrypoint
EXPOSE 80
ENTRYPOINT ["/entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]