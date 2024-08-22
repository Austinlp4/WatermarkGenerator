# Build frontend
FROM node:14-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install --no-cache
COPY frontend ./
# Add a build argument for cache busting
ARG CACHE_BUST=1
# Use the build argument in a RUN command to force a rebuild
RUN echo "Cache bust: ${CACHE_BUST}"
RUN npm run build

# Build backend
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
# Add a build argument for cache busting
ARG CACHE_BUST=1
# Use the build argument in a RUN command to force a rebuild
RUN echo "Cache bust: ${CACHE_BUST}"
COPY . .
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
RUN go mod download
RUN go build -o watermark-generator

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=backend-builder /app/watermark-generator .
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
EXPOSE 8080
CMD ["./watermark-generator"]