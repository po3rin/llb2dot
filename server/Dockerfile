# for NOW deploy
FROM golang:alpine
ADD . /go/src/github.com/po3rin/llb2dot/server
RUN go install github.com/po3rin/llb2dot/server
CMD ["/go/bin/server"]
EXPOSE 8080
