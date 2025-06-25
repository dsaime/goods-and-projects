# Build stage
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o goods-and-projects cmd/goods-and-projects/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/goods-and-projects ./
CMD ["./goods-and-projects"]
ENTRYPOINT ["./goods-and-projects"]