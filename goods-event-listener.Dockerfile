# Build stage
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o goods-event-listener cmd/goods-event-listener/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/goods-event-listener ./
CMD ["./goods-event-listener"]
ENTRYPOINT ["./goods-event-listener"]