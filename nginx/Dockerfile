# Nginx Dockerfile

# Build Frontend Stage
FROM node:18-alpine AS frontend_builder

WORKDIR /app

# Install dependencies
COPY frontend/package*.json ./
RUN npm install

# Copy frontend source code
COPY frontend/ .

# Build the frontend
RUN npm run build

# Nginx Stage
FROM nginx:stable-alpine

# Copy main Nginx configuration
COPY nginx/nginx.conf /etc/nginx/nginx.conf

# Copy default server configuration
COPY nginx/default.conf /etc/nginx/conf.d/default.conf

# Copy frontend build artifacts to Nginx html directory
COPY --from=frontend_builder /app/dist /usr/share/nginx/html

# Print out configuration files to verify their contents in the build stage
RUN echo "Contents of /etc/nginx/nginx.conf:" && cat /etc/nginx/nginx.conf
RUN echo "Contents of /etc/nginx/conf.d/default.conf:" && cat /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
