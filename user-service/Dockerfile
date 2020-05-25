FROM golang:1.12-alpine as build-env


ENV TZ Asia/Singapore
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV GO111MODULE=on

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev

RUN mkdir /user-service
RUN mkdir -p /user-service/proto 

WORKDIR /user-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./proto/user.pb.go ./proto
COPY ./service.go .


RUN go build -o user-service .

CMD ./user-service