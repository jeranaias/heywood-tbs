# ============================================================
# Heywood TBS — Multi-stage Docker build
# Same image for commercial Azure + Azure Gov (MCEN/IL5)
# ============================================================

# Stage 1: Build React SPA
FROM node:22-alpine AS web-builder
WORKDIR /build
COPY web/package.json web/package-lock.json ./
RUN npm ci --ignore-scripts
COPY web/ ./
RUN npm run build

# Stage 2: Build Go binary
FROM golang:1.22-alpine AS go-builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ cmd/
COPY internal/ internal/
# GOEXPERIMENT=boringcrypto enables FIPS-validated crypto
# CGO_ENABLED=0 for static binary
RUN CGO_ENABLED=0 GOEXPERIMENT=boringcrypto go build -trimpath \
    -ldflags="-s -w" -o heywood ./cmd/server

# Stage 3: Final minimal image
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 1000 heywood
WORKDIR /app

# Copy artifacts
COPY --from=go-builder /build/heywood .
COPY --from=web-builder /build/dist web/dist/
COPY data/ data/

# Run as non-root (STIG requirement)
USER heywood

EXPOSE 8080

# Health check for Azure App Service
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
    CMD wget -q --spider http://localhost:8080/api/v1/auth/me || exit 1

ENTRYPOINT ["./heywood"]
CMD ["-port", "8080", "-data", "data"]
