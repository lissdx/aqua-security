FROM golang:1.22-alpine AS build

COPY . /go/src/github.com/lissdx/aqua-security/

ENV GOOS=linux
ENV GO111MODULE=auto

WORKDIR /go/src/github.com/lissdx/aqua-security/

RUN go install github.com/lissdx/aqua-security/cmd/monitoring
RUN go install github.com/lissdx/aqua-security/cmd/migration

FROM alpine:latest

WORKDIR /app

COPY --from=build /go/bin/* .
COPY ./configs ./configs
COPY ./db ./db
COPY ./start.sh .
RUN mkdir -p test_dir/input && mkdir -p test_dir/output && chmod +x start.sh


