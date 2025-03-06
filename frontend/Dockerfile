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

# Remove default nginx config if needed
RUN rm -f /etc/nginx/conf.d/default.conf

# Copy the template
COPY nginx.template.conf /etc/nginx/nginx.template.conf

# Replace placeholder in the template
RUN envsubst '$NGINX_SERVER_NAME' < /etc/nginx/nginx.template.conf > /etc/nginx/nginx.conf

COPY --from=builder /app/dist /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
