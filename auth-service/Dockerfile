# syntax=docker/dockerfile:1

# ---------- build stage ----------
FROM golang:1.23-alpine AS builder
WORKDIR /src

# 1. Dependencies first (cached until go.mod/go.sum change)
COPY go.mod go.sum ./
RUN go mod download

# 2. Only what this service needs: shared + this service
COPY shared/ ./shared/
COPY auth-service/ ./auth-service/

# 3. Build
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
    -o /out/auth-service ./auth-service/cmd/server

# ---------- runtime stage ----------
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/auth-service /auth-service
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/auth-service"]