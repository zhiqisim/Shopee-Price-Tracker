FROM golang:1.12-alpine as build-env

ENV GO111MODULE=on

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev

RUN mkdir /user-service
RUN mkdir -p /user-service/proto 

WORKDIR /user-service

COPY ./proto/service.pb.go /user-service/proto
COPY ./server/service.go .
COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go build -o user-service .

CMD ./user-service