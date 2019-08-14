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
RUN apk update
RUN apk add --no-cache ca-certificates tzdata
RUN update-ca-certificates 2>/dev/null || true

# -----------------------------------------------------------------------------
# Stage application
# -----------------------------------------------------------------------------
# Agent binaries directory should be set as build context
COPY prometheus-eventhubs-adapter /var/tmp/
COPY prometheus-eventhubs-adapter.toml /var/tmp/

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
      description="Uses the remote write features of Prometheus to send samples intended for processing and storage on Azure."

# -----------------------------------------------------------------------------
# Import from builder
# -----------------------------------------------------------------------------
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /var/tmp/prometheus-eventhubs-adapter /usr/bin/
COPY --from=builder /var/tmp/prometheus-eventhubs-adapter.toml /etc/prometheus-eventhubs-adapter/

# -----------------------------------------------------------------------------
# Provide a non-root user to run the process.
# -----------------------------------------------------------------------------
RUN echo "===> Add user prometheus ....." && \
    mkdir -p /var/lib/prometheus-eventhubs-adapter && \
    groupadd --gid 1001 prometheus && \
    adduser --uid 1001 --gid 1001 \
      --home-dir /var/lib/prometheus-eventhubs-adapter \
      --no-create-home \
      --comment "prometheus services" \
      prometheus && \
    echo "===> Set permissions ....." && \
    chown -R 1001:1001 /usr/bin/prometheus-eventhubs-adapter /var/lib/prometheus-eventhubs-adapter /etc/prometheus-eventhubs-adapter && \
    chmod -R 755 /usr/bin/prometheus-eventhubs-adapter /var/lib/prometheus-eventhubs-adapter && \
    chmod -R 644 /etc/prometheus-eventhubs-adapter

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
ENTRYPOINT ["/usr/bin/prometheus-eventhubs-adapter", "--listen_address=:9201"]
