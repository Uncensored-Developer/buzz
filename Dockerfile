FROM golang:1.22.4-alpine3.20 AS builder

COPY . /app/
WORKDIR /app/

RUN go build -o ./.bin/app ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/.bin/app .

CMD ["./app"]
