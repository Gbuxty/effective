FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /effective ./cmd/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /effective .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env ./

EXPOSE 8080
CMD ["./effective"]