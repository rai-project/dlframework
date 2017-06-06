//go:generate rice -i ./models embed-go
//go:generate protoc --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go --proto_path=../../../..:../../../../github.com:. --gogoslick_out=plugins=grpc:. model.proto

package mxnet
