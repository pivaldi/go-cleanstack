# Dockerfile

# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/cleanstack ./main.go

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/bin/cleanstack /usr/local/bin/cleanstack
COPY --from=builder /app/config_development.toml /config_development.toml
COPY --from=builder /app/config_default.toml /config_default.toml

EXPOSE 4224

ENTRYPOINT ["cleanstack"]
CMD ["serve"]
