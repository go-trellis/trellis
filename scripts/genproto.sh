#!/bin/sh
GO_SOURCE_DIR=$GOPATH/src
DIR=`pwd`

echo "DIR ${DIR}"
echo "GO_SOURCE_DIR ${GO_SOURCE_DIR}"

find $DIR -path ${DIR}/vendor -prune -o -name '*.pb.go' -exec rm {} \;
find $DIR/proto -path ${DIR}/vendor -prune -o -name '*.proto' \
  -exec protoc --go_out=plugins=grpc:${GO_SOURCE_DIR} -I=${GO_SOURCE_DIR} {} \;
find $DIR -path ${DIR}/vendor -prune -o -name '*.pb.go' -exec ./proto/bin/protoc-go-inject-tag -input={} \;

