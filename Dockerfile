FROM golang:1.8

#RUN echo $GOPATH => /go

RUN apt-get -y update && apt-get install -y git

RUN go get -u github.com/hiromaily/fresh && \
go get -u github.com/hiromaily/go-server/...

RUN mkdir -p /go/src/github.com/hiromaily/go-server/tmp/log

WORKDIR /go/src/github.com/hiromaily/go-server
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/go-server ./cmd/main.go

EXPOSE 80
#CMD ["/go/bin/go-server -p 80"]
