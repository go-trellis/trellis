.PHONY: proto build

proto:
	sh ./scripts/genproto.sh

build:
	go build ./...

test:
	go test -v ./...