# first stage
FROM golang:1.19-alpine as builder
WORKDIR /app
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# second stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Install tzdata
RUN apk add --no-cache tzdata

EXPOSE 4001
EXPOSE 50051

CMD ["./main"]
