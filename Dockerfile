FROM golang:latest

RUN apt-get update && apt-get install -y git nano bash \
    && go get github.com/lib/pq \
    && go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate
WORKDIR /usr/src/app

COPY go.mod go.sum /usr/src/app/

RUN go mod download

ENV GO111MODULE on