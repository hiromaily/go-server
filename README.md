# go-server
web framework including http/2 functionalities.

## Installation
```
$ go get github.com/hiromaily/go-server ./...

$ docker-compose build
$ docker-compose up

```

## Docker Container Composition
### 1.web
Golang web framework with sample code

### 2.headless
Headless Chromium is used for end to end testing.
[headless](https://chromium.googlesource.com/chromium/src/+/lkgr/headless/README.md)

### 3.devtool
This is test code using [chromedp](https://github.com/knq/chromedp) on headless Chromium

