# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.23-alpine AS builder

# Install git for version info and ca-certificates for HTTPS
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build arguments for version info
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

# Build static binary
RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w \
        -X github.com/anowarislam/ado/internal/meta.Version=${VERSION} \
        -X github.com/anowarislam/ado/internal/meta.Commit=${COMMIT} \
        -X github.com/anowarislam/ado/internal/meta.BuildTime=${BUILD_TIME}" \
    -o /ado ./cmd/ado

# Runtime stage - using scratch for minimal image
FROM scratch

# Copy CA certificates for HTTPS support
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /ado /ado

# Run as non-root (numeric UID for scratch compatibility)
USER 65534:65534

ENTRYPOINT ["/ado"]
