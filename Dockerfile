FROM golang:1.20-bullseye as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/

# Path: Dockerfile
FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/bin/ /app/
COPY --from=builder /app/config.toml /app/

ENTRYPOINT ["/app/bin/main"]
#ENTRYPOINT ["tail", "-f", "/dev/null"]