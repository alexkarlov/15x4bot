FROM golang:1.10 as builder

WORKDIR /go/src/github.com/alexkarlov/15x4bot
COPY ./ ./
RUN go get -d -v ./... && \
    go install -v ./... && \
    go test ./... && \
    CGO_ENABLED=0 go build -o 15x4bot .

FROM alpine:3.10 as production

RUN addgroup -S gogroup && adduser -S gorunner -G gogroup && \
    apk update && apk add ca-certificates && rm -rf /var/cache/apk/* \
    && apk add --no-cache tzdata
USER gorunner
WORKDIR /home/gorunner
COPY --from=builder /go/src/github.com/alexkarlov/15x4bot/15x4bot .
CMD [ "./15x4bot" ]
