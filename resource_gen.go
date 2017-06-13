//go:generate protoc --plugin=protoc-gen-go=${GOPATH}/bin/protoc-gen-go --proto_path=../../..:. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. --gogofaster_out=plugins=grpc:. dlframework.proto

package dlframework
