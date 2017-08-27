FROM golang:1.9

#RUN echo $GOPATH => /go
ARG encKey
ARG encIV
ENV ENC_KEY=${encKey}
ENV ENC_IV=${encIV}

RUN apt-get -y update && apt-get install -y git

RUN mkdir -p /go/src/github.com/hiromaily/go-server/tmp/log

WORKDIR /go/src/github.com/hiromaily/go-server
COPY . .

RUN go get -u github.com/hiromaily/fresh && \
go get -d -v ./cmd/
#go get -u github.com/hiromaily/go-server/...


RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/go-server ./cmd/main.go

EXPOSE 8080
#CMD ["/go/bin/go-server -p 80"]
