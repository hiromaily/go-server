# Note: tabs by space can't not used for Makefile!

###############################################################################
# Docker
###############################################################################
up:
	docker-compose up

up2:
	#This doesn't work yet.
	# API
	docker exec -it gonats_api_1 bash -c "go get -d -v"
	docker exec -it gonats_api_1 bash -c "go build -o /go/bin/api main.go"
	#docker exec -it gonats_api_1 bash -c "kill -HUP 1"

	# WORKER
	docker exec -it gonats_worker_1 bash -c "go get -d -v"
	docker exec -it gonats_worker_1 bash -c "go build -o /go/bin/worker main.go"
	#docker exec -it gonats_worker_1 bash -c "kill -HUP 1"

up_product:
	docker-compose -f docker-compose.yml up

dcbld:
	docker-compose build --no-cache

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
# Build Local
###############################################################################
keygen:
	sudo go run /usr/local/Cellar/go/1.8/libexec/src/crypto/tls/generate_cert.go --host hy


###############################################################################
# Build Local
###############################################################################
bld:
	go build -i -v -o ${GOPATH}/bin/goserver ./cmd/


###############################################################################
# Execution Local
###############################################################################
run:
	#sudo go run ./cmd/main.go
	go build -i -v -o ${GOPATH}/bin/goserver ./cmd/
	sudo goserver

exec:
	sudo goserver

# https://localhost/


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
