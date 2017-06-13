//go:generate protoc --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go --proto_path=../../..:. --gogofaster_out=plugins=grpc:. dlframework.proto

package dlframework
