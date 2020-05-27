# Note: tabs by space can't not used for Makefile!

###############################################################################
# PKG Dependencies
###############################################################################
update:
	go get -u github.com/rakyll/hey
	go get -u github.com/davecheney/httpstat
	go get -u -d -v ./...


keygen:
	sudo go run ${GOROOT}/src/crypto/tls/generate_cert.go --host hy
	#sudo go run /usr/local/Cellar/go/1.8/libexec/src/crypto/tls/generate_cert.go --host hy

###############################################################################
# Golang formatter and detection
###############################################################################
.PHONY: imports
imports:
	./scripts/imports.sh

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lintfix
lintfix:
	golangci-lint run --fix


###############################################################################
# Build Local
###############################################################################
bld:
	go build -i -v -o ${GOPATH}/bin/goserver ./cmd/


###############################################################################
# Execution Local
###############################################################################
run:
	go run ./cmd/main.go -p 8080 -f ./data/config.toml
	#go build -i -race -v -o ${GOPATH}/bin/goserver ./cmd/
	#goserver -p 8080 -f data/config.toml

run2:
	#sudo go run ./cmd/main.go
	go build -i -race -v -o ${GOPATH}/bin/goserver ./cmd/
	sudo goserver -tsl 1

exec:
	goserver -p 8080 -f ./data/config.toml

exec2:
	sudo goserver -tsl 1 -f ./data/config.toml


###############################################################################
# Docker
###############################################################################
up:
	docker-compose up

up2:
	docker-compose up web


serverin:
	docker exec -it go-server bash -c "echo ${GOROOT}"

up_product:
	docker-compose -f docker-compose.yml up

dcbld:
	docker-compose build --no-cache


###############################################################################
# Test
###############################################################################
test:
	go test -run TestGetRequestOnTable -v ./cmd/*.go -f ../data/config.toml


###############################################################################
# Build on Docker
###############################################################################
devtool:
	docker exec -it devtool bash -c "devtool -d headless -n localhost -h chrome-headless"

devtoolbld:
	docker exec -it devtool bash -c "CGO_ENABLED=0 GOOS=linux go build -race -o /go/bin/devtool ./main.go"

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
