all: build run clean
.PHONY: all

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo -o lbh main.go
	chmod a+x lbh

run:
	sudo ${ROOT_DIR}/lbh --port 443 --tlsCertFile ${ROOT_DIR}/pki/server.crt --tlsKeyFile ${ROOT_DIR}/pki/server.key

clean:
	rm -rf lbh