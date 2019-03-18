FROM golang:1.12 AS stage0

WORKDIR /go
ADD ./ /go
RUN go build -o stage0_bin

FROM golang:1.12 AS stage1

WORKDIR /go
ADD ./ /go
RUN go build -o stage1_bin

FROM golang:1.12

COPY --from=stage0 /go/stage0_bin /
COPY --from=stage1 /go/stage1_bin /
