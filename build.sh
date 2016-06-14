#!/usr/bin/env bash

#update dependency packages
#go get -v
#go get -u github.com/golang/lint/golint

#export GOPATH=${HOME}/miidas/miidas_go

#cd ${GOPATH}/src/github.com/aws/aws-sdk-go
#git checkout v0.9.17

#cd ${GOPATH}/src/

#prerequisite build
go fmt ./...
go vet ./...
#go vet `go list ./... | grep -v '/vendor/'`
#golint ./...


# pkg install
# go build -i

#go install pkgname
#go build -i xxxx

#build
go build -o ./server ./server.go

#../bin/main

#oepn another tab on terminal and execute it.
#go tool pprof http://127.0.0.1:8080/debug/pprof/profile