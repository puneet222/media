# base go image

FROM golang:1.19-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o monk ./cmd/api

RUN chmod +x /app/monk

# build a tiny docker image that run the binary from builder image

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/monk /app

CMD [ "/app/monk" ]