# syntax=docker/dockerfile:1

# ---------- build stage ----------
FROM golang:1.26-alpine AS builder
WORKDIR /src

# 1. Dependencies first (cached until go.mod/go.sum change)
COPY go.mod go.sum ./
RUN go mod download

# 2. Only what this service needs: shared + this service
COPY shared/ ./shared/
COPY user-service/ ./user-service/

# 3. Build
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
    -o /out/user-service ./user-service/cmd

# ---------- runtime stage ----------
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/user-service /user-service
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/user-service"]