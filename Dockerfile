FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o task-manager-mcp .

FROM scratch
COPY --from=builder /app/task-manager-mcp /task-manager-mcp

EXPOSE 8080

ENTRYPOINT ["/task-manager-mcp"]
