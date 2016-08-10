#!/bin/sh
###
# initialize for docker environment
###

go get -d -v ./...
go run server.go
