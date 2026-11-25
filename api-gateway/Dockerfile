
# syntax=docker/dockerfile:1

# ---------- build stage ----------
FROM golang:1.26-alpine AS builder
WORKDIR /src

# 1. Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# 2. Copy source code
COPY shared/ ./shared/
COPY api-gateway/ ./api-gateway/

# 3. Build
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
    -o /out/api-gateway ./api-gateway/cmd

# ---------- runtime stage ----------
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/api-gateway /api-gateway
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/api-gateway"]
