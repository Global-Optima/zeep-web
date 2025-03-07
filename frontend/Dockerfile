# ====================
# Build Stage
# ====================
FROM oven/bun:latest AS builder

# Set the working directory
WORKDIR /app

# Build arguments for passing env variables at build time
ARG VITE_API_URL
ARG VITE_WS_URL

# Copy only package files first for better caching
COPY package*.json ./

# Install dependencies using bun
RUN bun install

# Copy the rest of the project files
COPY . .

# Inject environment variables for Vite build
ENV VITE_API_URL=$VITE_API_URL
ENV VITE_WS_URL=$VITE_WS_URL

# Build the production files
RUN bun run build

# ====================
# Final Stage (Nginx)
# ====================
FROM nginx:alpine AS final

# Remove default config and copy template
RUN rm -f /etc/nginx/conf.d/default.conf
COPY nginx.template.conf /etc/nginx/nginx.template.conf

# Copy static frontend files
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy entrypoint script
COPY scripts/frontend-entrypoint.sh /frontend-entrypoint.sh
RUN chmod +x /frontend-entrypoint.sh

EXPOSE 80 443

ENTRYPOINT ["/frontend-entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]
