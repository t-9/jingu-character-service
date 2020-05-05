FROM golang:1.14.2-alpine3.11

RUN apk update && \
    apk add git

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/t-9/jingu-character-service

RUN dep init