# ------------------------------------------------
# Build stage (official Golang image)
# ------------------------------------------------
FROM golang:1.25-bookworm AS build

WORKDIR /app

# Install librdkafka build deps (glibc based)
RUN apt-get update && apt-get install -y \
    gcc \
    libc6-dev \
    librdkafka-dev \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*


# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary (CGO REQUIRED for Kafka IAM)
RUN CGO_ENABLED=1 \
    go build -o app ./cmd/main.go

# ------------------------------------------------
# Runtime stage (Amazon Linux)
# ------------------------------------------------
FROM amazonlinux:2023

WORKDIR /app

# Install runtime dependencies only
RUN dnf install -y ca-certificates && dnf clean all

# Copy binary from build stage
COPY --from=build /app/app .

# Run as non-root (ECS best practice)
USER 10001

CMD ["./app"]