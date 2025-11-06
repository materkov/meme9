# Stage 1: Build
FROM node:20-alpine AS build
WORKDIR /app

# Install build dependencies for native modules
RUN apk add --no-cache python3 make g++

# Install dependencies (including optional dependencies for native modules)
COPY package.json package-lock.json* ./
RUN npm install --include=optional

# Copy source
COPY . .

# Build static assets
RUN npm run build

# Stage 2: Serve static files with nginx
FROM nginx:1.27-alpine

# Copy built files to nginx html directory
COPY --from=build /app/dist /usr/share/nginx/html
# Nginx config to listen on 8081
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose default nginx port
EXPOSE 8081

# Use the default nginx configuration
CMD ["nginx", "-g", "daemon off;"]
