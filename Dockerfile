FROM golang:1.22.4-alpine3.20 AS builder

COPY . /app/
WORKDIR /app/

ENV CGO_ENABLED=1

RUN apk add cmake make gcc libtool musl-dev

RUN go build -o ./.bin/app ./cmd/app/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/.bin/app .

CMD ["./app"]
