FROM golang:1.12-alpine as build-env

ENV GO111MODULE=on

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev

RUN mkdir /user-service
RUN mkdir -p /user-service/proto 

WORKDIR /user-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./proto/user.pb.go ./proto
COPY ./server/service.go .


RUN go build -o user-service .

CMD ./user-service