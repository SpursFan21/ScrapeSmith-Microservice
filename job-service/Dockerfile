# Use official Node.js image as the base
FROM node:20-alpine AS builder

# Set working directory
WORKDIR /app

# Install build tools (for any native modules)
RUN apk add --no-cache build-base

# Copy package files first to leverage cache
COPY package.json package-lock.json ./

# Install all dependencies (not just production) so 'mongodb' is included
RUN npm install

# Copy the entire project (including src/)
COPY . .

# Set working directory to the source folder
WORKDIR /app/src

# Expose the port the service listens on
EXPOSE 3007

# Start the app
CMD ["node", "server.js"]
