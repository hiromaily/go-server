#!/bin/sh

###########################################################
# Variable
###########################################################
#export GOTRACEBACK=single
GOTRACEBACK=all
HOSTADD=http://127.0.0.1:8080

###########################################################
# Update all package
###########################################################
#go get -u -v ./...

#boom like Apache Bench
#go get -u github.com/rakyll/boom
#go get -u github.com/tsenart/vegeta

###########################################################
# Adjust version dependency of projects
###########################################################
#cd ${GOPATH}/src/github.com/aws/aws-sdk-go
#git checkout v0.9.17
#git checkout master


###########################################################
# go fmt and go vet
###########################################################
go fmt ./...
go vet ./...
#go vet `go list ./... | grep -v '/vendor/'`


###########################################################
# go build and install
###########################################################
go build -i -v -o ./server ./server.go


###########################################################
# go run
###########################################################
./server

exit 0

###########################################################
# Profiling procedure And Test
###########################################################
#${GOPAHT}/bin/boom -n 1000 -c 100 http://localhost:8080
${GOPAHT}/bin/boom -n 100000 -c 100 ${HOSTADD}

#Access below on browser
${HOSTADD}/debug/pprof/profile

# Check Profile
#  usage: pprof [options] [binary] <profile source> ...
go tool pprof server ${HOME}/Downloads/profile

#e.g.
#(pprof)top 7

#oepn another tab on terminal and execute it.
#go tool pprof http://127.0.0.1:8080/debug/pprof/profile

