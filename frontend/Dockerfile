# ====================
# Build Stage
# ====================
FROM node:lts-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy only package files first for better caching
COPY package*.json ./

# Use npm ci for clean & reproducible builds
RUN npm ci

# Copy the rest of the project files
COPY . .

# Build the production files
RUN npm run build

# ====================
# Final Stage
# ====================
FROM nginx:alpine AS final

# Remove default nginx config if you have your own
RUN rm /etc/nginx/conf.d/default.conf

# Copy custom nginx configuration if needed
COPY nginx.conf /etc/nginx/nginx.conf

# Copy the production build from the builder stage
COPY --from=builder /app/dist /usr/share/nginx/html


# Expose port 80
EXPOSE 80

# Run nginx in the foreground
CMD ["nginx", "-g", "daemon off;"]
