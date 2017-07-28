# go-server

[![Go Report Card](https://goreportcard.com/badge/github.com/hiromaily/go-server)](https://goreportcard.com/report/github.com/hiromaily/go-server)
[![codebeat badge](https://codebeat.co/badges/4d6a94a0-529b-43e2-88fb-bfa1b8efbdb8)](https://codebeat.co/projects/github-com-hiromaily-go-server-master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/5867b50f6ce54a668f660d78a28b1c29)](https://www.codacy.com/app/hiromaily2/go-server?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=hiromaily/go-server&amp;utm_campaign=Badge_Grade)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://raw.githubusercontent.com/hiromaily/go-gin-wrapper/master/LICENSE)


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


## Web
### functionality
* http/2 based
* Web Push Notifications

