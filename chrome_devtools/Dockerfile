FROM golang:1.8

#RUN echo $GOPATH => /go

RUN apt-get -y update && apt-get install -y wget

RUN mkdir -p /workspace && mkdir -p /tmp/webdrive/images
WORKDIR /workspace

#RUN go get -u github.com/knq/chromedp
COPY main.go .

RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/devtool ./main.go

#CMD ["/go/bin/devtool"]
