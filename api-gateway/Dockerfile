FROM golang:1.12-alpine as build-env

ENV TZ Asia/Singapore
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV GO111MODULE=on

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev

RUN mkdir /api
RUN mkdir -p /api/proto 

WORKDIR /api


COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./proto ./proto
COPY ./api.go .

RUN go build -o api-gateway .

CMD ./api-gateway