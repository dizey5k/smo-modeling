FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o generator ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/generator .
EXPOSE 8081
CMD ["./generator"]