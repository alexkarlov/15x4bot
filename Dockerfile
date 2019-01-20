FROM golang:1.10

WORKDIR /go/src/reminder

COPY ./ ./

RUN set -e; \
    apt update; \
    apt install -y postgresql-client curl

CMD [ "go", "run", "./main.go"]
