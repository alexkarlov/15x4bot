FROM golang:1.10

WORKDIR /go/src/github.com/alexkarlov/15x4bot
COPY ./ ./
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["15x4bot"]
