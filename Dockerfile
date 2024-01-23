
FROM golang:1.21.6 AS builder

WORKDIR /app

COPY ./src/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o main -ldflags "-s" cmd/backend/main.go


FROM alpine:latest

WORKDIR /app

# ENV LOG_LEVEL=4
ENV DEBUG=false

COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main"]
