# Build frontend
FROM node:14-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend ./
RUN npm run build

# Build backend
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
RUN go mod download
RUN go build -o watermark-generator

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=1 /app/watermark-generator .
COPY --from=0 /app/frontend/dist ./frontend/dist
EXPOSE 8080
CMD ["./watermark-generator"]