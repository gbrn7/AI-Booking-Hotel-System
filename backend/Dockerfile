# Stage 1: Build binary
FROM docker.io/library/golang:alpine AS builder
WORKDIR /app
ENV GOTOOLCHAIN=auto
RUN apk add --no-cache git ca-certificates tzdata
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/cli ./cmd/cli

# Stage 2: Minimal runtime
FROM docker.io/library/alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/bin/server /app/server
COPY --from=builder /app/bin/cli /app/cli
COPY --from=builder /app/prompts /app/prompts
EXPOSE 8080
ENTRYPOINT ["/app/server"]
