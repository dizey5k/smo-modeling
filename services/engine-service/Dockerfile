FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o engine ./cmd/server/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/engine .
EXPOSE 8080
CMD ["./engine"]