
FROM golang:1.12-alpine
# install glide
# RUN go get github.com/Masterminds/glide
# create a working directory
WORKDIR /items-service
# add glide.yaml and glide.lock
# ADD glide.yaml glide.yaml
# ADD glide.lock glide.lock
# install packages
# RUN glide install
# add source code
ADD src src
# build main.go
RUN go build src/service.go
# run the binary
CMD ["./service"]