# Stage 1: Builder
FROM --platform=$BUILDPLATFORM golang:1.25-bookworm AS builder

WORKDIR /src

# Build arguments for cross-compilation
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

# Cache module downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Extract version
RUN VERSION=$(cat VERSION | tr -d '[:space:]') && \
    echo "Building version: $VERSION" > /tmp/version

# Build binary from root directory
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    set -eux; \
    mkdir -p /out; \
    GOARM=""; \
    if [ "$TARGETARCH" = "arm" ]; then \
        case "$TARGETVARIANT" in \
            v5) GOARM=5 ;; \
            v6) GOARM=6 ;; \
            v7|"") GOARM=7 ;; \
            *) exit 1 ;; \
        esac; \
    fi; \
    VERSION=$(cat /tmp/version); \
    env CGO_ENABLED=0 \
        GOOS="${TARGETOS:-linux}" \
        GOARCH="${TARGETARCH:-amd64}" \
        GOARM="${GOARM}" \
        go build -trimpath -ldflags="-s -w" -o /out/netdash .

# Stage 2: Runtime image
FROM debian:bookworm-slim

ENV DEBIAN_FRONTEND=noninteractive \
    TZ=UTC

# Install system dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        tzdata \
    && rm -rf /var/lib/apt/lists/*

ARG TARGETARCH
ARG TARGETVARIANT

# Download and install Ookla Speedtest CLI
RUN set -eux; \
    case "${TARGETARCH}/${TARGETVARIANT}" in \
        "amd64/")   OOKLA_ARCH="x86_64" ;; \
        "arm64/")   OOKLA_ARCH="aarch64" ;; \
        "arm/v7")   OOKLA_ARCH="armhf" ;; \
        "arm/v6")   OOKLA_ARCH="armel" ;; \
        "arm/v5")   OOKLA_ARCH="armel" ;; \
        *)          exit 1 ;; \
    esac; \
    DOWNLOAD_URL=$(curl -sL https://www.speedtest.net/apps/cli | \
        grep -oE "https://install\.speedtest\.net/app/cli/ookla-speedtest-[0-9.]+-linux-${OOKLA_ARCH}\.tgz" | \
        head -n 1); \
    [ -n "$DOWNLOAD_URL" ] || exit 1; \
    curl -fSL "$DOWNLOAD_URL" -o speedtest.tgz; \
    tar -xvf speedtest.tgz speedtest; \
    mv speedtest /usr/local/bin/speedtest; \
    rm speedtest.tgz; \
    speedtest --version

# Set working directory
WORKDIR /netdash

# Copy binary
COPY --from=builder /out/netdash /netdash/netdash

# Expose application port
EXPOSE 80

# Set entrypoint
ENTRYPOINT ["/netdash/netdash"]