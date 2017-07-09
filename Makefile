# Note: tabs by space can't not used for Makefile!

###############################################################################
# PKG Dependencies
###############################################################################
update:
	go get -u github.com/golang/dep/...
	go get -u github.com/rakyll/hey
	go get -u github.com/davecheney/httpstat
	go get -u github.com/client9/misspell/cmd/misspell
	go get -u github.com/gordonklaus/ineffassign
	go get -u github.com/pilu/fresh

	go get -u github.com/alecthomas/gometalinter
	#gometalinter --install

	go get -d -v ./cmd/

depinit:
	cd cmd/;dep init

dep:
	cd cmd/;dep ensure -update

depcln:
	cd cmd/;rm -rf vendor lock.json manifest.json


###############################################################################
# Golang formatter and detection
###############################################################################
fmt:
	go fmt `go list ./... | grep -v '/vendor/'`

vet:
	go vet `go list ./... | grep -v '/vendor/'`

fix:
	go fix `go list ./... | grep -v '/vendor/'`

lint:
	golint ./... | grep -v '^vendor\/' || true
	misspell `find . -name "*.go" | grep -v '/vendor/'`
	ineffassign .

chk:
	go fmt `go list ./... | grep -v '/vendor/'`
	go vet `go list ./... | grep -v '/vendor/'`
	go fix `go list ./... | grep -v '/vendor/'`
	golint ./... | grep -v '^vendor\/' || true
	misspell `find . -name "*.go" | grep -v '/vendor/'`
	ineffassign .


###############################################################################
# Docker
###############################################################################
up:
	docker-compose up

serverin:
	docker exec -it go-server bash -c "echo ${GOROOT}"

up_product:
	docker-compose -f docker-compose.yml up

dcbld:
	docker-compose build --no-cache


###############################################################################
# Initial Settings
###############################################################################
keygen:
	sudo go run ${GOROOT}/src/crypto/tls/generate_cert.go --host hy
	#sudo go run /usr/local/Cellar/go/1.8/libexec/src/crypto/tls/generate_cert.go --host hy

#submodule:
#    git submodule add https://github.com/BuckyMaler/global.git submodules/global
#	#ln -s ${GOPATH}/src/github.com/hiromaily/go-goa/goa/swagger ./public/swagger


###############################################################################
# Build Local
###############################################################################
bld:
	go build -i -v -o ${GOPATH}/bin/goserver ./cmd/

bld2:
	go build -i -v -o ${GOPATH}/bin/devtool ./chrome_devtools/


###############################################################################
# Execution Local
###############################################################################
run:
	#sudo go run ./cmd/main.go
	go build -i -v -o ${GOPATH}/bin/goserver ./cmd/
	goserver -p 8080 -f data/config.toml

run2:
	#sudo go run ./cmd/main.go
	go build -i -v -o ${GOPATH}/bin/goserver ./cmd/
	sudo goserver -tsl 1

exec:
	goserver -p 8080

exec2:
	sudo goserver -tsl 1


###############################################################################
# Build on Docker
###############################################################################
devtool:
	docker exec -it devtool bash -c "devtool -d headless -n localhost -h chrome-headless"

devtoolbld:
	docker exec -it devtool bash -c "CGO_ENABLED=0 GOOS=linux go build -o /go/bin/devtool ./main.go"

devtoolin:
	docker exec -it devtool bash

# https://localhost/
# https://localhost:8080/

###############################################################################
# Clean
###############################################################################
cln:
	go clean -n

clnok:
	go clean


###############################################################################
# httpie
###############################################################################
http:
	# httpie #brew install httpie
	# jq     #brew install jq

	# Static files
	http localhost:8080/

###############################################################################
# Bench
###############################################################################
bench:
	hey -n 20000 -c 50 -m GET http://localhost:8080/api/user


###############################################################################
# HTTP Stat
###############################################################################
httpstat:
	httpstat http://localhost:8080/api/user
