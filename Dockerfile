# --- Stage 1: Build the application ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files and download them to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application as a static binary
# Make sure to point to your main package, e.g., ./cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd/main.go


# --- Stage 2: Create the final, minimal image ---
FROM alpine:latest

# It's good practice to run as a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /home/appuser

# Copy ONLY the compiled binary from the builder stage
COPY --from=builder /app/main .

EXPOSE 8080

# Run the binary
ENTRYPOINT ["./main"]