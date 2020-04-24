# Building stage
FROM golang:1.13-alpine3.10

RUN apk add --update make

WORKDIR /go/src/github.com/duyanghao/velero-volume-controller

# Source code, building tools and dependences
COPY cmd /go/src/github.com/duyanghao/velero-volume-controller/cmd
COPY pkg /go/src/github.com/duyanghao/velero-volume-controller/pkg
COPY Makefile /go/src/github.com/duyanghao/velero-volume-controller

COPY go.mod /go/src/github.com/duyanghao/velero-volume-controller
COPY go.sum /go/src/github.com/duyanghao/velero-volume-controller

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GO111MODULE=on

ENV TIMEZONE "Asia/Shanghai"

RUN make build
RUN mv build/velero-volume-controller/velero-volume-controller /go/bin

# Production stage
FROM alpine:3.10
WORKDIR /go/bin

# copy the go binaries from the building stage
COPY --from=0 /go/bin /go/bin

# copy the config files from the current working dir
COPY examples /velero-volume-controller/examples

ENTRYPOINT ["./velero-volume-controller", "-c", "/velero-volume-controller/examples/config.yml", "-v", "4", "-logtostderr", "true"]
