FROM golang:alpine as builder

WORKDIR /code/go-tools
ADD . ./

ENV GO111MODULE on
ENV GOPROXY https://goproxy.io

RUN  CGO_ENABLED=0 GOOS=linux go build -o gotool -ldflags '-s -w'


FROM alpine:3.8

ENV TZ Asia/Shanghai

WORKDIR /code/go-tools

COPY --from=builder /code/go-tools/gotool .
CMD  ["sh", "-c", "./gotool"]

