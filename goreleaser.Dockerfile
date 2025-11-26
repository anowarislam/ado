# GoReleaser Dockerfile - uses pre-built binaries
# This is used by GoReleaser for efficient multi-arch builds

FROM scratch

# Copy CA certificates for HTTPS support (pinned version for reproducibility)
COPY --from=alpine:3.21 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data for proper time handling
COPY --from=alpine:3.21 /usr/share/zoneinfo /usr/share/zoneinfo

# Copy license and documentation
COPY LICENSE /LICENSE
COPY README.md /README.md

# Copy the pre-built binary (GoReleaser provides this)
COPY ado /ado

# Run as non-root (numeric UID for scratch compatibility)
USER 65534:65534

ENTRYPOINT ["/ado"]
