FROM golang:1.10

WORKDIR /go/src/github.com/15x4bot
COPY ./ ./

CMD [ "go", "run", "./main.go"]
