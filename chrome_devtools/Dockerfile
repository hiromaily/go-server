FROM golang:1.8

#RUN echo $GOPATH => /go

RUN apt-get -y update && apt-get install -y wget

RUN mkdir -p /workspace/images
RUN mkdir -p /go/src/github.com/hiromaily/go-server/chrome_devtools

WORKDIR /go/src/github.com/hiromaily/go-server/chrome_devtools

COPY main.go .
COPY pages/ ./pages/

#RUN go get -u github.com/knq/chromedp
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/devtool ./main.go

#CMD ["/go/bin/devtool"]
