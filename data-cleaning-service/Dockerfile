# Use official Node.js image as the base
FROM node:20-alpine AS builder
WORKDIR /app

# Install build tools (if necessary)
RUN apk add --no-cache build-base

# Copy package.json and package-lock.json first to leverage Docker layer caching
COPY package.json package-lock.json ./

# Install dependencies
RUN npm install --production

# Copy source code
COPY . .

# Expose the port the service will run on
EXPOSE 3004

# Start the application
CMD ["node", "src/app.js"]
