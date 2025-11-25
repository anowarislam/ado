# GoReleaser Dockerfile - uses pre-built binaries
# This is used by GoReleaser for efficient multi-arch builds

FROM scratch

# Copy CA certificates for HTTPS support
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the pre-built binary (GoReleaser provides this)
COPY ado /ado

# Run as non-root (numeric UID for scratch compatibility)
USER 65534:65534

ENTRYPOINT ["/ado"]
