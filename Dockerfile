# ---------- builder ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Add git for module dependencies if needed
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Add build flags for optimization
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o api ./cmd

# ---------- runtime ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/database/migrations ./database/migrations

EXPOSE 8080

CMD ["./api"]