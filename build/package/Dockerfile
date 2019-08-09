# =============================================================================
# bryanklewis/prometheus-eventhubs-adapter
#
# Project:
#     https://github.com/bryanklewis/prometheus-eventhubs-adapter
#
#
# =============================================================================
#
# -----------------------------------------------------------------------------
# Base image (builder)
# -----------------------------------------------------------------------------
FROM golang:alpine AS builder

# -----------------------------------------------------------------------------
# Install system packages
# -----------------------------------------------------------------------------
RUN echo "===> Updating Alpine Linux ....." && \
    apk update && \
    echo "===> Installing additional packages ....." && \
    apk add --no-cache \
      ca-certificates \
      tzdata && \
    echo "===> Update certificates folder ....." && \
    update-ca-certificates

# -----------------------------------------------------------------------------
# Base image (run)
# -----------------------------------------------------------------------------
FROM scratch

# -----------------------------------------------------------------------------
# Set image metadata
# -----------------------------------------------------------------------------
LABEL name="bryanklewis/prometheus-eventhubs-adapter" \
      maintainer="Bryan Lewis <dbre@micron.com>" \
      vendor="Micron Technology, Inc." \
      summary="A Prometheus remote storage adapter for Azure Event Hubs." \
      description="prometheus-eventhubs-adapter is a Prometheus remote storage adapter for Azure Event Hubs. The remote write features of Prometheus allow transparently sending samples intended for long term storage."

# -----------------------------------------------------------------------------
# Import from builder
# -----------------------------------------------------------------------------
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# -----------------------------------------------------------------------------
# Add application
# -----------------------------------------------------------------------------


# -----------------------------------------------------------------------------
# Provide a non-root user to run the process.
# -----------------------------------------------------------------------------
RUN echo "===> Add user prometheus ....." && \
    groupadd --gid 1001 prometheus && \
    adduser --uid 1001 --gid 1001 \
      --home-dir /usr/share/prometheus --no-create-home \
      prometheus

# -----------------------------------------------------------------------------
# Set path
# -----------------------------------------------------------------------------
WORKDIR /app

# -----------------------------------------------------------------------------
# Define execution user
# -----------------------------------------------------------------------------
USER prometheus

# -----------------------------------------------------------------------------
# Expose service ports
# -----------------------------------------------------------------------------
EXPOSE 9201 5671 5672

# -----------------------------------------------------------------------------
# Run container
# -----------------------------------------------------------------------------
ENTRYPOINT ["/bin/echo", "-listen_address=:9201"]
