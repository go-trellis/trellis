.PHONY: proto build

proto:
	sh ./scripts/proto.sh

build:
	go build ./...

test:
	go test -v ./...